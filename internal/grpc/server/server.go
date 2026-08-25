package server

import (
	"context"

	"github.com/nipalab/nipa/internal/grpc/pb"
	"github.com/nipalab/nipa/internal/usecase"
)

type usecaseContainer interface {
	Auth() *usecase.Auth
	User() *usecase.User
}

type nipaServer struct {
	pb.UnimplementedNipaServiceServer
	uc usecaseContainer
}

func New(usecaseContainer usecaseContainer) *nipaServer {
	return &nipaServer{uc: usecaseContainer}
}

func (n *nipaServer) GetTreeManifest(ctx context.Context, req *pb.GetTreeManifestRequest) (*pb.GetTreeManifestResponse, error) {
	return nil, nil
}
