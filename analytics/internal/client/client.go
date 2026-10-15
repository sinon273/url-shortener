package analytics_client

import urlv1 "github.com/sinon273/url-shortener/shared/pkg/proto/url/v1"

type Client struct {
	urlv1.URLServiceClient
}

func NewClient(URLServiceClient urlv1.URLServiceClient) *Client {
	return &Client{URLServiceClient: URLServiceClient}
}
