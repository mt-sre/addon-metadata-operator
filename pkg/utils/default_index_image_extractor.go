package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/operator-framework/operator-registry/pkg/containertools"
	"github.com/operator-framework/operator-registry/pkg/lib/indexer"
	log "github.com/sirupsen/logrus"
)

const (
	// (sblaisdo) sretoolbox also uses /tmp/mtcli to extract binary archive
	DefaultDownloadPath  = "/tmp/mtcli-07b10894-0673-4d95-b6ef-0cbd9701c9c3"
	DefaultCacheDir      = "/tmp/mtcli-07b10894-0673-4d95-b6ef-0cbd9701c9c3"
	DefaultCacheFileName = ".cache"
)

type DefaultIndexImageExtractor struct {
	downloadPath string
	cacheDir     string
	indexImage   string
}

func (obj DefaultIndexImageExtractor) ExtractBundlesFromImage(indexImage, extractTo string) error {
	// Write all index image extaction logs to /dev/null
	logger := log.StandardLogger()
	indexExporter := indexer.NewIndexExporter(
		containertools.NewContainerTool("", containertools.NoneTool),
		log.NewEntry(logger),
	)
	request := indexer.ExportFromIndexRequest{
		Index:         indexImage,
		Packages:      []string{},
		DownloadPath:  extractTo,
		ContainerTool: containertools.NewContainerTool("", containertools.NoneTool),
		SkipTLS:       false,
	}
	return indexExporter.ExportFromIndex(request)
}

func (obj DefaultIndexImageExtractor) ManifestsPath(addonName string) string {
	return filepath.Join(obj.downloadPath, obj.indexImage, addonName)
}

func (obj DefaultIndexImageExtractor) ExtractionPath() string {
	return filepath.Join(obj.downloadPath, obj.indexImage)
}

func (obj DefaultIndexImageExtractor) CacheLocation() string {
	return filepath.Join(obj.cacheDir, obj.indexImage, DefaultCacheFileName)
}

func (obj DefaultIndexImageExtractor) CacheKey(indexImage, addonName string) string {
	return strings.Join(
		[]string{indexImage, addonName},
		"<>",
	)
}

func (obj DefaultIndexImageExtractor) CacheHit(key string) bool {
	cacheLocation := obj.CacheLocation()
	if checkFileExists(cacheLocation) {
		data, err := os.ReadFile(cacheLocation)
		if err != nil {
			return false
		}
		return string(data) == key
	}
	return false
}

func (obj DefaultIndexImageExtractor) WriteToCache(value string) error {
	cacheLocation := obj.CacheLocation()
	return os.WriteFile(cacheLocation, []byte(value), 0644)
}

func checkFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// validate interface implemented
var _ = types.IndexImageExtractor(DefaultIndexImageExtractor{})
