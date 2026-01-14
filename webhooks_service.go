package playvideo

import (
	"context"
	"fmt"
)

// WebhooksService handles webhook-related API operations
type WebhooksService struct {
	client *httpClient
}

// List retrieves all webhooks
func (s *WebhooksService) List(ctx context.Context) (*WebhookListResponse, error) {
	var result WebhookListResponse
	err := s.client.get(ctx, "/webhooks", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new webhook
func (s *WebhooksService) Create(ctx context.Context, params *CreateWebhookParams) (*CreateWebhookResponse, error) {
	var result CreateWebhookResponse
	err := s.client.post(ctx, "/webhooks", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a webhook by ID with recent deliveries
func (s *WebhooksService) Get(ctx context.Context, id string) (*WebhookWithDeliveries, error) {
	var result WebhookWithDeliveries
	err := s.client.get(ctx, fmt.Sprintf("/webhooks/%s", encodeURLPath(id)), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update updates a webhook
func (s *WebhooksService) Update(ctx context.Context, id string, params *UpdateWebhookParams) (*Webhook, error) {
	var result Webhook
	err := s.client.patch(ctx, fmt.Sprintf("/webhooks/%s", encodeURLPath(id)), params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Test sends a test event to a webhook
func (s *WebhooksService) Test(ctx context.Context, id string) (*TestWebhookResponse, error) {
	var result TestWebhookResponse
	err := s.client.post(ctx, fmt.Sprintf("/webhooks/%s/test", encodeURLPath(id)), nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a webhook
func (s *WebhooksService) Delete(ctx context.Context, id string) error {
	var result map[string]string
	err := s.client.delete(ctx, fmt.Sprintf("/webhooks/%s", encodeURLPath(id)), &result)
	if err != nil {
		return err
	}
	return nil
}
