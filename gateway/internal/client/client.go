package gateway_client

import (
	urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"
	userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"
)

type Client struct {
	userClient userv1.UserServiceClient
	urlClient  urlv1.URLServiceClient
}

func NewClient(userClient userv1.UserServiceClient, urlClient urlv1.URLServiceClient) *Client {
	return &Client{
		userClient: userClient,
		urlClient:  urlClient,
	}
}

func (c *Client) URL() urlv1.URLServiceClient {
	return c.urlClient
}

func (c *Client) User() userv1.UserServiceClient {
	return c.userClient
}
