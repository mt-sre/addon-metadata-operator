package testutils

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"

	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

/*
   Helper functions used to load the reference addon. The reference addon will
   be our perfect addon, successfully passing all validators. Hardcoding stage
   as the only environment for now.

   Bundles live here:
   	<internal_gitlab>/managed-tenants-bundles/addons/reference-addon/**

   Catalog/Bundles images live here:
    - https://quay.io/osd-addons/reference-addon-bundle:<tag>
    - https://quay.io/osd-addons/reference-addon-index:<tag>
*/

type ReferenceAddonStage struct {
	MetaIndexImage *addonsv1alpha1.AddonMetadataSpec
	MetaImageSet   *addonsv1alpha1.AddonMetadataSpec
	Env            string
}

var (
	ReferenceAddonImageSetDir   = path.Join(AddonsImagesetDir, "reference-addon")
	ReferenceAddonIndexImageDir = path.Join(AddonsIndexImageDir, "reference-addon")
	instanceMutex               sync.Mutex
	metaBundleMutex             sync.Mutex
	instance                    *ReferenceAddonStage
)

// GetReferenceAddonStage - uses singleton pattern to avoid loading yaml manifests over and over
// currently supports:
// - (DEPRECATED) static indexImage reference-addon
// - imageSet reference-addon
func GetReferenceAddonStage() (*ReferenceAddonStage, error) {
	defer instanceMutex.Unlock()
	instanceMutex.Lock()

	if instance == nil {
		instance = &ReferenceAddonStage{Env: "stage"}
		metaIndexImage, err := instance.GetMetadata(false)
		if err != nil {
			return nil, fmt.Errorf("Could not load indexImage metadata for reference-addon, got %v.", err)
		}
		metaImageSet, err := instance.GetMetadata(true)
		if err != nil {
			return nil, fmt.Errorf("Could not load imageSet metadata for reference-addon, got %v.", err)
		}
		instance.MetaIndexImage = metaIndexImage
		instance.MetaImageSet = metaImageSet
	}
	return instance, nil
}

func (r *ReferenceAddonStage) ImageSetDir() string {
	return path.Join(AddonsImagesetDir, "reference-addon")
}

func (r *ReferenceAddonStage) IndexImageDir() string {
	return path.Join(AddonsIndexImageDir, "reference-addon")
}

func (r *ReferenceAddonStage) GetMetadata(useImageSet bool) (*addonsv1alpha1.AddonMetadataSpec, error) {
	var metaPath string
	if useImageSet {
		metaPath = path.Join(r.ImageSetDir(), "metadata", r.Env, "addon.yaml")
	} else {
		metaPath = path.Join(r.IndexImageDir(), "metadata", r.Env, "addon.yaml")
	}
	data, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	meta := &addonsv1alpha1.AddonMetadataSpec{}
	err = meta.FromYAML(data)
	return meta, err
}

func (r *ReferenceAddonStage) GetImageSet(version string) (*addonsv1alpha1.AddonImageSetSpec, error) {
	baseDir := path.Join(r.ImageSetDir(), "addonimagesets", r.Env)
	target := fmt.Sprintf("reference-addon.v%v.yaml", version)
	if version == "latest" {
		latest, err := utils.GetLatestImageSetVersion(baseDir)
		if err != nil {
			return nil, err
		}
		target = latest
	} else if version == "" {
		// fallback to default value from addon.yaml metadata
		target = fmt.Sprintf("reference-addon.v%v.yaml", *r.MetaImageSet.ImageSetVersion)
	}
	imageSetPath := path.Join(baseDir, target)
	data, err := ioutil.ReadFile(imageSetPath)
	if err != nil {
		return nil, err
	}
	imageSet := &addonsv1alpha1.AddonImageSetSpec{}
	err = imageSet.FromYAML(data)
	return imageSet, err
}

func (r *ReferenceAddonStage) GetMetaBundle(version string) (*utils.MetaBundle, error) {
	defer metaBundleMutex.Unlock()
	metaBundleMutex.Lock()

	imageSet, err := r.GetImageSet(version)
	if err != nil {
		return nil, fmt.Errorf("Could not get reference-addon imageset got %v.", err)
	}
	bundles, err := utils.ExtractAndParseAddons(imageSet.IndexImage, "reference-addon")
	if err != nil {
		return nil, fmt.Errorf("Could not extract reference-addon bundles, got %v.", err)
	}
	// patch metadata - deep copy because we don't want to mess up the singleton
	meta := r.MetaImageSet.DeepCopy()
	if err := meta.PatchWithImageSet(imageSet); err != nil {
		return nil, fmt.Errorf("Could not patch metadata with imageset, got %v.", err)
	}

	return utils.NewMetaBundle(meta, bundles), nil
}
