package cmd

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/mt-sre/addon-metadata-operator/pkg/validate"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	_ "github.com/mt-sre/addon-metadata-operator/pkg/validator/register"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

func init() {
	validateCmd.Flags().StringVar(&validateEnv, "env", validateEnv, "integration, stage or production")
	validateCmd.Flags().StringVar(&validateVersion, "version", validateVersion, "addon imageset version")
	validateCmd.Flags().StringVar(&validateDisabled, "disabled", validateDisabled, "Disable specific validators, separated by ','. Can't be combined with --enabled.")
	validateCmd.Flags().StringVar(&validateEnabled, "enabled", validateEnabled, "Enable specific validators, separated by ','. Can't be combined with --disabled.")
	mtcli.AddCommand(validateCmd)
}

var (
	validateEnv      = "stage"
	validateVersion  = ""
	validateDisabled = ""
	validateEnabled  = ""
	validateExamples = []string{
		"  # Validate an addon in staging. Uses the latest version if it supports imageset.",
		"  mtcli validate --env stage --version latest internal/testdata/addons-imageset/reference-addon",
		"  # Validate a version 1.0.0 of a production addon using imageset.",
		"  mtcli validate --env production --version 1.0.0 <path/to/addon_dir>",
		"  # Validate a staging addon that is not using imageset, but a static indexImage.",
		"  mtcli validate --env stage <path/to/addon_dir>",
		"  # Validate an integration addon using imageset, disabling validators 001_foo and 002_bar.",
		"  mtcli validate --env integration --disabled AM0001,AM0002 <path/to/addon_dir>",
		"  # Validate an integration addon using imageset, enabled only 001_foo.",
		"  mtcli validate --env integration --enabled AM0001 <path/to/addon_dir>",
	}
	validateLong = "Validate an addon metadata and it's bundles against custom validators."
	validateCmd  = &cobra.Command{
		Use:     "validate",
		Short:   "Validate addon metadata, bundles and imagesets.",
		Long:    validateLong,
		Example: strings.Join(validateExamples, "\n"),
		Args:    cobra.ExactArgs(1),
		Run:     validateMain,
	}
)

func validateMain(cmd *cobra.Command, args []string) {
	addonDir, err := parseAddonDir(args[0])
	if err != nil {
		fail(1, "unable to parse the provided directory '%s': %v", args[0], err)
	}

	if err := verifyArgsAndFlags(addonDir); err != nil {
		fail(1, "unable to process flag or argument: %v", err)
	}

	meta, err := utils.NewMetaLoader(addonDir, validateEnv, validateVersion).Load()
	if err != nil {
		fail(1, "unable to load addon metadata from file '%v': %v", addonDir, err)
	}

	bundles, err := utils.ExtractAndParseAddons(*meta.IndexImage, meta.PackageName)
	if err != nil {
		fail(1, "unable to extract and parse bundles from the given index image: %v", err)
	}

	if validateDisabled != "" && validateEnabled != "" {
		fail(1, "'--disabled' and '--enabled' are mutually exclusive options")
	}

	filter, err := generateFilter(validateDisabled, validateEnabled)
	if err != nil {
		fail(1, "could not filter validators: %v", err)
	}

	runner, err := validator.NewRunner(
		validator.WithMiddleware{
			validator.NewRetryMiddleware(),
		},
	)
	if err != nil {
		fail(1, "unable to initialize validator: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mb := types.MetaBundle{
		AddonMeta: meta,
		Bundles:   bundles,
	}

	var results validator.ResultList

	for res := range runner.Run(ctx, mb, filter) {
		results = append(results, res)
	}

	sort.Sort(results)

	table := validate.NewResultTable()
	for _, res := range results {
		table.WriteResult(res)
	}

	fmt.Printf("%v\n\n", table.String())
	fmt.Println("Please consult corresponding validator wikis: https://github.com/mt-sre/addon-metadata-operator/wiki/<code>.")

	if err := runner.CleanUp(); err != nil {
		fail(1, "unable to release resources: %s", err)
	}

	if errs := results.Errors(); len(errs) > 0 {
		utils.PrintValidationErrors(errs)
		os.Exit(1)
	}

	if results.HasFailure() {
		os.Exit(1)
	}
}

func parseAddonDir(dir string) (string, error) {
	if !path.IsAbs(dir) {
		return filepath.Abs(dir)
	}
	return dir, nil
}

func fail(code int, msg string, args ...interface{}) {
	fmt.Printf("A fatal error occurred while preparing validations: "+msg+"\n", args...)

	os.Exit(code)
}

func verifyArgsAndFlags(addonDir string) error {
	if err := verifyAddonDir(addonDir); err != nil {
		return err
	}
	if err := verifyEnv(validateEnv); err != nil {
		return err
	}
	return verifyVersion(validateVersion)
}

// addonDir is an absolute path at this point
func verifyAddonDir(addonDir string) error {
	dir, err := os.Stat(addonDir)
	if err != nil {
		return fmt.Errorf("error while reading directory: %w", err)
	}
	if !dir.IsDir() {
		return fmt.Errorf("'%s' is not a directory", addonDir)
	}
	return nil
}

func verifyEnv(env string) error {
	if env != "integration" && env != "stage" && env != "production" {
		return fmt.Errorf("'%s' is not a valid environment; must be one of 'integration', 'stage' or 'production'", env)
	}
	return nil
}

func verifyVersion(version string) error {
	// unset version is OK, will fallback to meta.addonImageSetVersion
	if version == "" {
		return nil
	}
	// semver.IsValid(...) requires the following format vMAJOR.MINOR.PATCH
	// so we temporarily prefix the 'v' character
	if version != "latest" && !semver.IsValid(fmt.Sprintf("v%v", version)) {
		return fmt.Errorf("'%s' is not a valid version; must be one of 'latest' or match 'MAJOR.MINOR.PATCH'", version)
	}
	return nil
}

func generateFilter(disabled, enabled string) (validator.Filter, error) {
	if disabled == "" && enabled == "" {
		return nil, nil
	}

	if disabled != "" {
		codes, err := parseCodeList(disabled)
		if err != nil {
			return nil, fmt.Errorf("unable to process '--disabled' option argument: %w", err)
		}

		return validator.Not(validator.MatchesCodes(codes...)), nil
	}

	codes, err := parseCodeList(enabled)
	if err != nil {
		return nil, fmt.Errorf("unable to process '--enabled' option argument: %w", err)
	}

	return validator.MatchesCodes(codes...), nil
}

func parseCodeList(maybeList string) ([]validator.Code, error) {
	rawStrings := strings.Split(maybeList, ",")

	res := make([]validator.Code, 0, len(rawStrings))

	for _, s := range rawStrings {
		c, err := validator.ParseCode(s)
		if err != nil {
			return nil, fmt.Errorf("invalid code list '%s': %w", maybeList, err)
		}

		res = append(res, c)
	}

	return res, nil
}
