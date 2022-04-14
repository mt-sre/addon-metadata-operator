//go:build mage
// +build mage

package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"go.uber.org/multierr"
)

const repository = "quay.io/mtsre/addon-metadata-operator"

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

	root, err := sh.Output("git", "rev-parse", "--show-toplevel")
	if err != nil {
		panic("failed to get working directory")
	}

	return root
}()

var _tag = func() string {
	tag, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		panic("failed to get current tag")
	}

	return tag
}()

var Aliases = map[string]interface{}{
	"check":     All.Check,
	"clean":     All.Clean,
	"run-hooks": Hooks.Run,
	"test":      All.Test,
}

type All mg.Namespace

func (All) Check(ctx context.Context) {
	mg.SerialCtxDeps(ctx,
		Check.Tidy,
		Check.Verify,
		Check.Lint,
	)
}

func (All) Clean() {
	mg.Deps(
		Build.CleanCLI,
		Build.CleanOperator,
		Test.Clean,
	)
}

func (All) Test() {
	mg.Deps(
		Test.Unit,
		Test.Integration,
	)
}

type Check mg.Namespace

func (Check) Tidy() error {
	return sh.Run(mg.GoCmd(), "mod", "tidy")
}

func (Check) Verify() error {
	return sh.Run(mg.GoCmd(), "mod", "verify")
}

func (Check) Lint(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdateGolangCILint)

	return sh.Run(filepath.Join(_depBin, "golangci-lint"), "run",
		"--timeout=10m",
		"-E", "unused,gofmt,goimports,gosimple,staticcheck",
		"--skip-dirs-use-default",
		"--verbose",
	)
}

type Test mg.Namespace

func (Test) Unit() error {
	return sh.RunWith(map[string]string{
		"CGO_CFLAGS": "-DSQLITE_ENABLE_JSON1",
	}, mg.GoCmd(), "test",
		"./api...",
		"./cmd...",
		"./internal...",
		"./pkg...",
	)
}

func (Test) Integration() error {
	e2eBin := filepath.Join(_projectRoot, ".cache/mtcli")

	if err := sh.RunWith(map[string]string{
		"CGO_ENABLED": "1",
		"CGO_CFLAGS":  "-DSQLITE_ENABLE_JSON1",
	}, mg.GoCmd(), "build", "-a", "-o", e2eBin, "./cmd/mtcli"); err != nil {
		return err
	}

	return sh.RunWith(map[string]string{
		"E2E_MTCLI_PATH": e2eBin,
	}, mg.GoCmd(), "test", "-count=1", "-race", "./integration...")
}

func (Test) Clean() error {
	var files []string

	if err := filepath.WalkDir(_projectRoot, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}

		matches, err := filepath.Match("index_tmp_*", filepath.Base(path))
		if err != nil {
			return err
		}

		if matches {
			files = append(files, path)
		}

		return nil
	}); err != nil {
		return err
	}

	var errCollector error

	for _, f := range files {
		multierr.AppendInto(&errCollector, fmt.Errorf("removing %s: %w", f, sh.Rm(f)))
	}

	return errCollector
}

type Build mg.Namespace

func (Build) CLI() error {
	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "1",
		"CGO_CFLAGS":  "-DSQLITE_ENABLE_JSON1",
	}, mg.GoCmd(), "build", "-a", "-o", filepath.Join(_projectRoot, "bin", "mtcli"), filepath.Join("cmd", "mtcli", "main.go"))
}

func (Build) CleanCLI() error {
	return sh.Rm(filepath.Join(_projectRoot, "bin", "mtcli"))
}

func (Build) Operator() error {
	return sh.RunWith(map[string]string{
		"CGO_ENABLED": "0",
	}, mg.GoCmd(), "build", "-a", "-o", filepath.Join(_projectRoot, "bin", "addon-metadata-operator"), filepath.Join("cmd", "addon-metadata-operator", "main.go"))
}

func (Build) CleanOperator() error {
	return sh.Rm(filepath.Join(_projectRoot, "bin", "addon-metadata-operator"))
}

func (Build) OperatorImage() error {
	runtime := runtime()
	if runtime == "" {
		return errors.New("could not find container runtime")
	}

	return sh.Run(runtime, "build", "-t", fmt.Sprintf("%s:%s", repository, _tag), "-f", "Dockerfile.build", _projectRoot)
}

type Generate mg.Namespace

func (Generate) Manifests(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdateControllerGen)

	return sh.Run(filepath.Join(_depBin, "controller-gen"),
		"crd", "rbac:roleName=manager-role",
		"webhook", `paths="./..."`,
		"output:crd:artifacts:config=config/crd/bases",
	)
}

func (Generate) Boilerplate(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdateControllerGen)

	return sh.Run(filepath.Join(_depBin, "controller-gen"),
		`object:headerFile="hack/boilerplate.go.txt"`,
		`paths="./..."`,
	)
}

type Release mg.Namespace

