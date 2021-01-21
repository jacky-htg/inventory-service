package service

import (
	"context"
	"inventory-service/pb/users"
	"io"

	"inventory-service/internal/pkg/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func contextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}

func getMetadata(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	userID := md["user_id"]
	if len(userID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "user_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("userID"), userID[0])

	companyID := md["company_id"]
	if len(companyID) == 0 {
		return ctx, status.Errorf(codes.Unauthenticated, "company_id is not provided")
	}

	ctx = context.WithValue(ctx, app.Ctx("companyID"), companyID[0])

	return ctx, nil
}

func setMetadata(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{
		"user_id":    ctx.Value(app.Ctx("userID")).(string),
		"company_id": ctx.Value(app.Ctx("companyID")).(string),
	})

	return metadata.NewOutgoingContext(ctx, md)
}

func isYourBranch(
	ctx context.Context,
	userClient users.UserServiceClient,
	regionClient users.RegionServiceClient,
	branchClient users.BranchServiceClient,
	branchID string) error {
	userLogin, err := getUserLogin(ctx, userClient)
	if err != nil {
		return err
	}

	if len(userLogin.GetBranchId()) > 0 {
		if userLogin.GetBranchId() != branchID {
			return status.Error(codes.Unauthenticated, "its not your branch")
		}
	} else if len(userLogin.GetRegionId()) > 0 {
		region, err := getRegion(ctx, regionClient, users.Region{Id: userLogin.GetRegionId()})
		if err != nil {
			return err
		}
		err = checkYourBranch(region.GetBranches(), branchID)
		if err != nil {
			return err
		}
	} else {
		branches, err := getBranches(ctx, branchClient)
		if err != nil {
			return err
		}
		err = checkYourBranch(branches, branchID)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkYourBranch(branches []*users.Branch, branchID string) error {
	isYourBranch := false
	for _, branch := range branches {
		if branch.GetId() == branchID {
			isYourBranch = true
			break
		}
	}

	if !isYourBranch {
		return status.Error(codes.Unauthenticated, "its not your branch")
	}

	return nil
}

func getUserLogin(ctx context.Context, userClient users.UserServiceClient) (*users.User, error) {
	userLogin, err := userClient.View(setMetadata(ctx), &users.Id{Id: ctx.Value(app.Ctx("userID")).(string)})

	if err != nil {
		return &users.User{}, status.Errorf(codes.Internal, "Error when calling user service: %v", err)
	}

	return userLogin, nil
}

func getRegion(ctx context.Context, regionClient users.RegionServiceClient, r users.Region) (*users.Region, error) {
	region, err := regionClient.View(setMetadata(ctx), &users.Id{Id: r.GetId()})

	if err != nil {
		return &users.Region{}, status.Errorf(codes.Internal, "Error when calling region service: %v", err)
	}

	return region, nil
}

func getBranches(ctx context.Context, branchClient users.BranchServiceClient) ([]*users.Branch, error) {
	var list []*users.Branch
	var err error
	stream, err := branchClient.List(setMetadata(ctx), &users.ListBranchRequest{})
	if err != nil {
		return list, status.Errorf(codes.Internal, "Error when calling branches service: %s", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return list, status.Errorf(codes.Internal, "cannot receive %v", err)
		}
		list = append(list, resp.GetBranch())
	}
	return list, err
}

func getBranch(ctx context.Context, branchClient users.BranchServiceClient, id string) (*users.Branch, error) {
	branch, err := branchClient.View(setMetadata(ctx), &users.Id{Id: id})
	if err != nil {
		return &users.Branch{}, status.Errorf(codes.Internal, "Error when calling branch service: %v", err)
	}

	return branch, nil
}
