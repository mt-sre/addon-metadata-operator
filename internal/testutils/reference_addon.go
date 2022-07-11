package testutils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/extractor"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
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

type singleton struct {
	MetaIndexImage *addonsv1alpha1.AddonMetadataSpec
	MetaImageSet   *addonsv1alpha1.AddonMetadataSpec
	Env            string
	Extractor      extractor.Extractor
}

var (
	instance *singleton
	lock     = sync.Mutex{}
)

// GetReferenceAddonStage - uses singleton pattern to avoid loading yaml manifests over and over
// currently supports:
// - (DEPRECATED) static indexImage reference-addon
// - imageSet reference-addon
func GetReferenceAddonStage() (*singleton, error) {
	lock.Lock()
	defer lock.Unlock()

	if instance != nil {
		return instance, nil
	}
	instance = &singleton{
		Env:       "stage",
		Extractor: extractor.New(),
	}
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
	return instance, nil
}

func (r *singleton) ImageSetDir() string {
	return filepath.Join(RootDir().TestData().MetadataV1().ImageSets(), "reference-addon")
}

func (r *singleton) IndexImageDir() string {
	return filepath.Join(RootDir().TestData().MetadataV1().Legacy(), "reference-addon")
}

func (r *singleton) GetMetadata(useImageSet bool) (*addonsv1alpha1.AddonMetadataSpec, error) {
	var metaPath string
	if useImageSet {
		metaPath = filepath.Join(r.ImageSetDir(), "metadata", r.Env, "addon.yaml")
	} else {
		metaPath = filepath.Join(r.IndexImageDir(), "metadata", r.Env, "addon.yaml")
	}
	data, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}
	meta := &addonsv1alpha1.AddonMetadataSpec{}
	err = meta.FromYAML(data)
	return meta, err
}

func (r *singleton) GetImageSet(version string) (*addonsv1alpha1.AddonImageSetSpec, error) {
	baseDir := filepath.Join(r.ImageSetDir(), "addonimagesets", r.Env)
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
	imageSetPath := filepath.Join(baseDir, target)
	data, err := ioutil.ReadFile(imageSetPath)
	if err != nil {
		return nil, err
	}
	imageSet := &addonsv1alpha1.AddonImageSetSpec{}
	err = imageSet.FromYAML(data)
	return imageSet, err
}

func (r *singleton) GetMetaBundle(version string) (*types.MetaBundle, error) {
	imageSet, err := r.GetImageSet(version)
	if err != nil {
		return nil, fmt.Errorf("Could not get reference-addon imageset got %v.", err)
	}
	bundles, err := r.Extractor.ExtractBundles(imageSet.IndexImage, "reference-addon")
	if err != nil {
		return nil, fmt.Errorf("Could not extract reference-addon bundles, got %v.", err)
	}

	// resolve metadata - deep copy because we don't want to mess up the singleton
	meta := r.MetaImageSet.DeepCopy()
	combinedMeta, err := meta.CombineWithImageSet(imageSet)
	if err != nil {
		return nil, fmt.Errorf("Could not combine metadata with imageset, got %v.", err)
	}

	return types.NewMetaBundle(combinedMeta, bundles), nil
}
