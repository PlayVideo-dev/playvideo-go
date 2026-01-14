package playvideo

import (
	"context"
)

// AccountService handles account-related API operations
type AccountService struct {
	client *httpClient
}

// Get retrieves the current account information
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	var result Account
	err := s.client.get(ctx, "/account", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates the account settings
func (s *AccountService) Update(ctx context.Context, params *UpdateAccountParams) (*UpdateAccountResponse, error) {
	var result UpdateAccountResponse
	err := s.client.patch(ctx, "/account", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
