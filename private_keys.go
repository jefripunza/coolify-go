package coolify

import (
	"context"
	"fmt"
	"net/http"
)

// PrivateKeysService handles communication with SSH private key endpoints of the Coolify API.
type PrivateKeysService struct {
	client *Client
}

// CreatePrivateKeyRequest represents request parameters to add a new SSH Private Key.
type CreatePrivateKeyRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	PrivateKey  string  `json:"private_key"`
}

// List retrieves all SSH private keys configured in the authenticated team.
func (s *PrivateKeysService) List(ctx context.Context) ([]PrivateKey, error) {
	req, err := s.client.newRequest(http.MethodGet, "security/keys", nil)
	if err != nil {
		return nil, err
	}

	var keys []PrivateKey
	_, err = s.client.do(ctx, req, &keys)
	return keys, err
}

// Get retrieves detailed configurations of a specific SSH key (excluding raw private key file).
func (s *PrivateKeysService) Get(ctx context.Context, uuid string) (*PrivateKey, error) {
	path := fmt.Sprintf("security/keys/%s", uuid)
	req, err := s.client.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	key := new(PrivateKey)
	_, err = s.client.do(ctx, req, key)
	return key, err
}

// Create registers a new SSH Private Key inside Coolify.
func (s *PrivateKeysService) Create(ctx context.Context, reqBody CreatePrivateKeyRequest) (*PrivateKey, error) {
	req, err := s.client.newRequest(http.MethodPost, "security/keys", reqBody)
	if err != nil {
		return nil, err
	}

	key := new(PrivateKey)
	_, err = s.client.do(ctx, req, key)
	return key, err
}

// Delete removes an SSH private key by its UUID.
func (s *PrivateKeysService) Delete(ctx context.Context, uuid string) error {
	path := fmt.Sprintf("security/keys/%s", uuid)
	req, err := s.client.newRequest(http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	_, err = s.client.do(ctx, req, nil)
	return err
}
