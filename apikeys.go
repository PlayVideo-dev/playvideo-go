package playvideo

import (
	"context"
	"fmt"
)

// APIKeysService handles API key-related operations
type APIKeysService struct {
	client *httpClient
}

// List retrieves all API keys
func (s *APIKeysService) List(ctx context.Context) (*APIKeyListResponse, error) {
	var result APIKeyListResponse
	err := s.client.get(ctx, "/api-keys", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new API key
func (s *APIKeysService) Create(ctx context.Context, params *CreateAPIKeyParams) (*CreateAPIKeyResponse, error) {
	var result CreateAPIKeyResponse
	err := s.client.post(ctx, "/api-keys", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes an API key by ID
func (s *APIKeysService) Delete(ctx context.Context, id string) error {
	var result map[string]string
	err := s.client.delete(ctx, fmt.Sprintf("/api-keys/%s", encodeURLPath(id)), &result)
	if err != nil {
		return err
	}
	return nil
}
