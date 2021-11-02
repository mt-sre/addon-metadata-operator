package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/operator-framework/operator-registry/pkg/containertools"
	"github.com/operator-framework/operator-registry/pkg/lib/indexer"
	log "github.com/sirupsen/logrus"
)

const (
	defaultDownloadPath  = "/tmp/mtcli"
	defaultCacheLocation = "/tmp/mtcli/.cache"
)

type DefaultIndexImageExtractor struct {
	downloadPath  string
	cacheLocation string
}

func (obj DefaultIndexImageExtractor) ExtractBundlesFromImage(indexImage, extractTo string) error {
	indexExporter := indexer.NewIndexExporter(
		containertools.NewContainerTool("", containertools.NoneTool),
		log.NewEntry(log.StandardLogger()),
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
	return filepath.Join(obj.downloadPath, addonName)
}

func (obj DefaultIndexImageExtractor) ExtractionPath() string {
	return obj.downloadPath
}

func (obj DefaultIndexImageExtractor) CacheLocation() string {
	return obj.cacheLocation
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
