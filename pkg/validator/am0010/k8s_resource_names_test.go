package am0010

import (
	"testing"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	v1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator/testutils"
	"github.com/stretchr/testify/require"
)

func TestK8sResourcesAndFieldNamesValid(t *testing.T) {
	t.Parallel()

	bundles, err := testutils.DefaultValidBundleMap()
	require.NoError(t, err)

	for name, bundle := range map[string]types.MetaBundle{
		"all valid namespaces, annotations, and labels": {
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
		},
	} {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewK8SResourceAndFieldNames)
	tester.TestValidBundles(bundles)
}

func TestK8sResourcesAndFieldNamesInvalid(t *testing.T) {
	t.Parallel()

	bundles := make(map[string]types.MetaBundle)

	for name, bundle := range metaBundlesWithBadNamespaces() {
		bundles[name] = bundle
	}

	for name, bundle := range metaBundlesWithBadLabels() {
		bundles[name] = bundle
	}

	for name, bundle := range metaBundlesWithBadAnnotations() {
		bundles[name] = bundle
	}

	for name, bundle := range metaBundlesWithBadSecrets() {
		bundles[name] = bundle
	}

	tester := testutils.NewValidatorTester(t, NewK8SResourceAndFieldNames)
	tester.TestInvalidBundles(bundles)
}

func metaBundlesWithBadNamespaces() map[string]types.MetaBundle {
	badNamespaces := map[string]string{
		"namespace exceeds 63 characters":              "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"namespace begins with non-alphanumeric":       "-foo-namespace",
		"namespace ends with non-alphanumeric":         "foo-namespace-",
		"namespace name includes upper-case":           "Foo-Namespace",
		"namespace name includes disallowed character": "foo_namespace",
	}

	bundles := make(map[string]types.MetaBundle)

	for name, ns := range badNamespaces {
		bundles["targetNamespace: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: ns,
				Label:           "foo-label",
			},
		}

		bundles["Namespaces: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           "foo-label",
				Namespaces: []string{
					ns,
				},
			},
		}
	}

	return bundles
}

func metaBundlesWithBadLabels() map[string]types.MetaBundle {
	badLabels := map[string]string{
		"empty prefix":                        "/label",
		"empty name":                          "foo.bar.com/",
		"multiple prefix separators":          "foo.com/bar/label",
		"prefix contains unallowed character": "foo*com/label",
		// Note: this is the primary difference between label and annotation
		// validations as annotations permit uppercase prefixes
		"prefix contains uppercase":               "FOO.com/label",
		"prefix does not begin with alphanumeric": ".foo.com/label",
		"prefix does not end with alphanumeric":   "foo.com./label",
		"name does not start with alphanumeric":   "foo.com/-label",
		"name does not end with alphanumeric":     "foo.com/label-",
		"prefix exceeds 253 characters": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaa/label",
		"name exceeds 63 characters": "foo.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	bundles := make(map[string]types.MetaBundle)

	for name, label := range badLabels {
		bundles["Label: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           label,
			},
		}

		bundles["NamespaceLabels: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           "foo-label",
				NamespaceLabels: map[string]string{
					label: "",
				},
			},
		}

		bundles["CommonLabels: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           "foo-label",
				CommonLabels: &map[string]string{
					label: "",
				},
			},
		}
	}

	return bundles
}

func metaBundlesWithBadAnnotations() map[string]types.MetaBundle {
	badAnnotations := map[string]string{
		"empty prefix":                            "/annotation",
		"empty name":                              "foo.bar.com/",
		"multiple prefix separators":              "foo.com/bar/annotation",
		"prefix contains disallowed character":    "foo*com/annotation",
		"prefix does not start with alphanumeric": ".foo.com/annotation",
		"prefix does not end with alphanumeric":   "foo.com./annotation",
		"name does not start with alphanumeric":   "foo.com/-annotation",
		"name does not end with alphanumeric":     "foo.com/annotation-",
		"prefix exceeds 253 characters": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaa/label",
		"name exceeds 63 characters": "foo.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	bundles := make(map[string]types.MetaBundle)

	for name, annotation := range badAnnotations {
		bundles["NamespaceAnnotations: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           "foo-label",
				NamespaceAnnotations: map[string]string{
					annotation: "",
				},
			},
		}

		bundles["CommonAnnotations: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           "foo-label",
				CommonAnnotations: &map[string]string{
					annotation: "",
				},
			},
		}
	}

	return bundles
}

func metaBundlesWithBadSecrets() map[string]types.MetaBundle {
	badSecrets := map[string]string{
		"empty secret name":                       "",
		"secret contains disallowed character":    "foo*com",
		"secret does not start with alphanumeric": ".foo.com",
		"secret does not end with alphanumeric":   "foo.com.",
		"secret contains uppercase":               "Foo.com",
		"name exceeds 253 characters": "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" +
			"aaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}

	bundles := make(map[string]types.MetaBundle)

	for name, secret := range badSecrets {
		bundles["PagerDuty.SecretName: "+name] = types.MetaBundle{
			AddonMeta: &v1alpha1.AddonMetadataSpec{
				TargetNamespace: "foo-namespace",
				Label:           "foo-label",
				PagerDuty: &v1.PagerDuty{
					SecretName: secret,
				},
			},
		}
	}

	return bundles
}
