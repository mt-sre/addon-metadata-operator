//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/mt-sre/go-ci/command"
	"github.com/mt-sre/go-ci/container"
	"github.com/mt-sre/go-ci/file"
	"github.com/mt-sre/go-ci/web"
	"go.uber.org/multierr"
)

var _depBin = filepath.Join(_dependencyDir, "bin")

var _dependencyDir = func() string {
	if dir, ok := os.LookupEnv("DEPENDENCY_DIR"); ok {
		return dir
	}

	return filepath.Join(_projectRoot, ".cache", "dependencies")
}()

var _projectRoot = func() string {
	if root, ok := os.LookupEnv("PROJECT_ROOT"); ok {
		return root
	}

	toplevel := git(
		command.WithArgs{"rev-parse", "--show-toplevel"},
	)

	if err := toplevel.Run(); err != nil || !toplevel.Success() {
		panic("failed to get working directory")
	}

	return strings.TrimSpace(toplevel.Stdout())
}()

var _tag = func() string {
	shortRev := git(
		command.WithArgs{"rev-parse", "--short", "HEAD"},
	)

	if err := shortRev.Run(); err != nil || !shortRev.Success() {
		panic("failed to get current tag")
	}

	return strings.TrimSpace(shortRev.Stdout())
}()

var git = command.NewCommandAlias("git")

var Aliases = map[string]interface{}{
	"check":     All.Check,
	"clean":     All.Clean,
	"run-hooks": Hooks.Run,
	"test":      All.Test,
}

type All mg.Namespace

func (All) Check(ctx context.Context) {
	mg.SerialCtxDeps(
		ctx,
		Check.Tidy,
		Check.Verify,
		Check.Lint,
	)
}

func (All) Clean(ctx context.Context) {
	mg.CtxDeps(
		ctx,
		Build.CleanCLI,
		Build.CleanOperator,
		Test.Clean,
	)
}

func (All) Test(ctx context.Context) {
	mg.CtxDeps(
		ctx,
		Test.Unit,
		Test.Integration,
	)
}

type Check mg.Namespace

func (Check) Tidy(ctx context.Context) error {
	tidy := gocmd(
		command.WithArgs{"mod", "tidy"},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := tidy.Run(); err != nil {
		return fmt.Errorf("starting tidy: %w", err)
	}

	if tidy.Success() {
		return nil
	}

	return fmt.Errorf("running tidy: %w", tidy.Error())
}

func (Check) Verify(ctx context.Context) error {
	verify := gocmd(
		command.WithArgs{"mod", "verify"},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := verify.Run(); err != nil {
		return fmt.Errorf("starting verification: %w", err)
	}

	if verify.Success() {
		return nil
	}

	return fmt.Errorf("running verification: %w", verify.Error())
}

func (Check) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdateGolangCILint)

	lint := golangci(
		command.WithArgs{"run",
			"--timeout=10m",
			"-E", "unused,gofmt,goimports,gosimple,staticcheck",
			"--skip-dirs-use-default",
			"--verbose",
		},
		command.WithContext{Context: ctx},
	)

	if err := lint.Run(); err != nil {
		return fmt.Errorf("starting linter: %w", err)
	}

	fmt.Fprint(os.Stdout, lint.CombinedOutput())

	if lint.Success() {
		return nil
	}

	return fmt.Errorf("running linter: %w", lint.Error())
}

var golangci = command.NewCommandAlias(filepath.Join(_depBin, "golangci-lint"))

type Test mg.Namespace

