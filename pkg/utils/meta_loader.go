package utils

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	log "github.com/sirupsen/logrus"

	addonsv1alpha1 "github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
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
		combinedMeta, err := meta.CombineWithImageSet(imageSet)
		if err != nil {
			return nil, fmt.Errorf("Could not combine metadata and imageset, got %v.", err)
		}
		return combinedMeta, nil
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
	return filepath.Join(l.AddonDir, "metadata", l.Env, "addon.yaml")
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
	baseDir := filepath.Join(l.AddonDir, "addonimagesets", l.Env)
	target := fmt.Sprintf("%s.v%s.yaml", l.AddonName, version)
	if version == "latest" {
		latest, err := GetLatestImageSetVersion(baseDir)
		if err != nil {
			return "", err
		}
		target = latest
	}
	return filepath.Join(baseDir, target), nil
}

func GetLatestImageSetVersion(dir string) (string, error) {
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
