package url_client

import userv1 "github.com/sinon273/url-shortener/shared/pkg/proto/user/v1"

type Client struct {
	userClient userv1.UserServiceClient
}

func NewClient(userClient userv1.UserServiceClient) *Client {
	return &Client{userClient: userClient}
}
