package testutils

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/mt-sre/addon-metadata-operator/internal/testutils"
	"github.com/mt-sre/addon-metadata-operator/pkg/types"
	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewValidatorTester(t *testing.T, init validator.Initializer, opts ...ValidatorTesterOption) ValidatorTester {
	t.Helper()

	vt := ValidatorTester{T: t}

	for _, opt := range opts {
		opt(&vt)
	}

	if vt.log == nil {
		vt.log = logr.Discard()
	}

	var err error

	// This also ensures that a validator implements the validator.Validator interface
	vt.Val, err = init(validator.Dependencies{
		Logger:    vt.log,
		OCMClient: vt.ocm,
	})
	require.NoError(t, err)

	return vt
}

type ValidatorTester struct {
	*testing.T
	Val validator.Validator
	log logr.Logger
	ocm validator.OCMClient
}

func (v *ValidatorTester) TestSingleBundle(mb types.MetaBundle) validator.Result {
	return v.Val.Run(context.Background(), mb)
}

func (v *ValidatorTester) TestValidBundles(bundles map[string]types.MetaBundle) {
	v.Helper()

	v.testBundles(bundles, assert.True)
}

func (v *ValidatorTester) TestInvalidBundles(bundles map[string]types.MetaBundle) {
	v.Helper()

	v.testBundles(bundles, assert.False)
}

func (v *ValidatorTester) testBundles(bundles map[string]types.MetaBundle, assert assert.BoolAssertionFunc) {
	v.Helper()

	for name, bundle := range bundles {
		bundle := bundle

		v.Run(name, func(t *testing.T) {
			t.Parallel()

			res := v.Val.Run(context.Background(), bundle)
			require.False(t, res.IsError())
			assert(t, res.IsSuccess(), "Actual Result: %+v", res)
		})
	}
}

func (v *ValidatorTester) Option(opt ValidatorTesterOption) { opt(v) }

type ValidatorTesterOption func(*ValidatorTester)

func ValidatorTesterLogger(l logr.Logger) ValidatorTesterOption {
	return func(v *ValidatorTester) {
		v.log = l
	}
}

func ValidatorTesterOCMClient(ocm validator.OCMClient) ValidatorTesterOption {
	return func(v *ValidatorTester) {
		v.ocm = ocm
	}
}

func DefaultValidBundleMap() (map[string]types.MetaBundle, error) {
	res := make(map[string]types.MetaBundle)

	refAddonStage, err := testutils.GetReferenceAddonStage()
	if err != nil {
		return nil, fmt.Errorf("unable to get reference addon: %w", err)
	}

	refAddonMetaBundle, err := refAddonStage.GetMetaBundle(*refAddonStage.MetaImageSet.ImageSetVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to get reference addon meta bundle: %w", err)
	}

	res["reference addon"] = *refAddonMetaBundle

	return res, nil
}
