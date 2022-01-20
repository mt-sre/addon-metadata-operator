package am0010

import (
	"context"
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/internal/kube"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
)

func init() {
	validator.Register(NewK8SResourceAndFieldNames)
}

const (
	code = 10
	name = "k8s_resource_and_field_names"
	desc = "Validates k8s namespaces, labels, and annotations within Addon metadata against k8s standards"
)

func NewK8SResourceAndFieldNames(deps validator.Dependencies) (validator.Validator, error) {
	base, err := validator.NewBase(
		code,
		validator.BaseName(name),
		validator.BaseDesc(desc),
	)
	if err != nil {
		return nil, err
	}

	return &K8SResourceAndFieldNames{
		Base: base,
	}, nil
}

type K8SResourceAndFieldNames struct {
	*validator.Base
}

func (k *K8SResourceAndFieldNames) Run(ctx context.Context, mb types.MetaBundle) validator.Result {
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
		return k.Fail(msgs...)
	}

	return k.Success()
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
