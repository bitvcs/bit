package grpc

import (
	"context"
	"log/slog"

	"github.com/nipalab/nipa/internal/client/domain"
	pb "github.com/nipalab/nipa/internal/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type tokenProvider interface {
	LoginWithRefreshToken(ctx context.Context, host, refreshToken string) error
	GetToken(ctx context.Context, host string) (string, error)
}

type Client struct {
	url           string
	clientConn    *grpc.ClientConn
	tokenProvider tokenProvider
}

func NewClient(tokenProvider tokenProvider) *Client {
	return &Client{
		tokenProvider: tokenProvider,
	}
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
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(c.unaryAuthInterceptor()))
	if err != nil {
		return err
	}
	c.url = url
	c.clientConn = conn
	return nil
}

func (c *Client) unaryAuthInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if method == "/greet.NipaService/LoginWithUsernamePassword" || method == "/greet.NipaService/LoginWithRefreshToken" {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		accToken, err := c.tokenProvider.GetToken(ctx, c.url)
		if err != nil {
			return status.Error(codes.Unauthenticated, "unable to get access token")
		}
		authCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+accToken)

		err = invoker(authCtx, method, req, reply, cc, opts...)
		if status.Code(err) != codes.Unauthenticated {
			return err
		}

		err = c.tokenProvider.LoginWithRefreshToken(ctx, c.url, "")
		if err != nil {
			return err
		}

		accToken, err = c.tokenProvider.GetToken(ctx, c.url)
		if err != nil {
			return status.Error(codes.Unauthenticated, "unable to get access token after refresh")
		}
		retryCtx := metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+accToken)
		return invoker(retryCtx, method, req, reply, cc, opts...)
	}
}

func (c *Client) LoginWithUsernamePassword(ctx context.Context, host, username, password string) (*domain.LoginResult, error) {
	if err := c.Connect(host); err != nil {
		return nil, err
	}
	client := pb.NewNipaServiceClient(c.clientConn)
	res, err := client.LoginWithUsernamePassword(ctx, &pb.LoginUsernamePasswordRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return &domain.LoginResult{
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
		ExpiresIn:    int(res.GetExpiresIn()),
		Host:         host,
	}, nil
}

func (c *Client) LoginWithRefreshToken(ctx context.Context, host, refreshToken string) (*domain.LoginResult, error) {
	if err := c.Connect(host); err != nil {
		return nil, err
	}
	client := pb.NewNipaServiceClient(c.clientConn)
	res, err := client.LoginWithRefreshToken(ctx, &pb.LoginWithRefreshRequest{
		RefreshToken: refreshToken,
	})
	if err != nil {
		return nil, err
	}
	return &domain.LoginResult{
		AccessToken:  res.GetAccessToken(),
		RefreshToken: res.GetRefreshToken(),
		ExpiresIn:    int(res.GetExpiresIn()),
		Host:         host,
	}, nil
}
