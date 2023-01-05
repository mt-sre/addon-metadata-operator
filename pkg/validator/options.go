package validator

// WithConnector applies the given OCMConnector implementation.
type WithConnector struct{ Connector OCMConnector }

func (w WithConnector) ConfigureOCMClient(c *OCMClientConfig) {
	c.Connector = w.Connector
}

// WithConnectOptions applies the given OCMConnectOption's
type WithConnectOptions []OCMConnectOption

func (w WithConnectOptions) ConfigureOCMClient(c *OCMClientConfig) {
	c.ConnectOptions = []OCMConnectOption(w)
}

// WithAPIURL applies the given API URL.
type WithAPIURL string

func (w WithAPIURL) ConfigureOCMConnection(c *OCMConnectionConfig) {
	c.APIURL = string(w)
}

// WithAccessToken applies the given access token.
type WithAccessToken string

func (w WithAccessToken) ConfigureOCMConnection(c *OCMConnectionConfig) {
	c.AccessToken = string(w)
}

// WithClientID applies the given client ID.
type WithClientID string

func (w WithClientID) ConfigureOCMConnection(c *OCMConnectionConfig) {
	c.ClientID = string(w)
}

// WithClientSecret applies the given client secret.
type WithClientSecret string

func (w WithClientSecret) ConfigureOCMConnection(c *OCMConnectionConfig) {
	c.ClientSecret = string(w)
}
