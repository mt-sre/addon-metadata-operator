package validate

import (
	"errors"
	"fmt"

	"github.com/spf13/pflag"
	"golang.org/x/mod/semver"
)

type options struct {
	Env                string
	Version            string
	Disabled           string
	Enabled            string
	ExcludedNamespaces []string
}

func (o *options) AddEnvFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&o.Env,
		"env",
		o.Env,
		"integration, stage or production",
	)
}

func (o *options) AddVersionFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&o.Version,
		"version",
		o.Version,
		"addon imageset version",
	)
}

func (o *options) AddDisabledFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&o.Disabled,
		"disabled",
		o.Disabled,
		"Disable specific validators, separated by ','. Can't be combined with --enabled.",
	)
}

func (o *options) AddEnabledFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&o.Enabled,
		"enabled",
		o.Enabled,
		"Enable specific validators, separated by ','. Can't be combined with --disabled.",
	)
}

func (o *options) AddExcludedNamespacesFlag(flags *pflag.FlagSet) {
	flags.StringArrayVar(
		&o.ExcludedNamespaces,
		"excluded-namespaces",
		o.ExcludedNamespaces,
		"Excludes the given namespaces from validation.",
	)
}

func (o *options) VerifyFlags() error {
	if !isValidEnv(o.Env) {
		return fmt.Errorf("'%s' is not a valid environment; must be one of 'integration', 'stage' or 'production'", o.Env)
	}

	// unset version is OK, will fallback to meta.addonImageSetVersion
	if o.Version == "" {
		return nil
	}
	// semver.IsValid(...) requires the following format vMAJOR.MINOR.PATCH
	// so we temporarily prefix the 'v' character
	if o.Version != "latest" && !semver.IsValid(fmt.Sprintf("v%v", o.Version)) {
		return fmt.Errorf("'%s' is not a valid version; must be one of 'latest' or match 'MAJOR.MINOR.PATCH'", o.Version)
	}

	if o.Disabled != "" && o.Enabled != "" {
		return errors.New("'--disabled' and '--enabled' are mutually exclusive options")
	}

	return nil
}

func isValidEnv(env string) bool {
	switch env {
	case "stage", "integration", "production":
		return true
	default:
		return false
	}
}
