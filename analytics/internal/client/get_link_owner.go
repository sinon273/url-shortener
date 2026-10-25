package analytics_client

import (
	"context"
	"fmt"

	analytics_domain "github.com/sinon273/url-shortener/analytics/internal/domain"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Client) GetLinkOwner(ctx context.Context, shortCode string) (string, error) {
	getUrlResponse, err := c.GetURL(ctx, &urlv1.GetURLRequest{ShortCode: shortCode})
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return "", analytics_domain.ErrNotFound
		case codes.Unavailable:
			return "", analytics_domain.ErrUnavailable
		default:
			return "", fmt.Errorf("get link owner: %w", err)
		}
	}
	return getUrlResponse.UserId, nil
}
