package objects

import (
	"fmt"

	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	hivev1 "github.com/openshift/hive/apis/hive/v1"

	// pdv1alpha1 "github.com/openshift/pagerduty-operator/pkg/apis/pagerduty/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// NewDeploySelectorSyncSet - creates a deploy SelectorSyncSet from the addonMetadata
func NewDeploySelectorSyncSet(addonMetadata *addonsv1alpha1.AddonMetadata) (*hivev1.SelectorSyncSet, error) {
	sss := &hivev1.SelectorSyncSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("addon-%s", addonMetadata.Spec.ID),
		},
		Spec: hivev1.SelectorSyncSetSpec{
			ClusterDeploymentSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					addonMetadata.Spec.Label: "true",
				},
			},
			SyncSetCommonSpec: hivev1.SyncSetCommonSpec{
				// TODO: chgange to Upsert? (jgwosdz)
				ResourceApplyMode: hivev1.SyncResourceApplyMode,
				// TODO - start from line 118 in template
				//	- Addon CR
				//  - Secrets
				//  - ExtraResources
				Resources: []runtime.RawExtension{},
			},
		},
	}
	return sss, nil
}

// NewDeleteSelectorSyncSet - creates a delete SelectorSyncSet from the addonMetadata
func NewDeleteSelectorSyncSet(addonMetadata *addonsv1alpha1.AddonMetadata) (*hivev1.SelectorSyncSet, error) {
	return nil, nil
}

// NewPagerDutyIntegration - creates a PagerDutyIntegration from addonMetadata
// func NewPagerDutyIntegration(addonMetadata *addonsv1alpha1.AddonMetadata) (*pdv1alpha1.PagerDutyIntegration, error) {
// 	return nil, nil
// }
