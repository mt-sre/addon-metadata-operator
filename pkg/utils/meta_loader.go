package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"

	log "github.com/sirupsen/logrus"

	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	ocmv1 "github.com/mt-sre/addon-metadata-operator/pkg/ocm/v1"
)

type MetaLoader interface {
	Load() (*addonsv1alpha1.AddonMetadataSpec, error)
}

type defaultMetaLoader struct {
	AddonDir  string
	AddonName string
	Env       string
	Version   string
}

// NewMetaLoader - returns default implementation of the AddonMetaLoader
func NewMetaLoader(addonDir, env, version string) MetaLoader {
	return defaultMetaLoader{
		AddonDir:  addonDir,
		AddonName: path.Base(addonDir),
		Env:       env,
		Version:   version,
	}
}

// Load - loads the addon metadata and imageSet
func (l defaultMetaLoader) Load() (*addonsv1alpha1.AddonMetadataSpec, error) {
	meta, err := l.readMeta()
	if err != nil {
		return nil, err
	}
	// invalid - legacy addon
	if meta.IndexImage == nil && meta.ImageSetVersion == nil {
		return nil, errors.New("No validation support for legacy addon. Please use the imageSet feature.")
	}
	// invalid - misconfiguration
	if meta.IndexImage != nil && meta.ImageSetVersion != nil {
		return nil, errors.New("Can't set both the 'indexImage' and the 'imageSetVersion' field.")
	}
	// imageSet
	if meta.ImageSetVersion != nil {
		imageSet, err := l.readImageSet(*meta.ImageSetVersion)
		if err != nil {
			return nil, fmt.Errorf("Could not read imageSet, got %v.\n", err)
		}
		// deepcopy to be safe
		meta.IndexImage = &imageSet.IndexImage
		// TODO(ykukreja): function to perform safer and consistent deep copies and use that function here
		if meta.AddOnParameters != nil && imageSet.AddOnParameters != nil {
			*meta.AddOnParameters = make([]ocmv1.AddOnParameter, len(*imageSet.AddOnParameters))
			copy(*meta.AddOnParameters, *imageSet.AddOnParameters)
		}
		if meta.AddOnRequirements != nil && imageSet.AddOnRequirements != nil {
			*meta.AddOnRequirements = make([]ocmv1.AddOnRequirement, len(*imageSet.AddOnRequirements))
			copy(*meta.AddOnRequirements, *imageSet.AddOnRequirements)
		}
		if meta.SubOperators != nil && imageSet.SubOperators != nil {
			*meta.SubOperators = make([]ocmv1.AddOnSubOperator, len(*imageSet.SubOperators))
			copy(*meta.SubOperators, *imageSet.SubOperators)
		}
	}
	return meta, nil
}

func (l defaultMetaLoader) readMeta() (*addonsv1alpha1.AddonMetadataSpec, error) {
	data, err := ioutil.ReadFile(l.getMetadataPath())
	if err != nil {
		return nil, err
	}
	log.Debugf("Raw metadata read from addon: %v. \n%v\n", l.AddonName, string(data))
	meta := &addonsv1alpha1.AddonMetadataSpec{}
	err = meta.FromYAML(data)
	return meta, err
}

func (l defaultMetaLoader) getMetadataPath() string {
	return path.Join(l.AddonDir, "metadata", l.Env, "addon.yaml")
}

func (l defaultMetaLoader) readImageSet(defaultVersion string) (*addonsv1alpha1.AddonImageSetSpec, error) {
	version := l.getImageSetVersion(defaultVersion)
	imageSetPath, err := l.getImagesetPath(version)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(imageSetPath)
	if err != nil {
		return nil, err
	}
	log.Debugf("Raw imageSet read from addon: %v. \n%v\n", l.AddonName, string(data))
	imageSet := &addonsv1alpha1.AddonImageSetSpec{}
	err = imageSet.FromYAML(data)
	return imageSet, err
}

// defaultVersion == meta.ImageSetVersion
// Can be overriden by providing the --version CLI flag
func (l defaultMetaLoader) getImageSetVersion(defaultVersion string) string {
	if l.Version != "" {
		return l.Version
	}
	return defaultVersion
}

func (l defaultMetaLoader) getImagesetPath(version string) (string, error) {
	baseDir := path.Join(l.AddonDir, "addonimagesets", l.Env)
	target := fmt.Sprintf("%s.v%s.yaml", l.AddonName, version)
	if version == "latest" {
		latest, err := getLatestImageSetVersion(baseDir)
		if err != nil {
			return "", err
		}
		target = latest
	}
	return path.Join(baseDir, target), nil
}

func getLatestImageSetVersion(dir string) (string, error) {
	sortDescending := func(files []os.FileInfo) {
		sort.Slice(files, func(i, j int) bool {
			return files[i].Name() > files[j].Name()
		})
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("No imageset present in the directory.")
	}
	sortDescending(files)
	return files[0].Name(), nil
}
