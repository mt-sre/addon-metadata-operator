package testutils

import (
	"context"

	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
)

// NewMockOCMClient returns a *MockOCMClient configured with the supplied options.
func NewMockOCMClient(opts ...MockOCMClientOption) *MockOCMClient {
	var client MockOCMClient

	for _, opt := range opts {
		client.Option(opt)
	}

	return &client
}

// MockOCMClient satisfies the types.OCMClient interface and is used to test validators
// interacting with OCM.
type MockOCMClient struct {
	validQuotaNames []string
}

// Option applies the provided MockOCMClientOption to a MockOCMClientInstance
func (c *MockOCMClient) Option(opt MockOCMClientOption) {
	opt(c)
}

func (m *MockOCMClient) GetSKURules(ctx context.Context, ocmQuotaName string) ([]*amv1.SkuRule, error) {
	if !contains(ocmQuotaName, m.validQuotaNames) {
		return []*amv1.SkuRule{}, nil
	}

	return []*amv1.SkuRule{new(amv1.SkuRule)}, nil
}

func (m *MockOCMClient) CloseConnection() error {
	return nil
}

type MockOCMClientOption func(c *MockOCMClient)

// MockOCMClientValidQuotaNames populates the MockOCMClient instance with
// a list of dummy ocmQuotaNames that are guaranteed to return a SKU Rule.
func MockOCMClientValidQuotaNames(names ...string) MockOCMClientOption {
	return func(c *MockOCMClient) {
		c.validQuotaNames = names
	}
}

func contains(elem string, slice []string) bool {
	for _, s := range slice {
		if elem == s {
			return true
		}
	}

	return false
}
