package testutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var (
	RootDir             string = getAndValidateRootDir()
	TestdataDir         string = path.Join(RootDir, "internal", "testdata")
	AddonsImagesetDir   string = path.Join(TestdataDir, "addons-imageset")
	AddonsIndexImageDir string = path.Join(TestdataDir, "addons-indeximage")
)

func getAndValidateRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	root := path.Join(filepath.Dir(b), "..", "..")

	if !dirContainsGoMod(root) {
		log.Fatal("could not find go.mod in root directory: ", root)
	}
	return root
}

func dirContainsGoMod(root string) bool {
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal("can't read root directory, got: ", err)
	}
	for _, file := range files {
		if file.Name() == "go.mod" {
			return true
		}
	}
	return false
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
func DefaultSucceedingCandidates() ([]utils.MetaBundle, error) {
	var res []utils.MetaBundle
	refAddonStage, err := GetReferenceAddonStage()
	if err != nil {
		return nil, fmt.Errorf("Could not load reference-addon, got %v.", err)
	}
	refAddonMetaBundle, err := refAddonStage.GetMetaBundle(*refAddonStage.MetaImageSet.ImageSetVersion)
	if err != nil {
		return nil, fmt.Errorf("Could not get reference-addon meta bundles, got %v.", err)
	}
	res = append(res, *refAddonMetaBundle)

	return res, nil
}
