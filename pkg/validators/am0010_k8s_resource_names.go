package validators

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/internal/kube"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
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

func Validatek8sResourceAndFieldNames(mb types.MetaBundle) types.ValidatorResult {
	subValidators := []func(types.MetaBundle) []string{
		validateLabel,
		validateTargetNamespace,
		validateAllNamespaces,
		validateCommonAnnotations,
		validateCommonLabels,
		validatePagerDutyFields,
	}

	var msgs []string

	for _, v := range subValidators {
		msgs = append(msgs, v(mb)...)
	}

	if len(msgs) > 0 {
		return Fail(msgs...)
	}

	return Success()
}

func validateLabel(mb types.MetaBundle) (failures []string) {
	if msg := kube.IsValidk8sLabelName(mb.AddonMeta.Label); msg != "" {
		failures = append(failures, prefixedFailureMsg("label", msg))
	}

	return failures
}

func validateTargetNamespace(mb types.MetaBundle) (failures []string) {
	if msg := kube.IsValidk8sNamespaceName(mb.AddonMeta.TargetNamespace); msg != "" {
		failures = append(failures, prefixedFailureMsg("targetNamespace", msg))
	}

	return failures
}

func validateAllNamespaces(mb types.MetaBundle) []string {
	var result []string

	if msgs := kube.AreValidk8sNamespaceNames(mb.AddonMeta.Namespaces...); len(msgs) > 0 {
		for _, msg := range msgs {
			result = append(result, prefixedFailureMsg(
				"namespaces", msg,
			))
		}
	}

	if msgs := kube.AreValidk8sAnnotationNames(keys(mb.AddonMeta.NamespaceAnnotations)...); len(msgs) > 0 {
		for _, msg := range msgs {
			result = append(result, prefixedFailureMsg(
				"namespaceAnnotations", msg,
			))
		}
	}

	if msgs := kube.AreValidk8sLabelNames(keys(mb.AddonMeta.NamespaceLabels)...); len(msgs) > 0 {
		for _, msg := range msgs {
			result = append(result, prefixedFailureMsg(
				"namespaceLabels", msg,
			))
		}
	}

	return result
}

func validateCommonAnnotations(mb types.MetaBundle) []string {
	if mb.AddonMeta.CommonAnnotations == nil {
		return []string{}
	}

	annotations := *mb.AddonMeta.CommonAnnotations

	result := make([]string, 0, len(annotations))

	if msgs := kube.AreValidk8sAnnotationNames(keys(annotations)...); len(msgs) > 0 {
		for _, msg := range msgs {
			result = append(result, prefixedFailureMsg(
				"commonAnnotations", msg,
			))
		}
	}

	return result
}

func validateCommonLabels(mb types.MetaBundle) []string {
	if mb.AddonMeta.CommonLabels == nil {
		return []string{}
	}

	labels := *mb.AddonMeta.CommonLabels

	result := make([]string, 0, len(labels))

	if msgs := kube.AreValidk8sLabelNames(keys(labels)...); len(msgs) > 0 {
		for _, msg := range msgs {
			result = append(result, prefixedFailureMsg(
				"commonLabels", msg,
			))
		}
	}

	return result
}

func validatePagerDutyFields(mb types.MetaBundle) []string {
	pd := mb.AddonMeta.PagerDuty
	if pd == nil {
		return []string{}
	}

	var msgs []string

	if msg := kube.IsValidk8sNamespaceName(pd.SecretNamespace); msg != "" {
		msgs = append(msgs, msg)
	}

	if msg := kube.IsValidk8sSecretName(pd.SecretName); msg != "" {
		msgs = append(msgs, msg)
	}

	return msgs
}

func keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func prefixedFailureMsg(prefix string, msg string) string {
	return string(fmt.Sprintf("%s: %s", prefix, msg))
}
