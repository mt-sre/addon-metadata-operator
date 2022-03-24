package controllers

import (
	addonmetadata "github.com/mt-sre/addon-metadata-operator/internal/controllers/addonmetadata"
	ctrl "sigs.k8s.io/controller-runtime"
)

func SetupAddonMetadataReconciler(mgr ctrl.Manager) error {
	reconciler := &addonmetadata.AddonMetadataReconciler{
		Log:    ctrl.Log.WithName("controllers").WithName("AddonMetadata"),
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
	return reconciler.SetupWithManager(mgr)
}
