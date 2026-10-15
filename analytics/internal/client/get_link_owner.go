package analytics_client

import (
	"context"
	"fmt"

	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
)

func (c *Client) GetLinkOwner(ctx context.Context, shortCode string) (string, error) {
	getUrlResponse, err := c.URLServiceClient.GetURL(ctx, &urlv1.GetURLRequest{ShortCode: shortCode})
	if err != nil {
		return "", fmt.Errorf("get link owner: %w", err)
	}

	return getUrlResponse.UserId, nil
}
