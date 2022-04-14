package cmd

import (
	"fmt"
	"os"

	"github.com/mt-sre/addon-metadata-operator/api/v1alpha1"
	"github.com/mt-sre/addon-metadata-operator/pkg/utils"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type convertOptions struct {
	Env     string
	Version string
}

func init() {
	var opts convertOptions

	cmd := &cobra.Command{
		Use:   "convert ADDON_DIR",
		Short: "Convert metadata to Addon Custom Resource",
		RunE:  convert(&opts),
	}

	cmd.Flags().StringVar(&opts.Env, "environment", "stage", "specifies which environment's metadata to convert")
	cmd.Flags().StringVar(&opts.Version, "version", "", "specifies which imageset version to load")

	mtcli.AddCommand(cmd)
}

func convert(opts *convertOptions) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		addonDir, err := parseAddonDir(args[0])
		if err != nil {
			return fmt.Errorf("reading addon directory: %w", err)
		}

		meta, err := utils.NewMetaLoader(addonDir, opts.Env, opts.Version).Load()
		if err != nil {
			return fmt.Errorf("loading addon metadta: %w", err)
		}

		cr := v1alpha1.AddonMetadata{
			TypeMeta: metav1.TypeMeta{
				Kind:       v1alpha1.AddonMetadataKind,
				APIVersion: v1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: meta.ID,
			},
			Spec: *meta,
		}

		data, err := cr.ToYAML()
		if err != nil {
			return fmt.Errorf("converting addon to yaml: %w", err)
		}

		if _, err := fmt.Fprint(os.Stdout, string(data)); err != nil {
			return fmt.Errorf("printing output: %w", err)
		}

		return nil
	}
}
