package am0008

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	"golang.org/x/exp/slices"
)

func init() {
	validator.Register(NewNamespace)
}

const (
	code = 8
	name = "ensure_namespace"
	desc = "Ensure that the target namespace is listed in the set of channels listed"
)

func NewNamespace(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &Namespace{
		Base:               base,
		ExcludedNamespaces: deps.ValidatorConfig.ExcludedNamespaces,
	}, nil
}

type Namespace struct {
	*validator.Base
	ExcludedNamespaces []string
}

func (n *Namespace) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
	targetNamespace := mb.AddonMeta.TargetNamespace
	namespaceList := mb.AddonMeta.Namespaces
	valid := validateNamespacePresence(targetNamespace, namespaceList, namespacePresenceExceptions())
	if !valid {
		return n.Fail("Target namespace is not in the list of supplied namespaces")
	}

	allValid, failedNamespaces := validateNamespaceRegex(namespaceList, n.ExcludedNamespaces)
	if !allValid {
		return n.Fail(fmt.Sprintf("Some namespaces doesn't start with 'redhat-*' %v", failedNamespaces))
	}
	return n.Success()
}

func validateNamespacePresence(targetNamespace string, namespaceList []string, exceptionList []string) bool {
	// Return true if targetNamespace is in the exceptionList or
	// Return true if targetNamespace is in the namespaceList
	return slices.Contains(exceptionList, targetNamespace) || slices.Contains(namespaceList, targetNamespace)
}

func namespacePresenceExceptions() []string {
	return []string{
		"openshift-logging",
	}
}

func validateNamespaceRegex(namespaceList []string, exceptionList []string) (bool, []string) {
	namespaceRegex := regexp.MustCompile(`^redhat-.*$`)

	var failedNamespaces []string

	for _, namespace := range namespaceList {
		if slices.Contains(exceptionList, namespace) {
			continue
		}

		if namespaceRegex.MatchString(namespace) {
			continue
		}

		failedNamespaces = append(failedNamespaces, namespace)
	}

	return len(failedNamespaces) == 0, failedNamespaces
}