func (Test) Unit(ctx context.Context) error {
	test := gocmd(
		command.WithCurrentEnv(true),
		command.WithEnv{
			"CGO_CFLAGS": "-DSQLITE_ENABLE_JSON1",
		},
		command.WithArgs{
			"test", "-count=1", "-race", "-timeout", "15m",
			"./api...",
			"./cmd...",
			"./internal...",
			"./pkg...",
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := test.Run(); err != nil {
		return fmt.Errorf("starting unit tests: %w", err)
	}

	if test.Success() {
		return nil
	}

	return fmt.Errorf("running unit tests: %w", test.Error())
}

func (Test) Integration(ctx context.Context) error {
	e2eBin := filepath.Join(_projectRoot, ".cache/mtcli")

	build := gocmd(
		command.WithCurrentEnv(true),
		command.WithEnv{
			"CGO_ENABLED": "1",
			"CGO_CFLAGS":  "-DSQLITE_ENABLE_JSON1",
		},
		command.WithArgs{"build", "-a", "-o", e2eBin, "./cmd/mtcli"},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := build.Run(); err != nil || !build.Success() {
		if build.Error() != nil {
			err = build.Error()
		}

		return fmt.Errorf("building mtcli binary: %w", err)
	}

	integration := gocmd(
		command.WithCurrentEnv(true),
		command.WithEnv{
			"E2E_MTCLI_PATH": e2eBin,
		},
		command.WithArgs{
			"test", "-count=1", "-race", "./integration...",
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := integration.Run(); err != nil {
		return fmt.Errorf("starting integration tests: %w", err)
	}

	if integration.Success() {
		return nil
	}

	return fmt.Errorf("running integration tests: %w", integration.Error())
}

func (Test) Clean() error {
	files, err := file.Find(_projectRoot,
		file.WithEntType(file.EntTypeDir),
		file.WithName("index_tmp_*"),
	)
	if err != nil {
		return err
	}

	var errCollector error

	for _, f := range files {
		multierr.AppendInto(&errCollector, sh.Rm(f))
	}

	return errCollector
}

const repository = "quay.io/mtsre/addon-metadata-operator"

type Build mg.Namespace

func (Build) CLI(ctx context.Context) error {
	build := gocmd(
		command.WithCurrentEnv(true),
		command.WithEnv{
			"CGO_ENABLED": "1",
			"CGO_CFLAGS":  "-DSQLITE_ENABLE_JSON1",
		},
		command.WithArgs{
			"build", "-a",
			"-o", filepath.Join(_projectRoot, "bin", "mtcli"),
			filepath.Join("cmd", "mtcli", "main.go"),
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := build.Run(); err != nil {
		return fmt.Errorf("starting to build mtcli: %w", err)
	}

	if build.Success() {
		return nil
	}

	return fmt.Errorf("building mtcli: %w", build.Error())
}

func (Build) CleanCLI() error {
	return sh.Rm(filepath.Join(_projectRoot, "bin", "mtcli"))
}

func (Build) Operator(ctx context.Context) error {
	build := gocmd(
		command.WithCurrentEnv(true),
		command.WithEnv{
			"CGO_ENABLED": "0",
		},
		command.WithArgs{
			"build", "-a",
			"-o", filepath.Join(_projectRoot, "bin", "addon-metadata-operator"),
			filepath.Join("cmd", "addon-metadata-operator", "main.go"),
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := build.Run(); err != nil {
		return fmt.Errorf("starting to build amo: %w", err)
	}

	if build.Success() {
		return nil
	}

	return fmt.Errorf("building amo: %w", build.Error())
}

var gocmd = command.NewCommandAlias(mg.GoCmd())

func (Build) CleanOperator() error {
	return sh.Rm(filepath.Join(_projectRoot, "bin", "addon-metadata-operator"))
}

func (Build) OperatorImage(ctx context.Context) error {
	runtime, ok := container.Runtime()
	if !ok {
		return errors.New("could not find container runtime")
	}

	build := command.NewCommand(runtime,
		command.WithArgs{
			"build",
			"-t", fmt.Sprintf("%s:%s", repository, _tag),
			"-f", "Dockerfile.build", _projectRoot,
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := build.Run(); err != nil {
		return fmt.Errorf("starting to build amo image: %w", err)
	}

	if build.Success() {
		return nil
	}

	return fmt.Errorf("building amo image: %w", build.Error())
}

type Generate mg.Namespace

func (Generate) Manifests(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdateControllerGen)

	generate := controllergen(
		command.WithArgs{
			"crd", "rbac:roleName=manager-role",
			"webhook", `paths="./..."`,
			"output:crd:artifacts:config=config/crd/bases",
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := generate.Run(); err != nil {
		return fmt.Errorf("starting to generate manifests: %w", err)
	}

	if generate.Success() {
		return nil
	}

	return fmt.Errorf("generating manifests: %w", generate.Error())
}

func (Generate) Boilerplate(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdateControllerGen)

	generate := controllergen(
		command.WithArgs{
			`object:headerFile="hack/boilerplate.go.txt"`,
			`paths="./..."`,
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := generate.Run(); err != nil {
		return fmt.Errorf("starting to generate boilerplate: %w", err)
	}

	if generate.Success() {
		return nil
	}

	return fmt.Errorf("generating boilerplate: %w", generate.Error())
}

var controllergen = command.NewCommandAlias(filepath.Join(_depBin, "controller-gen"))

type Release mg.Namespace

func (Release) PushOperatorImage(ctx context.Context) {
	mg.SerialCtxDeps(
		ctx,
		Build.OperatorImage,
		mg.F(Release.tagImage, fmt.Sprintf("%s:%s", repository, _tag)),
		mg.F(Release.pushImage, fmt.Sprintf("%s:%s", repository, _tag)),
		mg.F(Release.tagImage, fmt.Sprintf("%s:%s", repository, "latest")),
		mg.F(Release.pushImage, fmt.Sprintf("%s:%s", repository, "latest")),
	)
}

func (Release) tagImage(ctx context.Context, ref string) error {
	tag := command.NewCommand("docker",
		command.WithArgs{"tag", ref},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := tag.Run(); err != nil {
		return fmt.Errorf("starting to tag image %q: %w", ref, err)
	}

	if tag.Success() {
		return nil
	}

	return fmt.Errorf("tagging image %q: %w", ref, tag.Error())
}

func (Release) pushImage(ctx context.Context, ref string) error {
	const creds_var = "DOCKER_CONF"

	creds, ok := os.LookupEnv(creds_var)
	if !ok {
		return fmt.Errorf("%q must be defined", creds_var)
	}

	push := command.NewCommand("docker",
		command.WithArgs{
			fmt.Sprintf("--config=%s", creds),
			"push", ref,
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := push.Run(); err != nil {
		return fmt.Errorf("starting to push image %q: %w", ref, err)
	}

	if push.Success() {
		return nil
	}

	return fmt.Errorf("pushing image %q: %w", ref, push.Error())
}

func (Release) CLI(ctx context.Context) error {
	return runGoreleaser(ctx)
}

func (Release) CLISnapshot(ctx context.Context) error {
	return runGoreleaser(ctx, "--snapshot")
}

func runGoreleaser(ctx context.Context, args ...string) error {
	mg.CtxDeps(
		ctx,
		Release.container,
	)

	runtime, ok := container.Runtime()
	if !ok {
		return errors.New("could not find container runtime")
	}

	run := command.NewCommand(runtime,
		command.WithArgs{
			"run", "--rm",
			"-e", "CGO_ENABLED=1",
			"-e", "CGO_CFLAGS=-DSQLITE_ENABLE_JSON1",
			"-e", fmt.Sprintf("GITHUB_TOKEN=%s", os.Getenv("GITHUB_TOKEN")),
			fmt.Sprintf("amo-release:%s", _tag), "release",
		},
		command.WithArgs(args),
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := run.Run(); err != nil {
		return fmt.Errorf("starting to run goreleaser: %w", err)
	}

	if run.Success() {
		return nil
	}

	return fmt.Errorf("running goreleaser: %w", run.Error())
}

func (Release) container(ctx context.Context) error {
	runtime, ok := container.Runtime()
	if !ok {
		return errors.New("could not find container runtime")
	}

	build := command.NewCommand(runtime,
		command.WithArgs{
			"build",
			"-t", fmt.Sprintf("amo-release:%s", _tag),
			"-f", "Dockerfile.release", _projectRoot,
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := build.Run(); err != nil {
		return fmt.Errorf("starting to build goreleaser: %w", err)
	}

	if build.Success() {
		return nil
	}

	return fmt.Errorf("building goreleaser: %w", build.Error())
}

type Deps mg.Namespace

func (Deps) UpdateControllerGen(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Deps.updateGoDependency, "sigs.k8s.io/controller-tools/cmd/controller-gen"))
}

func (Deps) UpdateGolangCILint(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Deps.updateGoDependency, "github.com/golangci/golangci-lint/cmd/golangci-lint"))
}

func (Deps) UpdateGoImports(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Deps.updateGoDependency, "golang.org/x/tools/cmd/goimports"))
}

func (Deps) UpdateKind(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Deps.updateGoDependency, "sigs.k8s.io/kind"))
}

func (Deps) UpdateKustomize(ctx context.Context) {
	mg.CtxDeps(ctx, mg.F(Deps.updateGoDependency, "sigs.k8s.io/kustomize/kustomize/v4"))
}

func (Deps) updateGoDependency(ctx context.Context, src string) error {
	if err := setupDepsBin(); err != nil {
		return fmt.Errorf("creating dependencies bin directory: %w", err)
	}

	toolsDir := filepath.Join(_projectRoot, "tools")

	tidy := gocmd(
		command.WithArgs{"mod", "tidy"},
		command.WithWorkingDirectory(toolsDir),
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := tidy.Run(); err != nil {
		return fmt.Errorf("starting to tidy tools dir: %w", err)
	}

	if !tidy.Success() {
		return fmt.Errorf("tidying tools dir: %w", tidy.Error())
	}

	install := gocmd(
		command.WithArgs{"install", src},
		command.WithWorkingDirectory(toolsDir),
		command.WithCurrentEnv(true),
		command.WithEnv{"GOBIN": _depBin},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := install.Run(); err != nil {
		return fmt.Errorf("starting to install command from source %q: %w", src, err)
	}

	if !install.Success() {
		return fmt.Errorf("installing command from source %q: %w", src, install.Error())
	}

	return nil
}

func (Deps) UpdatePreCommit(ctx context.Context) error {
	if err := setupDepsBin(); err != nil {
		return fmt.Errorf("creating dependencies bin directory: %w", err)
	}

	const urlPrefix = "https://github.com/pre-commit/pre-commit/releases/download"

	// pinning to version 2.17.0 since 2.18.0+ requires python>=3.7
	const version = "2.17.0"

	out := filepath.Join(_depBin, "pre-commit")

	if _, err := os.Stat(out); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspecting output location %q: %w", out, err)
		}

		if err := web.DownloadFile(ctx, urlPrefix+fmt.Sprintf("/v%s/pre-commit-%s.pyz", version, version), out); err != nil {
			return fmt.Errorf("downloading pre-commit: %w", err)
		}
	}

	return os.Chmod(out, 0775)
}

func setupDepsBin() error {
	return os.MkdirAll(_depBin, 0o774)
}

type Hooks mg.Namespace

func (Hooks) Enable(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	install := precommit(
		command.WithArgs{"install"},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := install.Run(); err != nil {
		return fmt.Errorf("starting to enable hooks: %w", err)
	}

	if install.Success() {
		return nil
	}

	return fmt.Errorf("enabling hooks: %w", install.Error())
}

func (Hooks) Disable(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	uninstall := precommit(
		command.WithArgs{"uninstall"},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := uninstall.Run(); err != nil {
		return fmt.Errorf("starting to disable hooks: %w", err)
	}

	if uninstall.Success() {
		return nil
	}

	return fmt.Errorf("disabling hooks: %w", uninstall.Error())
}

func (Hooks) Run(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	run := precommit(
		command.WithArgs{
			"run",
			"--show-diff-on-failure",
			"--from-ref", "origin/master", "--to-ref", "HEAD",
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := run.Run(); err != nil {
		return fmt.Errorf("starting to run hooks: %w", err)
	}

	if run.Success() {
		return nil
	}

	return fmt.Errorf("running hooks: %w", run.Error())
}

func (Hooks) RunAllFiles(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	runAll := precommit(
		command.WithArgs{
			"run", "--all-files",
		},
		command.WithConsoleOut(mg.Verbose()),
		command.WithContext{Context: ctx},
	)

	if err := runAll.Run(); err != nil {
		return fmt.Errorf("starting to run hooks for all files: %w", err)
	}

	if runAll.Success() {
		return nil
	}

	return fmt.Errorf("running hooks for all files: %w", runAll.Error())
}

var precommit = command.NewCommandAlias(filepath.Join(_depBin, "pre-commit"))
