package grpc

import (
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	url        string
	clientConn *grpc.ClientConn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Close() error {
	if c.clientConn != nil {
		return c.clientConn.Close()
	}
	return nil
}

func (c *Client) Connect(url string) error {
	if c.clientConn != nil && url != c.url {
		err := c.clientConn.Close()
		if err != nil {
			slog.Error("unable to close grpc connection", "url", url, "error", err)
		}
		c.clientConn = nil
		c.url = ""
	}
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	c.url = url
	c.clientConn = conn
	return nil
}
