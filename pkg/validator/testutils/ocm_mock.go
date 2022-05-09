package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
)

func NewMockOCMClient() *MockOCMClient {
	return &MockOCMClient{}
}

// MockOCMClient satisfies the types.OCMClient interface and is used to test validators
// interacting with OCM.
type MockOCMClient struct {
	mock.Mock
}

func (m *MockOCMClient) QuotaRuleExists(ctx context.Context, ocmQuotaName string) (bool, error) {
	args := m.Called(ctx, ocmQuotaName)

	return args.Bool(0), args.Error(1)
}

func (m *MockOCMClient) CloseConnection() error {
	args := m.Called()

	return args.Error(1)
}
