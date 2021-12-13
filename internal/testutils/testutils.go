package testutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
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
