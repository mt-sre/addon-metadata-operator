package validators

import (
	"fmt"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"k8s.io/apimachinery/pkg/api/validation"
	utilvalidation "k8s.io/apimachinery/pkg/util/validation"
)

var AM0010 = types.Validator{
	Code:        "AM0010",
	Name:        "k8s_resource_and_field_names",
	Description: "Validates k8s namespaces, labels, and annotations within Addon metadata against k8s standards",
	Runner:      Validatek8sResourceAndFieldNames,
}

func init() {
	Registry.Add(AM0010)
}

const (
	emptyMsg = failureMsg("")
)

type failureMsg string

func (m failureMsg) IsEmpty() bool {
	return m == ""
}

func prefixedFailureMsg(prefix string, msg failureMsg) failureMsg {
	return failureMsg(fmt.Sprintf("%s: %s", prefix, msg))
}

func joinFailureMsgs(msgs ...failureMsg) failureMsg {
	msgStrings := make([]string, 0, len(msgs))

	for _, m := range msgs {
		msgStrings = append(msgStrings, string(m))
	}

	return failureMsg(strings.Join(msgStrings, ", "))
}

func Validatek8sResourceAndFieldNames(mb types.MetaBundle) types.ValidatorResult {
	subValidators := []subValidator{
		validateLabel,
		validateTargetNamespace,
		validateAllNamespaces,
		validateCommonAnnotations,
		validateCommonLabels,
	}

	var msgs []failureMsg

	for _, v := range subValidators {
		msgs = append(msgs, v(mb)...)
	}

	if len(msgs) > 0 {
		return Fail(string(joinFailureMsgs(msgs...)))
	}

	return Success()
}

type subValidator func(types.MetaBundle) []failureMsg

func validateLabel(mb types.MetaBundle) (failures []failureMsg) {
	if msg := isValidLabelName(mb.AddonMeta.Label); !msg.IsEmpty() {
		failures = append(failures, prefixedFailureMsg("label", msg))
	}

	return failures
}

func validateTargetNamespace(mb types.MetaBundle) (failures []failureMsg) {
	if msg := isValidNamespaceName(mb.AddonMeta.TargetNamespace); !msg.IsEmpty() {
		failures = append(failures, prefixedFailureMsg("targetNamespace", msg))
	}

	return failures
}

func validateAllNamespaces(mb types.MetaBundle) []failureMsg {
	var result []failureMsg

	if msgs := areValidNamespaceNames(mb.AddonMeta.Namespaces...); len(msgs) > 0 {
		result = append(result, prefixedFailureMsg(
			"namespaces", joinFailureMsgs(msgs...),
		))
	}

	if msgs := areValidAnnotationNames(keys(mb.AddonMeta.NamespaceAnnotations)...); len(msgs) > 0 {
		result = append(result, prefixedFailureMsg(
			"namespaceAnnotations", joinFailureMsgs(msgs...),
		))
	}

	if msgs := areValidLabelNames(keys(mb.AddonMeta.NamespaceLabels)...); len(msgs) > 0 {
		result = append(result, prefixedFailureMsg(
			"namespaceLabels", joinFailureMsgs(msgs...),
		))
	}

	return result
}

func validateCommonAnnotations(mb types.MetaBundle) []failureMsg {
	if mb.AddonMeta.CommonAnnotations == nil {
		return []failureMsg{}
	}

	annotations := *mb.AddonMeta.CommonAnnotations

	result := make([]failureMsg, 0, len(annotations))

	if msgs := areValidAnnotationNames(keys(annotations)...); len(msgs) > 0 {
		result = append(result, prefixedFailureMsg(
			"commonAnnotations", joinFailureMsgs(msgs...),
		))
	}

	return result
}

func validateCommonLabels(mb types.MetaBundle) []failureMsg {
	if mb.AddonMeta.CommonLabels == nil {
		return []failureMsg{}
	}

	labels := *mb.AddonMeta.CommonLabels

	result := make([]failureMsg, 0, len(labels))

	if msgs := areValidLabelNames(keys(labels)...); len(msgs) > 0 {
		result = append(result, prefixedFailureMsg(
			"commonLabels", joinFailureMsgs(msgs...),
		))
	}

	return result
}

type k8sName string

const (
	k8sNameAnnotation = k8sName("annotation")
	k8sNameLabel      = k8sName("label")
	k8sNameNamespace  = k8sName("namespace")
)

func areValidAnnotationNames(names ...string) []failureMsg {
	return areValidNames(k8sNameAnnotation, names...)
}

func areValidLabelNames(names ...string) []failureMsg {
	return areValidNames(k8sNameLabel, names...)
}

func areValidNamespaceNames(names ...string) []failureMsg {
	return areValidNames(k8sNameNamespace, names...)
}

func areValidNames(nameType k8sName, names ...string) []failureMsg {
	validationFunc := k8sNameToValidationFunc(nameType)

	msgs := make([]failureMsg, 0, len(names))

	for _, name := range names {
		if msg := validationFunc(name); !msg.IsEmpty() {
			msgs = append(msgs, msg)
		}
	}

	return msgs
}

func k8sNameToValidationFunc(nameType k8sName) func(string) failureMsg {
	switch nameType {
	case k8sNameAnnotation:
		return isValidAnnotationName
	case k8sNameLabel:
		return isValidLabelName
	case k8sNameNamespace:
		return isValidNamespaceName
	default:
		// panic is preffered here since this indicates a programming error
		panic(fmt.Sprintf("no validation function defined for '%v'", nameType))
	}
}

func isValidNamespaceName(name string) failureMsg {
	if valid, failureReasons := reasonsToResult(validation.ValidateNamespaceName(name, false)); !valid {
		return prefixedFailureMsg(
			fmt.Sprintf("\"%s\" is not a valid kubernetes namespace name", name),
			failureMsg(failureReasons),
		)
	}

	return emptyMsg
}

func isValidAnnotationName(name string) failureMsg {
	return isQualifiedName(k8sNameAnnotation, name)
}

func isValidLabelName(name string) failureMsg {
	return isQualifiedName(k8sNameLabel, name)
}

func isQualifiedName(nameType k8sName, name string) failureMsg {
	var adjustedName string

	switch nameType {
	case k8sNameAnnotation:
		adjustedName = strings.ToLower(name)
	default:
		adjustedName = name
	}

	if valid, failureReasons := reasonsToResult(utilvalidation.IsQualifiedName(adjustedName)); !valid {
		return prefixedFailureMsg(
			fmt.Sprintf("\"%s\" is not a valid kubernetes %s name", name, nameType),
			failureMsg(failureReasons),
		)
	}

	return emptyMsg
}

func keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func reasonsToResult(reasons []string) (bool, string) {
	// Handles case where reasons = []string{""} which is
	// indistinguishable from []string{} when joined
	if len(reasons) > 0 {
		return false, strings.Join(reasons, ", ")
	}

	return true, ""
}
