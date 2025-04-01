package testclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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
		apiKey:  "test-api-key",
		baseURL: "https://api-v2-sandbox.chimoney.io/v0.2.4",
		http:    http.DefaultClient,
	}

	for _, opt := range options {
		opt(c)
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

func WithTestServer(server string) Option {
	return func(c *Client) {
		c.baseURL = server
	}
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

func (c *Client) Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	if params != nil && len(params) > 0 {
		url += "?"
		for k, v := range params {
			url += fmt.Sprintf("%s=%s&", k, v)
		}
		url = url[:len(url)-1]
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("chimoney: request failed with status %d", resp.StatusCode)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}
