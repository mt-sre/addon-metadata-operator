package validators_test

import (
	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	v1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validators"
)

func init() {
	TestRegistry.Add(TestAM0010{})
}

type TestAM0010 struct{}

func (val TestAM0010) Name() string {
	return validators.AM0010.Name
}

func (val TestAM0010) Run(mb types.MetaBundle) types.ValidatorResult {
	return validators.AM0010.Runner(mb)
}

func (val TestAM0010) SucceedingCandidates() ([]types.MetaBundle, error) {
	candidates, err := testutils.DefaultSucceedingCandidates()
	if err != nil {
		return nil, err
	}

	return append(candidates,
		types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				Label:           "api.openshift.com/foo-addon",
				TargetNamespace: "redhat-foo-operator",
				Namespaces: []string{
					"redhat-foo-operator",
				},
				NamespaceAnnotations: map[string]string{
					"Foo-bar.com/Foo-namespace_annotation.1":  "true",
					"Foo.bar0.com/Foo-namespace_annotation.2": "true",
					"Foo.bar.com/Foo-namespace_annotation.3":  "true",
				},
				NamespaceLabels: map[string]string{
					"foo-bar.com/Foo-namespace_label.1":  "true",
					"foo.bar0.com/Foo-namespace_label.2": "true",
					"foo.bar.com/Foo-namespace_label.3":  "true",
				},
				CommonAnnotations: &map[string]string{
					"foo-bar.com/Foo-Common_Annotation.1":  "true",
					"foo.bar0.com/Foo-Common_Annotation.2": "true",
					"foo.bar.com/Foo-Common_Annotation.3":  "true",
				},
				CommonLabels: &map[string]string{
					"foo-bar.com/Foo-common_label.1":  "true",
					"foo.bar0.com/Foo-common_label.2": "true",
					"foo.bar.com/Foo-common_label.3":  "true",
				},
				PagerDuty: &v1.PagerDuty{
					SecretName:      "pagerduty-secret",
					SecretNamespace: "redhat-foo-operator",
				},
			},
		}), nil
}

func (val TestAM0010) FailingCandidates() ([]types.MetaBundle, error) {
	var failingBundles []types.MetaBundle

	failingBundles = append(failingBundles, metaBundlesWithBadNamespaces()...)
	failingBundles = append(failingBundles, metaBundlesWithBadLabels()...)
	failingBundles = append(failingBundles, metaBundlesWithBadAnnotations()...)
	failingBundles = append(failingBundles, metaBundlesWithBadSecrets()...)

	return failingBundles, nil
}

func metaBundlesWithBadNamespaces() []types.MetaBundle {
	badNamespaces := []string{
		// Namespace name is empty
		"",
		// Namespace name exceeds 63 characters
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		// Namespace name begins with non-alphanumeric character
		"-foo-namespace",
		// Namespace name ends with non-alphanumeric charatcter
		"foo-namespace-",
		// Namespace name includes upper case alpha characters
		"Foo-Namespace",
		// Namespace name includes unallowed character '_'
		"foo_namespace",
	}

	bundles := make([]types.MetaBundle, 0, len(badNamespaces)*3)

	for _, ns := range badNamespaces {
		bundles = append(bundles,
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: ns,
					Label:           "foo-label",
				},
			},
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					Namespaces: []string{
						ns,
					},
				},
			},
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					PagerDuty: &v1.PagerDuty{
						SecretNamespace: ns,
					},
				},
			},
		)
	}

	return bundles
}

func metaBundlesWithBadLabels() []types.MetaBundle {
	badLabels := []string{
		// Label is empty
		"",
		// Empty Prefix
		"/label",
		// Empty Name
		"foo.bar.com/",
		// Multiple prefix seperators
		"foo.com/bar/label",
		// Prefix contains unallowed character
		"foo*com/label",
		// Prefix contains uppercase characters
		// Note: this is the primary difference between label and annotation
		// validations as annotations permit uppercase prefixes
		"FOO.com/label",
		// Prefix does not start with alphanumeric
		".foo.com/label",
		// Prefix does not end with alphanumeric
		"foo.com./label",
		// Name does not start with alphanumeric
		"foo.com/-label",
		// Name does not end with alphanumeric
		"foo.com/label-",
		// Prefix exceeds 253 characters
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaa/label",
		// Name exceeds 63 characters
		"foo.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	bundles := make([]types.MetaBundle, 0, len(badLabels)*3)

	for _, label := range badLabels {
		bundles = append(bundles,
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           label,
				},
			},
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					NamespaceLabels: map[string]string{
						label: "",
					},
				},
			},
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					CommonLabels: &map[string]string{
						label: "",
					},
				},
			},
		)
	}

	return bundles
}

func metaBundlesWithBadAnnotations() []types.MetaBundle {
	badAnnotations := []string{
		// Annotation is empty
		"",
		// Empty Prefix
		"/annotation",
		// Empty Name
		"foo.bar.com/",
		// Multiple prefix seperators
		"foo.com/bar/annotation",
		// Prefix contains unallowed character
		"foo*com/annotation",
		// Prefix does not start with alphanumeric
		".foo.com/annotation",
		// Prefix does not end with alphanumeric
		"foo.com./annotation",
		// Name does not start with alphanumeric
		"foo.com/-annotation",
		// Name does not end with alphanumeric
		"foo.com/annotation-",
		// Prefix exceeds 253 characters
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaa/label",
		// Name exceeds 63 characters
		"foo.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	bundles := make([]types.MetaBundle, 0, len(badAnnotations)*2)

	for _, annotation := range badAnnotations {
		bundles = append(bundles,
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					NamespaceAnnotations: map[string]string{
						annotation: "",
					},
				},
			},
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					CommonAnnotations: &map[string]string{
						annotation: "",
					},
				},
			},
		)
	}

	return bundles
}

func metaBundlesWithBadSecrets() []types.MetaBundle {
	badSecrets := []string{
		// Secret name is empty
		"",
		// Secret name contains unallowed character
		"foo*com",
		// Secret name does not start with alphanumeric
		".foo.com",
		// Secret name does not end with alphanumeric
		"foo.com.",
		// Secret name contains an uppercase character
		"Foo.com",
		// Name exceeds 253 characters
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	bundles := make([]types.MetaBundle, 0, len(badSecrets))

	for _, secret := range badSecrets {
		bundles = append(bundles,
			types.MetaBundle{
				AddonMeta: &v1alpha1.AddonMetadataSpec{
					TargetNamespace: "foo-namespace",
					Label:           "foo-label",
					PagerDuty: &v1.PagerDuty{
						SecretName: secret,
					},
				},
			},
		)
	}

	return bundles
}