func (Release) PushOperatorImage() error {
	mg.Deps(
		Build.OperatorImage,
	)

	const creds_var = "DOCKER_CONF"

	creds, ok := os.LookupEnv(creds_var)
	if !ok {
		return fmt.Errorf("%q must be defined", creds_var)
	}

	if err := sh.Run("docker", "tag",
		fmt.Sprintf("%s:%s", repository, _tag),
		fmt.Sprintf("%s:latest", repository),
	); err != nil {
		return fmt.Errorf("tagging operator image: %w", err)
	}

	if err := sh.Run("docker",
		fmt.Sprintf("--config=%s", creds),
		"push", fmt.Sprintf("%s:%s", repository, _tag),
	); err != nil {
		return fmt.Errorf("pushing operator image commit tag: %w", err)
	}

	if err := sh.Run("docker",
		fmt.Sprintf("--config=%s", creds),
		"push", fmt.Sprintf("%s:%s", repository, _tag),
	); err != nil {
		return fmt.Errorf("pushing operator image latest tag: %w", err)
	}

	return nil
}

func (Release) CLI() error {
	mg.Deps(
		Release.Container,
	)

	runtime := runtime()
	if runtime == "" {
		return errors.New("could not find container runtime")
	}

	return sh.Run(runtime, "run", "--rm",
		"-e", "CGO_ENABLED=1",
		"-e", fmt.Sprintf("GITHUB_TOKEN=%s", os.Getenv("GITHUB_TOKEN")),
		fmt.Sprintf("amo-release:%s", _tag), "release",
	)
}

func (Release) Container() error {
	runtime := runtime()
	if runtime == "" {
		return errors.New("could not find container runtime")
	}

	return sh.Run(runtime, "build",
		"-t", fmt.Sprintf("amo-release:%s", _tag),
		"-f", "Dockerfile.release", _projectRoot,
	)
}

func runtime() string {
	prefferedRuntimes := []string{
		"podman",
		"docker",
	}

	for _, runtime := range prefferedRuntimes {
		runtimePath, err := exec.LookPath(runtime)
		if err == nil {
			return runtimePath
		}
	}

	return ""
}

type Deps mg.Namespace

func (Deps) UpdateControllerGen(ctx context.Context) error {
	return updateGoDependency(ctx, "sigs.k8s.io/controller-tools/cmd/controller-gen")
}

func (Deps) UpdateGolangCILint(ctx context.Context) error {
	return updateGoDependency(ctx, "github.com/golangci/golangci-lint/cmd/golangci-lint")
}

func (Deps) UpdateGoImports(ctx context.Context) error {
	return updateGoDependency(ctx, "golang.org/x/tools/cmd/goimports")
}

func (Deps) UpdateKind(ctx context.Context) error {
	return updateGoDependency(ctx, "sigs.k8s.io/kind")
}

func (Deps) UpdateKustomize(ctx context.Context) error {
	return updateGoDependency(ctx, "sigs.k8s.io/kustomize/kustomize/v4")
}

func updateGoDependency(ctx context.Context, src string) error {
	if err := setupDepsBin(); err != nil {
		return fmt.Errorf("creating dependencies bin directory: %w", err)
	}

	toolsDir := filepath.Join(_projectRoot, "tools")

	tidy := exec.CommandContext(ctx, "go", "mod", "tidy")
	tidy.Dir = toolsDir

	if err := tidy.Run(); err != nil {
		return fmt.Errorf("starting to tidy tools dir: %w", err)
	}

	install := exec.CommandContext(ctx, "go", "install", src)
	install.Dir = toolsDir
	install.Env = append(os.Environ(), fmt.Sprintf("GOBIN=%s", _depBin))

	if err := install.Run(); err != nil {
		return fmt.Errorf("starting to install command from source %q: %w", src, err)
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

		if err := downloadFile(ctx, urlPrefix+fmt.Sprintf("/v%s/pre-commit-%s.pyz", version, version), out); err != nil {
			return fmt.Errorf("downloading pre-commit: %w", err)
		}
	}

	return os.Chmod(out, 0775)
}

func downloadFile(ctx context.Context, url, out string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("constructing request: %w", err)
	}

	var client http.Client

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading file: %w", err)
	}

	defer res.Body.Close()

	if status := res.StatusCode; status != http.StatusOK {
		return fmt.Errorf("request failed with status %d", status)
	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("creating file %q: %w", out, err)
	}

	defer f.Close()

	if _, err := io.Copy(f, res.Body); err != nil {
		return fmt.Errorf("copying response: %w", err)
	}

	return nil
}

func setupDepsBin() error {
	return os.MkdirAll(_depBin, 0o774)
}

type Hooks mg.Namespace

func (Hooks) Enable(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	return sh.Run(filepath.Join(_depBin, "pre-commit"), "install",
		"--hook-type", "pre-commit",
		"--hook-type", "pre-push",
	)
}

func (Hooks) Disable(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	return sh.Run(filepath.Join(_depBin, "pre-commit"), "install")
}

func (Hooks) Run(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	return sh.Run(filepath.Join(_depBin, "pre-commit"), "run",
		"--show-diff-on-failure",
		"--from-ref", "origin/master", "--to-ref", "HEAD",
	)
}

func (Hooks) RunAllFiles(ctx context.Context) error {
	mg.CtxDeps(ctx, Deps.UpdatePreCommit)

	return sh.Run(filepath.Join(_depBin, "pre-commit"), "run", "--all-files")
}
