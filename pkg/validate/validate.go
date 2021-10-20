package validate

import (
	"fmt"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
)

var (
	success = utils.Green("Success")
	failed  = utils.Red("Failed")
	running = utils.Bold("Running")
)

type Validator struct {
	Description string
	Runner      ValidateFunc
}

type ValidateFunc func(mb *MetaBundle) error

type MetaBundle struct {
	AddonMeta *v1alpha1.AddonMetadataSpec
	// TODO: add field for corresponding bundle
}

// TODO: This will return a MetaBundle with corresponding bundle
func NewMetaBundle(addonMeta *v1alpha1.AddonMetadataSpec) *MetaBundle {
	return &MetaBundle{
		AddonMeta: addonMeta,
	}
}

func (mb *MetaBundle) Validate(runMeta bool) []error {
	errs := []error{}

	if runMeta {
		validators := GetAllMetaValidators()
		fmt.Printf("%s\n", utils.Bold("Running metadata validators"))
		fmt.Println()
		for _, validator := range validators {
			fmt.Printf("\r%s\t\t%s", validator.Description, running)
			err := validator.Runner(mb)
			if err != nil {
				errs = append(errs, err)
				fmt.Printf("\r%s\t\t%s", validator.Description, failed)
			}
			fmt.Printf("\r%s\t\t%s", validator.Description, success)
			fmt.Println()
		}
	}
	return errs
}
