package aptos

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// httpClient handles HTTP communication with the Aptos node.
type httpClient struct {
	baseURL    string
	httpClient *http.Client
}

// newHTTPClient creates a new HTTP client for the Aptos API.
func newHTTPClient(baseURL string, client *http.Client) *httpClient {
	// Ensure base URL doesn't have trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")
	if client == nil {
		client = http.DefaultClient
	}
	return &httpClient{
		baseURL:    baseURL,
		httpClient: client,
	}
}

// get performs a GET request and decodes the JSON response.
func (c *httpClient) get(ctx context.Context, path string, result interface{}) (ResponseMetadata, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil, result)
}

// post performs a POST request with a JSON body and decodes the response.
func (c *httpClient) post(ctx context.Context, path string, body interface{}, result interface{}) (ResponseMetadata, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return ResponseMetadata{}, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}
	return c.doRequest(ctx, http.MethodPost, path, bodyReader, result)
}

// postBCS performs a POST request with a BCS body and decodes the JSON response.
func (c *httpClient) postBCS(ctx context.Context, path string, body []byte, result interface{}) (ResponseMetadata, error) {
	return c.doRequestWithContentType(ctx, http.MethodPost, path, bytes.NewReader(body), "application/x.aptos.signed_transaction+bcs", result)
}

func (c *httpClient) doRequest(ctx context.Context, method, path string, body io.Reader, result interface{}) (ResponseMetadata, error) {
	contentType := ""
	if body != nil {
		contentType = "application/json"
	}
	return c.doRequestWithContentType(ctx, method, path, body, contentType, result)
}

func (c *httpClient) doRequestWithContentType(ctx context.Context, method, path string, body io.Reader, contentType string, result interface{}) (ResponseMetadata, error) {
	url := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return ResponseMetadata{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return ResponseMetadata{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Parse response metadata from headers
	metadata := parseResponseHeaders(resp.Header)

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return metadata, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for error responses
	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			// If we can't parse the error, return a generic one
			return metadata, &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
		}
		apiErr.StatusCode = resp.StatusCode
		return metadata, &apiErr
	}

	// Decode successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return metadata, fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return metadata, nil
}

// parseResponseHeaders extracts metadata from Aptos API response headers.
func parseResponseHeaders(h http.Header) ResponseMetadata {
	return ResponseMetadata{
		ChainID:             parseHeaderUint8(h.Get("X-Aptos-Chain-Id")),
		LedgerVersion:       parseHeaderUint64(h.Get("X-Aptos-Ledger-Version")),
		LedgerOldestVersion: parseHeaderUint64(h.Get("X-Aptos-Ledger-Oldest-Version")),
		LedgerTimestampUsec: parseHeaderUint64(h.Get("X-Aptos-Ledger-TimestampUsec")),
		Epoch:               parseHeaderUint64(h.Get("X-Aptos-Epoch")),
		BlockHeight:         parseHeaderUint64(h.Get("X-Aptos-Block-Height")),
		OldestBlockHeight:   parseHeaderUint64(h.Get("X-Aptos-Oldest-Block-Height")),
		Cursor:              h.Get("X-Aptos-Cursor"),
	}
}

func parseHeaderUint8(s string) uint8 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseUint(s, 10, 8)
	return uint8(v)
}

func parseHeaderUint64(s string) uint64 {
	if s == "" {
		return 0
	}
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
