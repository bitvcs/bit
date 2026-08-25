package server

import (
	"context"
	"time"

	"github.com/nipalab/nipa/internal/domain"
	"github.com/nipalab/nipa/internal/grpc/pb"
	"github.com/nipalab/nipa/internal/snow"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/typ.v4/slices"
)

func (n *nipaServer) GetListBranch(ctx context.Context, req *pb.GetListBranchRequest) (*pb.GetListBranchResponse, error) {
	projectID, err := snow.ParseBase36(req.Context.ProjectId)
	if err != nil {
		return nil, handleError(err)
	}
	var lastUpdate *time.Time
	if req.LastUpdatedAt != nil {
		t := req.LastUpdatedAt.AsTime()
		lastUpdate = &t
	}
	var lastID snow.ID
	if req.LastId != nil {
		lastID, err = snow.ParseBase36(*req.LastId)
		if err != nil {
			return nil, handleError(err)
		}
	}
	branches, err := n.uc.Branch().ListBranches(ctx, projectID, int(req.Limit), lastUpdate, lastID)
	if err != nil {
		return nil, handleError(err)
	}

	return &pb.GetListBranchResponse{
		Branches: slices.Map(branches, func(b *domain.Branch) *pb.Branch {
			return domainBranchToPB(b)
		}),
	}, nil
}

func (n *nipaServer) GetBranch(ctx context.Context, req *pb.GetBranchRequest) (*pb.GetBranchResponse, error) {
	projectID, err := snow.ParseBase36(req.Context.ProjectId)
	if err != nil {
		return nil, handleError(err)
	}
	branchID, err := snow.ParseBase36(req.BranchId)
	if err != nil {
		return nil, handleError(err)
	}
	branch, err := n.uc.Branch().GetByProjectIDAndID(ctx, projectID, branchID)
	if err != nil {
		return nil, handleError(err)
	}
	return &pb.GetBranchResponse{
		Branch: domainBranchToPB(branch),
	}, nil
}

func domainBranchToPB(branch *domain.Branch) *pb.Branch {
	return &pb.Branch{
		Id:        branch.ID.Base36(),
		Name:      branch.Name,
		Protected: branch.Protected,
		CommitId:  branch.CommitID.Base36(),
		UpdatedAt: timestamppb.New(branch.UpdatedAt),
		CreatedAt: timestamppb.New(branch.CreatedAt),
	}
}
