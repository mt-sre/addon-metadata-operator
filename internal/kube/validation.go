package kube

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/validation"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
)

// AreValidk8sAnnotationNames validates the given names against the kubernetes
// format for annotations and returns a slice of validation failure messages
// if any issues are found. Otherwise an empty slice is returned.
func AreValidk8sAnnotationNames(names ...string) []string {
	return areValidk8sNames(IsValidk8sAnnotationName, names...)
}

// AreValidk8sLabelNames validates the given names against the kubernetes
// format for labels and returns a slice of validation failure messages
// if any issues are found. Otherwise an empty slice is returned.
func AreValidk8sLabelNames(names ...string) []string {
	return areValidk8sNames(IsValidk8sLabelName, names...)
}

// AreValidk8sNamespaceNames validates the given names against the kubernetes
// format for namespaces and returns a slice of validation failure messages
// if any issues are found. Otherwise an empty slice is returned.
func AreValidk8sNamespaceNames(names ...string) []string {
	return areValidk8sNames(IsValidk8sNamespaceName, names...)
}

func areValidk8sNames(validationFunc func(string) string, names ...string) []string {
	msgs := make([]string, 0, len(names))

	for _, name := range names {
		if msg := validationFunc(name); msg != "" {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

// IsValidk8sNamespaceName validates the given name against the kubernetes format for
// namespace names and returns a validation failure message if an issue is found.
// Otherwise an empty string is returned.
func IsValidk8sNamespaceName(name string) string {
	if valid, failureReasons := reasonsToResult(validation.ValidateNamespaceName(name, false)); !valid {
		return fmt.Sprintf("\"%s\" is not a valid kubernetes namespace name: %s", name, failureReasons)
	}

	return ""
}

// IsValidk8sSecretName validates the given name against the kubernetes format for
// secret names and returns a validation failure message if an issue is found.
// Otherwise an empty string is returned.
func IsValidk8sSecretName(name string) string {
	if valid, failureReasons := reasonsToResult(utilvalidation.IsDNS1123Subdomain(name)); !valid {
		return fmt.Sprintf("\"%s\" is not a valid kubernetes secret name: %s", name, failureReasons)
	}

	return ""
}

// IsValidk8sAnnotationName validates the given name against the kubernetes format for
// annotation names and returns a validation failure message if an issue is found.
// Otherwise an empty string is returned.
func IsValidk8sAnnotationName(name string) string {
	if msg := isQualifiedk8sName(strings.ToLower(name)); msg != "" {
		return fmt.Sprintf("\"%s\" is not a valid kubernetes annotation name: %s", name, msg)
	}

	return ""
}

// IsValidk8sLabelName validates the given name against the kubernetes format for
// label names and returns a validation failure message if an issue is found.
// Otherwise an empty string is returned.
func IsValidk8sLabelName(name string) string {
	if msg := isQualifiedk8sName(name); msg != "" {
		return fmt.Sprintf("\"%s\" is not a valid kubernetes label name: %s", name, msg)
	}

	return ""
}

func isQualifiedk8sName(name string) string {
	if valid, failureReasons := reasonsToResult(utilvalidation.IsQualifiedName(name)); !valid {
		return failureReasons
	}

	return ""
}

func reasonsToResult(reasons []string) (bool, string) {
	// Handles case where reasons = []string{""} which is
	// indistinguishable from []string{} when joined
	if len(reasons) > 0 {
		return false, strings.Join(reasons, ", ")
	}

	return true, ""
}
