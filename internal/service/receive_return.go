package service

import (
	"context"
	"database/sql"
	"inventory-service/internal/model"
	"inventory-service/pb/inventories"
	"inventory-service/pb/users"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveReturn struct
type ReceiveReturn struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}

// Create ReceiveReturn
func (u *ReceiveReturn) Create(ctx context.Context, in *inventories.ReceiveReturn) (*inventories.ReceiveReturn, error) {
	var receiveReturnModel model.ReceiveReturn
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetBranchId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid branch")
		}

		if len(in.GetReceive().GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid receiving")
		}

		if in.GetReturnDate().IsValid() {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid date")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			return &receiveReturnModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			return &receiveReturnModel.Pb, err
		}
	}

	err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	branch, err := getBranch(ctx, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	receiveReturnModel.Pb = inventories.ReceiveReturn{
		BranchId:   in.GetBranchId(),
		BranchName: branch.GetName(),
		Code:       in.GetCode(),
		ReturnDate: in.GetReturnDate(),
		Receive:    in.GetReceive(),
		Remark:     in.GetRemark(),
		Details:    in.GetDetails(),
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	err = receiveReturnModel.Create(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &receiveReturnModel.Pb, err
	}

	tx.Commit()

	return &receiveReturnModel.Pb, nil
}

// View ReceiveReturn
func (u *ReceiveReturn) View(ctx context.Context, in *inventories.Id) (*inventories.ReceiveReturn, error) {
	var receiveReturnModel model.ReceiveReturn
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		receiveReturnModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	err = receiveReturnModel.Get(ctx, u.Db)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	return &receiveReturnModel.Pb, nil
}
