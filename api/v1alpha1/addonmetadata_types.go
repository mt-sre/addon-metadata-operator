/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	mtsrev1 "github.com/mt-sre/addon-metadata-operator/pkg/mtsre/v1"
	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddonMetadataSpec defines the desired state of AddonMetadata
// View markers: $ controller-gen -www crd
// TODO add missing fields from schema, only required fields from jsonschema are present
type AddonMetadataSpec struct {
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,30}[A-Za-z0-9]$`
	// Unique ID of the addon
	ID string `json:"id" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[0-9A-Z\[\]][A-Za-z0-9-_ ()\[\]]+$`
	// Friendly name for the addon, displayed in the UI
	Name string `json:"name" validate:"required"`

	// +kubebuilder:validation:Required
	// Short description for the addon
	Description string `json:"description" validate:"required"`

	// +optional
	// +kubebuilder:validation:Pattern=`^http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+$`
	// Link to the addon documentation
	Link string `json:"link"`

	// +kubebuilder:validation:Required
	// Icon to be shown in UI. Should be around 200px and base64 encoded.
	Icon string `json:"icon" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^api\.openshift\.com/addon-[0-9a-z][0-9a-z-]{0,30}[0-9a-z]$`
	// Kubernetes label for the addon. Needs to match: 'api.openshift.com/<addon-id>'.
	Label string `json:"label" validate:"required"`

	// +kubebuilder:validation:Required
	// Set to true to allow installation of the addon.
	Enabled bool `json:"enabled" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([A-Za-z -]+ <[0-9A-Za-z_.-]+@redhat\.com>,?)+$`
	// Team or individual responsible for this addon. Needs to match: 'some name <some-email@redhat.com>'.
	AddonOwner string `json:"addonOwner" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^quay\.io/osd-addons/[a-z-]+$`
	// Quay repository for the addon operator. Needs to match: 'quay.io/osd-addons/<my-addon-repo>'.
	QuayRepo string `json:"quayRepo" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^quay\.io/[0-9A-Za-z._-]+/[0-9A-Za-z._-]+(:[A-Za-z0-9._-]+)?$`
	// Quay repository for the testHarness image. Needs to match: 'quay.io/<my-repo>/<my-test-harness>:<my-tag>'.
	TestHarness string `json:"testHarness" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum={AllNamespaces,OwnNamespace}
	// OLM InstallMode for the addon operator. One of: AllNamespaces or OwnNamespace.
	InstallMode string `json:"installMode" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,60}[A-Za-z0-9]$`
	// Namespace where the addon operator should be installed.
	TargetNamespace string `json:"targetNamespace" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:UniqueItems=true
	// Namespaces managed by the addon-operator. Need to include the TargetNamespace.
	Namespaces []string `json:"namespaces" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-_]{0,35}[A-Za-z0-9]$`
	// Refers to the SKU name for the addon.
	OcmQuotaName string `json:"ocmQuotaName" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=0
	// TODO: what is this?
	OcmQuotaCost int `json:"ocmQuotaCost" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^[A-Za-z0-9][A-Za-z0-9-]{0,30}[A-Za-z0-9]$`
	// Name of the addon operator.
	OperatorName string `json:"operatorName" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum={alpha,beta,stable,edge,rc}
	// OLM channel from which to install the addon-operator. One of: alpha, beta, stable, edge or rc.
	DefaultChannel string `json:"defaultChannel" validate:"required"`

	// +optional
	// Deprecated: List of channels where the addon operator is available.
	// Only needed for legacy addon builds.
	Channels *[]Channel `json:"channels"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default:{}
	// Labels to be applied on all listed namespaces.
	NamespaceLabels map[string]string `json:"namespaceLabels" validate:"required"`

	// +kubebuilder:validation:Required
	// +kubebuilder:default:{}
	// Annotations to be applied on all listed namespaces.
	NamespaceAnnotations map[string]string `json:"namespaceAnnotations" validate:"required"`

	// +optional
	// +kubebuilder:validation:Pattern=`^quay\.io/osd-addons/[a-z-]+`
	IndexImage *string `json:"indexImage"`

	// +optional
	// OCM representation of an add-on parameter
	AddOnParameters *[]ocmv1.AddOnParameter `json:"addOnParameters"`

	// +optional
	// OCM representation of an addon-requirement
	AddOnRequirements *[]ocmv1.AddOnRequirement `json:"addOnRequirements"`

	// +optional
	// OCM representation of an add-on sub operator. A sub operator is an
	// operator who's life cycle is controlled by the add-on umbrella operator.
	SubOperators *[]ocmv1.AddOnSubOperator `json:"subOperators"`

	// +optional
	// A string which specifies the imageset to use. Can either be 'latest' or a version string
	// MAJOR.MINOR.PATCH
	ImageSetVersion *string `json:"addonImageSetVersion"`

	// +optional
	HasExternalResources *bool `json:"hasExternalResources"`

	// +optional
	AddonNotifications *[]mtsrev1.Notification `json:"addonNotifications"`

	// +optional
	ManualInstallPlanApproval *bool `json:"manualInstallPlanApproval"`

	// +optional
	PullSecret string `json:"pullSecret"`

	// +optional
	// Labels to be applied to all objects created in the SelectorSyncSet.
	CommonLabels *map[string]string `json:"commonLabels"`

	// +optional
	// Annotations to be applied to all objects created in the SelectorSyncSet.
	CommonAnnotations *map[string]string `json:"commonAnnotations"`

	// +optional
	// Configuration parameters to be injected in the ServiceMonitor used for federation. The target prometheus server found by matchLabels needs to serve service-ca signed TLS traffic (https://docs.openshift.com/container-platform/4.6/security/certificate_types_descriptions/service-ca-certificates.html), and it needs to be runing inside the monitoring.namespace, with the service name 'prometheus'.
	Monitoring *mtsrev1.Monitoring `json:"monitoring"`

	// +optional
	// Deprecated: Replaced by SubscriptionConfig.
	BundleParameters *mtsrev1.BundleParameters `json:"bundleParameters"` //nolint: staticcheck // ignoring self-deprecation SA1019

	// +optional
	StartingCSV *string `json:"startingCSV"`

	// +optional
	PagerDuty *mtsrev1.PagerDuty `json:"pagerduty"`

	// +optional
	// Denotes the Deadmans Snitch Configuration which is supposed to be setup alongside the Addon.
	DeadmansSnitch *mtsrev1.DeadmansSnitch `json:"deadmanssnitch"`

	// +optional
	// Extra Resources to be applied to the Hive cluster.
	ExtraResources *[]string `json:"extraResources"`

	// +optional
	// Configs to be passed to the subscription OLM object.
	SubscriptionConfig *mtsrev1.SubscriptionConfig `json:"subscriptionConfig"`
}

// AddonMetadataStatus defines the observed state of AddonMetadata
type AddonMetadataStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
// AddonMetadata is the Schema for the AddonMetadata API
type AddonMetadata struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AddonMetadataSpec   `json:"spec,omitempty" validate:"required"`
	Status AddonMetadataStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// AddonMetadataList contains a list of AddonMetadata
type AddonMetadataList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AddonMetadata `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AddonMetadata{}, &AddonMetadataList{})
}

// *****
// Helper types
// *****

// Channel - list all channels for a given operator
type Channel struct {
	Name       string `json:"name"`
	CurrentCSV string `json:"currentCSV"`
}
