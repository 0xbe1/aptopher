package aptos

import (
	"net/http"
	"time"
)

// ClientConfig contains configuration for the Aptos client.
type ClientConfig struct {
	// NodeURL is the URL of the Aptos node REST API.
	NodeURL string

	// HTTPClient is an optional custom HTTP client.
	// If nil, a default client with 30 second timeout is used.
	HTTPClient *http.Client

	// Timeout is the default timeout for API requests.
	// If zero, defaults to 30 seconds.
	Timeout time.Duration
}

// Predefined network configurations.
var (
	// MainnetConfig is the configuration for Aptos mainnet.
	MainnetConfig = ClientConfig{
		NodeURL: "https://fullnode.mainnet.aptoslabs.com/v1",
	}

	// TestnetConfig is the configuration for Aptos testnet.
	TestnetConfig = ClientConfig{
		NodeURL: "https://fullnode.testnet.aptoslabs.com/v1",
	}

	// DevnetConfig is the configuration for Aptos devnet.
	DevnetConfig = ClientConfig{
		NodeURL: "https://fullnode.devnet.aptoslabs.com/v1",
	}

	// LocalnetConfig is the configuration for a local Aptos node.
	LocalnetConfig = ClientConfig{
		NodeURL: "http://127.0.0.1:8080/v1",
	}
)
