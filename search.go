package serpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Response wraps the SerpApi JSON response and provides helper methods.
type Response map[string]interface{}

// NextPageParams returns the parameters for the next page of results, if available.
func (r Response) NextPageParams() map[string]string {
	serpapiPagination, ok := r["serpapi_pagination"].(map[string]interface{})
	if !ok {
		return nil
	}
	nextLink, ok := serpapiPagination["next"].(string)
	if !ok {
		return nil
	}

	u, err := url.Parse(nextLink)
	if err != nil {
		return nil
	}

	params := make(map[string]string)
	for k, v := range u.Query() {
		// We only want the engine and other search params, not the API key or output
		if k == "api_key" || k == "output" {
			continue
		}
		if len(v) > 0 {
			params[k] = v[0]
		}
	}
	return params
}

// Search performs a standard search and returns the results.
func (c *Client) Search(ctx context.Context, params map[string]string) (Response, error) {
	return c.GetJSON(ctx, "/search", params)
}

// GetLocation retrieves supported locations.
func (c *Client) GetLocation(ctx context.Context, params map[string]string) ([]interface{}, error) {
	body, err := c.makeRequest(ctx, "/locations.json", params, "json")
	if err != nil {
		return nil, err
	}

	var result []interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

// GetAccount retrieves account information.
func (c *Client) GetAccount(ctx context.Context) (Response, error) {
	return c.GetJSON(ctx, "/account", nil)
}

// GetJSON makes a GET request and decodes the response as JSON.
func (c *Client) GetJSON(ctx context.Context, path string, params map[string]string) (Response, error) {
	body, err := c.makeRequest(ctx, path, params, "json")
	if err != nil {
		return nil, err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	if errorMsg, ok := result["error"]; ok {
		return nil, fmt.Errorf("serpapi error: %v", errorMsg)
	}

	return result, nil
}


// GetHTML makes a GET request and returns the raw HTML.
func (c *Client) GetHTML(ctx context.Context, path string, params map[string]string) (string, error) {
	body, err := c.makeRequest(ctx, path, params, "html")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) makeRequest(ctx context.Context, path string, params map[string]string, output string) ([]byte, error) {
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	if c.apiKey != "" {
		q.Set("api_key", c.apiKey)
	}
	if output != "" {
		q.Set("output", output)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp struct {
			Error string `json:"error"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.Error != "" {
			return nil, fmt.Errorf("serpapi error (%d): %s", resp.StatusCode, errResp.Error)
		}
		return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
