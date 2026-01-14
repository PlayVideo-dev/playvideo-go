package playvideo

import (
	"context"
)

// EmbedService handles embed settings API operations
type EmbedService struct {
	client *httpClient
}

// GetSettings retrieves the current embed settings
func (s *EmbedService) GetSettings(ctx context.Context) (*EmbedSettings, error) {
	var result EmbedSettings
	err := s.client.get(ctx, "/embed/settings", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// UpdateSettings updates the embed settings
func (s *EmbedService) UpdateSettings(ctx context.Context, params *UpdateEmbedSettingsParams) (*UpdateEmbedSettingsResponse, error) {
	var result UpdateEmbedSettingsResponse
	err := s.client.patch(ctx, "/embed/settings", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Sign generates a signed embed URL for a video
func (s *EmbedService) Sign(ctx context.Context, params *SignEmbedParams) (*SignEmbedResponse, error) {
	var result SignEmbedResponse
	err := s.client.post(ctx, "/embed/sign", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
