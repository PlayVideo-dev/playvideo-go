package playvideo

import (
	"context"
)

// UsageService handles usage and limits API operations
type UsageService struct {
	client *httpClient
}

// Get retrieves the current usage and plan limits
func (s *UsageService) Get(ctx context.Context) (*Usage, error) {
	var result Usage
	err := s.client.get(ctx, "/usage", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
