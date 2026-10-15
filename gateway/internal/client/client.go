package gateway_client

import (
	analyticsv1 "github.com/sinon273/url-shortener/shared/pkg/proto/analytics/v1"
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

type Client struct {
	userClient      userv1.UserServiceClient
	urlClient       urlv1.URLServiceClient
	analyticsClient analyticsv1.AnalyticsServiceClient
}

func NewClient(userClient userv1.UserServiceClient, urlClient urlv1.URLServiceClient, analyticsClient analyticsv1.AnalyticsServiceClient) *Client {
	return &Client{
		userClient:      userClient,
		urlClient:       urlClient,
		analyticsClient: analyticsClient,
	}
}

func (c *Client) URL() urlv1.URLServiceClient {
	return c.urlClient
}

func (c *Client) User() userv1.UserServiceClient {
	return c.userClient
}

func (c *Client) Analytics() analyticsv1.AnalyticsServiceClient {
	return c.analyticsClient
}
