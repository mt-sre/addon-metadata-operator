package validator

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mt-sre/client"
)

type QuayClient interface {
	HasReference(context.Context, ImageReference) (bool, error)
}

func NewQuayClient() *DefaultV2RegistryClient {
	return NewDefaultV2RegistryClient("https://quay.io")
}

func NewDefaultV2RegistryClient(url string) *DefaultV2RegistryClient {
	c := client.NewClient(client.WithWrapper{TransportWrapper: client.NewRetryWrapper()})

	return &DefaultV2RegistryClient{
		baseURL: url,
		client:  c,
	}
}

type DefaultV2RegistryClient struct {
	baseURL string
	client  *client.Client
}

func (c *DefaultV2RegistryClient) HasReference(ctx context.Context, ref ImageReference) (bool, error) {
	res, err := c.client.Head(ctx, fmt.Sprintf("%s/v2/%s/manifests/%s", c.baseURL, ref.ShortName(), ref.Tag()))
	if err != nil {
		return false, fmt.Errorf("sending HTTP request: %w", err)
	}

	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

type ImageReference interface {
	ShortName() string
	Tag() string
}
