package utils

import (
	"context"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
	amv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1"
)

const (
	apiURLStage    = "https://api.stage.openshift.com"
	ocmTokenEnvVar = "OCM_TOKEN"
)

// OCMResponseError is used to wrap HTTP error (400 - 599) response codes
// which are returned from a request to OCM.
type OCMResponseError int

func (e OCMResponseError) Error() string {
	return fmt.Sprintf("ocm responded with code %d", e)
}

func (e OCMResponseError) ServerSide() bool {
	code := int(e)

	return code >= 500 && code < 600
}

// NewDefaultOCMClient takes a variadic slice of options to configure a default
// OCM client and applies defaults if no appropriate option is given. An error
// may be returned if an unusable OCM token is provided or a connection cannot
// be initialized otherwise the default client is returned.
func NewDefaultOCMClient(opts ...DefaultOCMClientOption) (*DefaultOCMClient, error) {
	var client DefaultOCMClient

	for _, opt := range opts {
		client.Option(opt)
	}

	if client.apiURL == "" {
		client.apiURL = apiURLStage
	}

	if client.tp == nil {
		client.tp = NewEnvOCMTokenProvider(ocmTokenEnvVar)
	}

	if err := client.initConnection(); err != nil {
		return nil, err
	}

	return &client, nil
}

// DefaultOCMClient implements the 'types.OCMClient' interface and
// exposes methods by which validators can communicate with OCM.
type DefaultOCMClient struct {
	apiURL string
	conn   *sdk.Connection
	tp     OCMTokenProvider
}

func (c *DefaultOCMClient) initConnection() error {
	token, err := c.tp.ProvideToken()
	if err != nil {
		return err
	}

	c.conn, err = sdk.NewConnectionBuilder().
		URL(c.apiURL).
		Tokens(token).
		Build()

	return err
}

func (c *DefaultOCMClient) GetSKURules(ctx context.Context, quotaName string) ([]*amv1.SkuRule, error) {
	query := fmt.Sprintf("quota_id = '%s'", addonQuotaID(quotaName))

	req := c.conn.AccountsMgmt().V1().SkuRules().List().Search(query)

	res, err := req.SendContext(ctx)
	if err != nil {
		return nil, err
	}

	if isHTTPError(res.Status()) {
		return nil, OCMResponseError(res.Status())
	}

	return res.Items().Slice(), nil
}

func (c *DefaultOCMClient) CloseConnection() error {
	return c.conn.Close()
}

func addonQuotaID(quotaName string) string {
	return "add-on|" + quotaName
}

func isHTTPError(code int) bool {
	return code >= 400 && code < 600
}

// Option applies the given option to an instance of the DefaultOCMClient.
func (c *DefaultOCMClient) Option(opt DefaultOCMClientOption) {
	opt(c)
}

// DefaultOCMClientOption describes a function which configures a DefaultOCMClient instance.
type DefaultOCMClientOption func(c *DefaultOCMClient)

// DefaultOCMClientAPIURL updates the URL used to connect to OCM.
// This option only takes effect if applied prior to connection initialization.
func DefaultOCMClientAPIURL(url string) DefaultOCMClientOption {
	return func(c *DefaultOCMClient) {
		c.apiURL = url
	}
}

// DefaultOCMClientAPIURL updates the URL used to connect to OCM.
// This option only takes effect if applied prior to connection initialization.
func DefaultOCMClientTokenProvider(tp OCMTokenProvider) DefaultOCMClientOption {
	return func(c *DefaultOCMClient) {
		c.tp = tp
	}
}

// OCMTokenProvider provides OCM tokens to clients.
type OCMTokenProvider interface {
	// ProvideToken returns either an access token as an encoded JWT token
	// or an error if a token could not be provided.
	ProvideToken() (string, error)
}

// NewEnvOCMTokenProvider returns an instance of EnvOCMTokenProvider
// which retrieves OCM tokens from the environment by the given envVar name.
func NewEnvOCMTokenProvider(envVar string) EnvOCMTokenProvider {
	return EnvOCMTokenProvider{
		envVar: envVar,
	}
}

// EnvOCMTokenProvider provides OCM tokens from the current environment.
type EnvOCMTokenProvider struct {
	envVar string
}

func (tp EnvOCMTokenProvider) ProvideToken() (string, error) {
	if token, ok := os.LookupEnv(tp.envVar); ok {
		return token, nil
	}

	return "", fmt.Errorf("utils: ocm token environment variable '%s' not set", tp.envVar)
}
