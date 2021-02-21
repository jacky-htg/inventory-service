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

// DeliveryReturn struct
type DeliveryReturn struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}

// Create DeliveryReturn
func (u *DeliveryReturn) Create(ctx context.Context, in *inventories.DeliveryReturn) (*inventories.DeliveryReturn, error) {
	var deliveryReturnModel model.DeliveryReturn
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetBranchId()) == 0 {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid branch")
		}

		if len(in.GetDelivery().GetId()) == 0 {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid receiving")
		}

		if in.GetReturnDate().IsValid() {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid date")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			return &deliveryReturnModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			return &deliveryReturnModel.Pb, err
		}
	}

	err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	branch, err := getBranch(ctx, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	deliveryReturnModel.Pb = inventories.DeliveryReturn{
		BranchId:   in.GetBranchId(),
		BranchName: branch.GetName(),
		Code:       in.GetCode(),
		ReturnDate: in.GetReturnDate(),
		Delivery:   in.GetDelivery(),
		Remark:     in.GetRemark(),
		Details:    in.GetDetails(),
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	err = deliveryReturnModel.Create(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &deliveryReturnModel.Pb, err
	}

	tx.Commit()

	return &deliveryReturnModel.Pb, nil
}

// View DeliveryReturn
func (u *DeliveryReturn) View(ctx context.Context, in *inventories.Id) (*inventories.DeliveryReturn, error) {
	var deliveryReturnModel model.DeliveryReturn
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		deliveryReturnModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	err = deliveryReturnModel.Get(ctx, u.Db)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	return &deliveryReturnModel.Pb, nil
}
