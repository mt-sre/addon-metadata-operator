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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
)

// AddonMetadataReconciler reconciles a AddonMetadata object
type AddonMetadataReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=addonsflow.redhat.openshift.io,resources=AddonMetadata,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=addonsflow.redhat.openshift.io,resources=AddonMetadata/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=addonsflow.redhat.openshift.io,resources=AddonMetadata/finalizers,verbs=update

func (r *AddonMetadataReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&addonsv1alpha1.AddonMetadata{}).
		Complete(r)
}

// Reconcile - Implements the AddonMetadata flow
// 1. Find corresponding imageset
// 2. Validate the workload (extract bundles and index image)
// 3. POST payloads to OCM
func (r *AddonMetadataReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("addonmetadata", req.NamespacedName.String())

	meta := &addonsv1alpha1.AddonMetadata{}
	if err := r.Get(ctx, req.NamespacedName, meta); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	data, err := meta.ToJSON()
	if err != nil {
		log.Error(err, "could not marshal meta CR to JSON")
	}
	log.Info(string(data))

	// TODO
	// 1. find corresponding imageset (might not exist, might already be validated, can skip)
	// 2. validate (might already be valid, otherwise extract bundles + index)

	return ctrl.Result{}, nil
}
