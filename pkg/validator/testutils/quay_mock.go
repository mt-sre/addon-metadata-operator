package testutils

import (
	"context"

	"github.com/mt-sre/addon-metadata-operator/pkg/validator"
	"github.com/stretchr/testify/mock"
)

func NewMockQuayClient() *MockQuayClient {
	return &MockQuayClient{}
}

type MockQuayClient struct {
	mock.Mock
}

func (c *MockQuayClient) HasReference(ctx context.Context, ref validator.ImageReference) (bool, error) {
	args := c.Called(ctx, ref)

	return args.Bool(0), args.Error(1)
}
