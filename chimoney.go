// Package chimoney provides a Go client for the Chimoney API
package chimoney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chimoney/chimoney-go/modules/account"
	"github.com/chimoney/chimoney-go/modules/info"
	"github.com/chimoney/chimoney-go/modules/mobilemoney"
	"github.com/chimoney/chimoney-go/modules/payouts"
	"github.com/chimoney/chimoney-go/modules/redeem"
	"github.com/chimoney/chimoney-go/modules/subaccount"
	"github.com/chimoney/chimoney-go/modules/wallet"
)

type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client

	Account     *account.Account
	Info        *info.Info
	MobileMoney *mobilemoney.MobileMoney
	Payouts     *payouts.Payouts
	Redeem      *redeem.Redeem
	SubAccount  *subaccount.SubAccount
	Wallet      *wallet.Wallet
}

type Option func(*Client)

func New(options ...Option) *Client {
	c := &Client{
		apiKey:  os.Getenv("CHIMONEY_API_KEY"),
		baseURL: "https://api.chimoney.io/v0.2.4",
		http:    http.DefaultClient,
	}

	for _, opt := range options {
		opt(c)
	}

	if c.apiKey == "" {
		panic("chimoney: API key is required")
	}

	c.Account = account.New(c)
	c.Info = info.New(c)
	c.MobileMoney = mobilemoney.New(c)
	c.Payouts = payouts.New(c)
	c.Redeem = redeem.New(c)
	c.SubAccount = subaccount.New(c)
	c.Wallet = wallet.New(c)

	return c
}

func WithAPIKey(key string) Option {
	return func(c *Client) {
		c.apiKey = key
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.http = client
	}
}

func WithSandbox(enabled bool) Option {
	return func(c *Client) {
		if enabled {
			c.baseURL = "https://api-v2-sandbox.chimoney.io/v0.2.4"
		}
	}
}

func (c *Client) Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}


	// Create request with base URL and path
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reqBody)
	if err != nil {
		return err
	}

	// Add query parameters
	if params != nil {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}


	// Add query parameters if body is a map
	if params, ok := body.(map[string]interface{}); ok && method == "GET" {
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, fmt.Sprint(v))
		}
		req.URL.RawQuery = q.Encode()
	}


	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("chimoney: request failed with status %d", resp.StatusCode)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("chimoney: failed to decode response: %v", err)
		}
	}

	return nil
}
