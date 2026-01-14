package playvideo

import (
	"context"
	"fmt"
)

// CollectionsService handles collection-related API operations
type CollectionsService struct {
	client *httpClient
}

// List retrieves all collections
func (s *CollectionsService) List(ctx context.Context) (*CollectionListResponse, error) {
	var result CollectionListResponse
	err := s.client.get(ctx, "/collections", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Create creates a new collection
func (s *CollectionsService) Create(ctx context.Context, params *CreateCollectionParams) (*Collection, error) {
	var result Collection
	err := s.client.post(ctx, "/collections", params, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Get retrieves a collection by slug
func (s *CollectionsService) Get(ctx context.Context, slug string) (*CollectionWithVideos, error) {
	var result CollectionWithVideos
	err := s.client.get(ctx, fmt.Sprintf("/collections/%s", encodeURLPath(slug)), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Delete deletes a collection by slug
func (s *CollectionsService) Delete(ctx context.Context, slug string) error {
	var result map[string]string
	err := s.client.delete(ctx, fmt.Sprintf("/collections/%s", encodeURLPath(slug)), &result)
	if err != nil {
		return err
	}
	return nil
}
