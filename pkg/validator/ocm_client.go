package validator

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

// OCMClient abstracts behavior required for validators which request data
// from OCM to be implemented by OCM API clients.
type OCMClient interface {
	QuotaRuleGetter
}

type QuotaRuleGetter interface {
	// QuotaRuleExists takes a given quota rule name and returns a tuple
	// of ('ok', error) which returns 'true' if the quota rule exists
	// and false otherwise. An optional error is returned if any issues
	// occurred.
	QuotaRuleExists(context.Context, string) (bool, error)
}

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

// NewOCMClient takes a variadic slice of options to configure a default
// OCM client and applies defaults if no appropriate option is given. An error
// may be returned if an unusable OCM token is provided or a connection cannot
// be initialized otherwise the default client is returned.
func NewOCMClient(opts ...OCMClientOption) (*OCMClientImpl, error) {
	var cfg OCMClientConfig

	cfg.Option(opts...)
	cfg.Default()

	conn, err := cfg.Connector.Connect(cfg.ConnectOptions...)
	if err != nil {
		return nil, fmt.Errorf("connecting to OCM: %w", err)
	}

	return &OCMClientImpl{
		cfg:  cfg,
		conn: conn,
	}, nil
}

// OCMClientImpl implements the 'types.OCMClient' interface and
// exposes methods by which validators can communicate with OCM.
type OCMClientImpl struct {
	cfg  OCMClientConfig
	conn *sdk.Connection
}

func (c *OCMClientImpl) QuotaRuleExists(ctx context.Context, quotaName string) (bool, error) {
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

func isHTTPError(code int) bool {
	return code >= 400 && code < 600
}

// CloseConnection releases any resources held by the connection to
// OCM.
func (c *OCMClientImpl) CloseConnection() error { return c.conn.Close() }

type OCMClientConfig struct {
	Connector      OCMConnector
	ConnectOptions []OCMConnectOption
}

func (c *OCMClientConfig) Option(opts ...OCMClientOption) {
	for _, opt := range opts {
		opt.ConfigureOCMClient(c)
	}
}

func (c *OCMClientConfig) Default() {
	if c.Connector == nil {
		c.Connector = NewOCMConnector()
	}
}

type OCMClientOption interface {
	ConfigureOCMClient(*OCMClientConfig)
}

// OCMConnector establishes and returns a connection to OCM.
type OCMConnector interface {
	// Connect takes a variadic slice of OCMConnectOptions
	// to configure and return an open connection to OCM.
	// Returns an error if connection fails.
	Connect(opts ...OCMConnectOption) (*sdk.Connection, error)
}

type OCMConnectionConfig struct {
	APIURL       string
	AccessToken  string
	ClientID     string
	ClientSecret string
}

func (c *OCMConnectionConfig) Option(opts ...OCMConnectOption) {
	for _, opt := range opts {
		opt.ConfigureOCMConnection(c)
	}
}

const (
	apiURLStage = "https://api.stage.openshift.com"
)

func (c *OCMConnectionConfig) Default() {
	if c.APIURL == "" {
		c.APIURL = apiURLStage
	}
}

type OCMConnectOption interface {
	ConfigureOCMConnection(*OCMConnectionConfig)
}

// NewOCMConnector returns an initialized OCMConnector instance.
func NewOCMConnector() *OCMConnectorImpl {
	return &OCMConnectorImpl{}
}

type OCMConnectorImpl struct{}

func (c *OCMConnectorImpl) Connect(opts ...OCMConnectOption) (*sdk.Connection, error) {
	var cfg OCMConnectionConfig

	cfg.Option(opts...)
	cfg.Default()

	builder := sdk.NewConnectionBuilder().URL(cfg.APIURL)

	if cfg.ClientID != "" && cfg.ClientSecret != "" {
		builder = builder.Client(cfg.ClientID, cfg.ClientSecret)
	} else {
		builder = builder.Tokens(cfg.AccessToken)
	}

	return builder.Build()
}

// NewDisconnectedOCMClient returns an OCM Client which fails
// on any operations which call OCM. Helpful to trigger failure
// only for validators which depend on OCM.
func NewDisconnectedOCMClient() DisconnectedOCMClient { return DisconnectedOCMClient{} }

type DisconnectedOCMClient struct{}

var ErrDisconnectedOCMClient = errors.New("OCM client disconnected")

func (c DisconnectedOCMClient) QuotaRuleExists(_ context.Context, _ string) (bool, error) {
	return false, ErrDisconnectedOCMClient
}
