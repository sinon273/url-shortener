package url_client

import (
	"context"
	"fmt"

	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

func (c *Client) GetLimit(ctx context.Context, userID string) (int, error) {
	limit, err := c.userClient.GetLimit(ctx, &userv1.GetLimitRequest{UserId: userID})
	if err != nil {
		return 0, fmt.Errorf("get user limit: %w", err)
	}
	return int(limit.LinksLimit), nil
}
