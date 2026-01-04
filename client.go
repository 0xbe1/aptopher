package aptos

import (
	"context"
	"net/http"
	"time"
)

// Client is the main Aptos SDK client.
type Client struct {
	http    *httpClient
	chainID uint8
}

// NewClient creates a new Aptos client with the given configuration.
func NewClient(config ClientConfig) (*Client, error) {
	hc := config.HTTPClient
	if hc == nil {
		timeout := config.Timeout
		if timeout == 0 {
			timeout = 30 * time.Second
		}
		hc = &http.Client{Timeout: timeout}
	}

	return &Client{
		http: newHTTPClient(config.NodeURL, hc),
	}, nil
}

// GetLedgerInfo retrieves the current ledger information.
func (c *Client) GetLedgerInfo(ctx context.Context) (Response[LedgerInfo], error) {
	var info LedgerInfo
	metadata, err := c.http.get(ctx, "/", &info)
	if err != nil {
		return Response[LedgerInfo]{}, err
	}
	return Response[LedgerInfo]{Data: info, Metadata: metadata}, nil
}

// GetNodeInfo retrieves basic information about the node.
func (c *Client) GetNodeInfo(ctx context.Context) (Response[NodeInfo], error) {
	var info NodeInfo
	metadata, err := c.http.get(ctx, "/", &info)
	if err != nil {
		return Response[NodeInfo]{}, err
	}
	return Response[NodeInfo]{Data: info, Metadata: metadata}, nil
}

// HealthCheck checks if the node is healthy.
// Returns nil if healthy, or an error otherwise.
func (c *Client) HealthCheck(ctx context.Context) error {
	_, err := c.http.get(ctx, "/-/healthy", nil)
	return err
}

// EstimateGasPrice retrieves the current gas price estimation.
func (c *Client) EstimateGasPrice(ctx context.Context) (Response[GasEstimation], error) {
	var estimation GasEstimation
	metadata, err := c.http.get(ctx, "/estimate_gas_price", &estimation)
	if err != nil {
		return Response[GasEstimation]{}, err
	}
	return Response[GasEstimation]{Data: estimation, Metadata: metadata}, nil
}

// GetAccount retrieves account information including sequence number and authentication key.
func (c *Client) GetAccount(ctx context.Context, address AccountAddress, opts ...RequestOption) (Response[AccountData], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + options.BuildQueryParams()

	var account AccountData
	metadata, err := c.http.get(ctx, path, &account)
	if err != nil {
		return Response[AccountData]{}, err
	}
	return Response[AccountData]{Data: account, Metadata: metadata}, nil
}

// GetAccountResources retrieves all resources for an account.
func (c *Client) GetAccountResources(ctx context.Context, address AccountAddress, opts ...RequestOption) (Response[[]MoveResource], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + "/resources" + options.BuildQueryParams()

	var resources []MoveResource
	metadata, err := c.http.get(ctx, path, &resources)
	if err != nil {
		return Response[[]MoveResource]{}, err
	}
	return Response[[]MoveResource]{Data: resources, Metadata: metadata}, nil
}

// GetAccountResource retrieves a specific resource for an account.
func (c *Client) GetAccountResource(ctx context.Context, address AccountAddress, resourceType string, opts ...RequestOption) (Response[MoveResource], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + "/resource/" + resourceType + options.BuildQueryParams()

	var resource MoveResource
	metadata, err := c.http.get(ctx, path, &resource)
	if err != nil {
		return Response[MoveResource]{}, err
	}
	return Response[MoveResource]{Data: resource, Metadata: metadata}, nil
}

// GetAccountModules retrieves all modules for an account.
func (c *Client) GetAccountModules(ctx context.Context, address AccountAddress, opts ...RequestOption) (Response[[]MoveModuleBytecode], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + "/modules" + options.BuildQueryParams()

	var modules []MoveModuleBytecode
	metadata, err := c.http.get(ctx, path, &modules)
	if err != nil {
		return Response[[]MoveModuleBytecode]{}, err
	}
	return Response[[]MoveModuleBytecode]{Data: modules, Metadata: metadata}, nil
}

// GetAccountModule retrieves a specific module for an account.
func (c *Client) GetAccountModule(ctx context.Context, address AccountAddress, moduleName string, opts ...RequestOption) (Response[MoveModuleBytecode], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + "/module/" + moduleName + options.BuildQueryParams()

	var module MoveModuleBytecode
	metadata, err := c.http.get(ctx, path, &module)
	if err != nil {
		return Response[MoveModuleBytecode]{}, err
	}
	return Response[MoveModuleBytecode]{Data: module, Metadata: metadata}, nil
}

// GetAccountBalance retrieves the balance of a specific asset type for an account.
func (c *Client) GetAccountBalance(ctx context.Context, address AccountAddress, assetType string, opts ...RequestOption) (Response[uint64], error) {
	options := ApplyOptions(opts...)
	path := "/accounts/" + address.String() + "/balance/" + assetType + options.BuildQueryParams()

	var balance uint64
	metadata, err := c.http.get(ctx, path, &balance)
	if err != nil {
		return Response[uint64]{}, err
	}
	return Response[uint64]{Data: balance, Metadata: metadata}, nil
}
