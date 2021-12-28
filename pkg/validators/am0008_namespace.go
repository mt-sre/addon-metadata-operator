package validators

import (
	"fmt"
	"regexp"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var namespaceRegexExceptions = []string{
	"acm",
	"codeready-workspaces-operator",
	"codeready-workspaces-operator-qe",
	"addon-dba-operator",
	"openshift-storage",
	"prow",
}

var namespacePresenceExceptions = []string{
	"openshift-logging",
}

var namespaceRegex = regexp.MustCompile(`^redhat-.*$`)

var AM0008 = utils.Validator{
	Code:        "AM0008",
	Name:        "ensure_namespace",
	Description: "Ensure that the target namespace is listed in the set of channels listed",
	Runner:      ValidateNamespace,
}

func init() {
	Registry.Add(AM0008)
}

func ValidateNamespace(metabundle utils.MetaBundle) (bool, string, error) {
	targetNamespace := metabundle.AddonMeta.TargetNamespace
	namespaceList := metabundle.AddonMeta.Namespaces
	valid := validateNamespacePresence(targetNamespace, namespaceList, namespacePresenceExceptions)
	if !valid {
		return false, "Target namespace is not in the list of supplied namespaces", nil
	}

	allValid, failedNamespaces := validateNamespaceRegex(namespaceList, namespaceRegexExceptions)
	if !allValid {
		return false, fmt.Sprintf("Some namespaces doesn't start with 'redhat-*' %v", failedNamespaces), nil
	}
	return true, "", nil
}

func validateNamespaceRegex(namespaceList []string, exceptionList []string) (bool, []string) {
	failedNamespaces := []string{}
	for _, namespace := range namespaceList {
		if includes(namespace, exceptionList) {
			continue
		}
		matched := namespaceRegex.MatchString(namespace)
		if !matched {
			failedNamespaces = append(failedNamespaces, namespace)
		}
	}
	return len(failedNamespaces) == 0, failedNamespaces
}

func validateNamespacePresence(targetNamespace string, namespaceList []string, exceptionList []string) bool {
	// Return true if targetNamespace is in the exceptionList or
	// Return true if targetNamespace is in the namespaceList
	return includes(targetNamespace, exceptionList) || includes(targetNamespace, namespaceList)
}

func includes(item string, itemList []string) bool {
	for _, elem := range itemList {
		if elem == item {
			return true
		}
	}
	return false
}
