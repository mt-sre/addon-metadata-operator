package validator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	sdk "github.com/openshift-online/ocm-sdk-go"
)

// OCMError abstracts behavior required for validators to identify the underlying
// causes of OCM related errors.
type OCMError interface {
	// ServerSide returns 'true' if an instance of OCMError was caused
	// by a server-side issue.
	ServerSide() bool
}

// IsOCMServerSideError determines if the given error is both an instance of OCMError
// and was caused by a server-side issue.
func IsOCMServerSideError(err error) bool {
	ocmErr, ok := err.(OCMError)

	return ok && ocmErr.ServerSide()
}

type OCMClientError interface {
	IsAuthRelated() bool
}

func IsOCMClientAuthError(err error) bool {
	clientErr, ok := err.(OCMClientError)

	return ok && clientErr.IsAuthRelated()
}

// OCMClient abstracts behavior required for validators which request data
// from OCM to be implemented by OCM API clients.
type OCMClient interface {
	QuotaRuleGetter
	// Stuck here until OCM Client can be injected via Params...
	CloseConnection() error
}

type QuotaRuleGetter interface {
	QuotaRuleExists(context.Context, string) (bool, error)
}

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

func NewOCMClientAuthError(err error) ocmClientError {
	return ocmClientError{
		error:       err,
		authRelated: true,
	}
}

type ocmClientError struct {
	authRelated bool
	error
}

func (e ocmClientError) Unwrap() error { return e.error }
func (e ocmClientError) Error() string {
	return fmt.Sprintf("unable to setup OCM authentication: %v", e.error)
}

func (e ocmClientError) IsAuthRelated() bool { return e.authRelated }

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
		return NewOCMClientAuthError(err)
	}

	c.conn, err = sdk.NewConnectionBuilder().
		URL(c.apiURL).
		Tokens(token).
		Build()

	return err
}

func (c *DefaultOCMClient) QuotaRuleExists(ctx context.Context, quotaName string) (bool, error) {
	if c.conn == nil {
		return false, errors.New("no active OCM connection")
	}

	req := c.conn.
		Get().
		Path("/api/accounts_mgmt/v1/quota_rules").
		Parameter("search", fmt.Sprintf("name = '%s'", quotaName))

	res, err := req.SendContext(ctx)
	if err != nil {
		return false, fmt.Errorf("requesting quota rules: %w", err)
	}

	if isHTTPError(res.Status()) {
		return false, OCMResponseError(res.Status())
	}

	list := struct{ Size int }{}

	if err := json.Unmarshal(res.Bytes(), &list); err != nil {
		return false, fmt.Errorf("unmarshalling quota rules: %w", err)
	}

	return list.Size > 0, nil
}

func (c *DefaultOCMClient) CloseConnection() error {
	if c.conn == nil {
		return nil
	}

	return c.conn.Close()
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
