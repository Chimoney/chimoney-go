package subaccount

import (
	"context"
	"encoding/json"
	"errors"
)

var (
	ErrInvalidID = errors.New("invalid sub-account ID")
)

type Client interface {
	Do(ctx context.Context, method, path string, body interface{}, v interface{}, params map[string]string) error
}

type SubAccount struct {
	client Client
}

func New(client Client) *SubAccount {
	return &SubAccount{client: client}
}

type CreateRequest struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Description string `json:"description,omitempty"`
}

type SubAccountResponse struct {
	Status  string          `json:"status"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

/**
 * This function creates a new sub-account
 * @param {string} name The name of the sub-account
 * @param {string} email The email of the sub-account
 * @param {string?} description Optional description for the sub-account
 * @returns The response from the Chimoney API
 */
func (s *SubAccount) Create(ctx context.Context, req *CreateRequest) (*SubAccountResponse, error) {
	resp := new(SubAccountResponse)
	err := s.client.Do(ctx, "POST", "/sub-account", req, resp, nil)
	return resp, err
}

/**
 * This function lists all sub-accounts
 * @returns The response from the Chimoney API
 */
func (s *SubAccount) List(ctx context.Context) (*SubAccountResponse, error) {
	resp := new(SubAccountResponse)
	err := s.client.Do(ctx, "GET", "/sub-account/list", nil, resp, nil)
	return resp, err
}

/**
 * This function deletes a sub-account
 * @param {string} id The ID of the sub-account to delete
 * @returns The response from the Chimoney API
 */
func (s *SubAccount) Delete(ctx context.Context, id string) (*SubAccountResponse, error) {
	if id == "" {
		return nil, ErrInvalidID
	}

	req := map[string]string{
		"id": id,
	}

	resp := new(SubAccountResponse)
	err := s.client.Do(ctx, "DELETE", "/sub-account", req, resp, nil)
	return resp, err
}
