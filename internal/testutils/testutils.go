package testutils

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var (
	RootDir             string = getRootDir()
	TestdataDir         string = path.Join(RootDir, "internal", "testdata")
	AddonsImagesetDir   string = path.Join(TestdataDir, "addons-imageset")
	AddonsIndexImageDir string = path.Join(TestdataDir, "addons-indeximage")
)

func getRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return path.Join(filepath.Dir(b), "..", "..")
}

func RemoveDir(downloadDir string) {
	err := os.RemoveAll(downloadDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to cleanup download dir")
	}
}

func GetStringLiteralRef(s string) *string { return &s }

// DefaultSucceedingCandidates - returns a slice of valid metaBundles that are supposed
// to pass all validators successfully. If it is not the case, please make the required adjustments.
func DefaultSucceedingCandidates() []utils.MetaBundle {
	var res []utils.MetaBundle
	refAddonStage := GetReferenceAddonStage()
	refAddonMetaBundle, err := refAddonStage.GetMetaBundle(*refAddonStage.MetaImageSet.ImageSetVersion)
	if err != nil {
		log.Fatalf("Could not get reference-addon meta bundles, got %v.", err)
	}
	res = append(res, *refAddonMetaBundle)

	return res
}
