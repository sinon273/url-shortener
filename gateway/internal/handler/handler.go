package gateway_handler

import gateway_client "github.com/sinon273/url-shortener/gateway/internal/client"

type Handler struct {
	client *gateway_client.Client
}

func NewHandler(client *gateway_client.Client) *Handler {
	return &Handler{client: client}
}
