package balance

import (
	"context"
	"errors"
	"fmt"

	"github.com/brave-intl/bat-go/utils/clients"
	uuid "github.com/satori/go.uuid"
)

// Client abstracts over the underlying client
type Client interface {
	InvalidateBalance(ctx context.Context, id uuid.UUID) error
}

// HTTPClient wraps http.Client for interacting with the ledger server
type HTTPClient struct {
	client *clients.SimpleHTTPClient
}

// New returns a new HTTPClient, retrieving the base URL from the environment
func New(ctx context.Context) (*HTTPClient, error) {
	serverEnvKey := "BALANCE_SERVER"
	serverURL, err := appctx.GetConfValue(ctx, serverEnvKey)
	if err != nil {
		return nil, fmt.Errorf("error getting configuration value: %w", err)
	}
	if len(serverEnvKey) == 0 {
		return nil, errors.New(serverEnvKey + " was empty")
	}

	bToken, err := appctx.GetConfValue(ctx, "BALANCE_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("error getting configuration value: %w", err)
	}

	client, err := clients.New(serverURL, bToken)
	if err != nil {
		return nil, err
	}
	return &HTTPClient{client}, err
}

// InvalidateBalance invalidates the cached value on balance
func (c *HTTPClient) InvalidateBalance(ctx context.Context, id uuid.UUID) error {
	req, err := c.client.NewRequest(ctx, "DELETE", "v2/wallet/"+id.String()+"/balance", nil)
	if err != nil {
		return err
	}

	_, err = c.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}

	return nil
}
