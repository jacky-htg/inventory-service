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

// Update DeliveryReturn
func (u *DeliveryReturn) Update(ctx context.Context, in *inventories.DeliveryReturn) (*inventories.DeliveryReturn, error) {
	var deliveryReturnModel model.DeliveryReturn
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		deliveryReturnModel.Pb.Id = in.GetId()
	}

	// TODO : if any mutation_unit update will be blocked

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	err = deliveryReturnModel.Get(ctx, u.Db)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	if len(in.GetDelivery().GetId()) > 0 {
		deliveryReturnModel.Pb.Delivery = in.GetDelivery()
	}

	if in.GetReturnDate().IsValid() {
		deliveryReturnModel.Pb.ReturnDate = in.GetReturnDate()
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &deliveryReturnModel.Pb, status.Errorf(codes.Internal, "begin transaction: %v", err)
	}

	err = deliveryReturnModel.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &deliveryReturnModel.Pb, err
	}

	var newDetails []*inventories.DeliveryReturnDetail
	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			tx.Rollback()
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &deliveryReturnModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			tx.Rollback()
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &deliveryReturnModel.Pb, err
		}

		if len(detail.GetId()) > 0 {
			// operasi update
			deliveryReturnDetailModel := model.DeliveryReturnDetail{}
			deliveryReturnDetailModel.Pb.Id = detail.GetId()
			deliveryReturnDetailModel.Pb.DeliveryReturnId = deliveryReturnModel.Pb.GetId()
			err = deliveryReturnDetailModel.Get(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &deliveryReturnModel.Pb, err
			}

			deliveryReturnDetailModel.Pb.Product = detail.GetProduct()
			deliveryReturnDetailModel.Pb.Shelve = detail.GetShelve()
			deliveryReturnDetailModel.PbDeliveryReturn = inventories.DeliveryReturn{
				Id:         deliveryReturnModel.Pb.Id,
				BranchId:   deliveryReturnModel.Pb.BranchId,
				BranchName: deliveryReturnModel.Pb.BranchName,
				Delivery:   deliveryReturnModel.Pb.GetDelivery(),
				Code:       deliveryReturnModel.Pb.Code,
				ReturnDate: deliveryReturnModel.Pb.ReturnDate,
				Remark:     deliveryReturnModel.Pb.Remark,
				CreatedAt:  deliveryReturnModel.Pb.CreatedAt,
				CreatedBy:  deliveryReturnModel.Pb.CreatedBy,
				UpdatedAt:  deliveryReturnModel.Pb.UpdatedAt,
				UpdatedBy:  deliveryReturnModel.Pb.UpdatedBy,
				Details:    deliveryReturnModel.Pb.Details,
			}
			err = deliveryReturnDetailModel.Update(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &deliveryReturnModel.Pb, err
			}

			newDetails = append(newDetails, &deliveryReturnDetailModel.Pb)
			for index, data := range deliveryReturnModel.Pb.GetDetails() {
				if data.GetId() == detail.GetId() {
					deliveryReturnModel.Pb.Details = append(deliveryReturnModel.Pb.Details[:index], deliveryReturnModel.Pb.Details[index+1:]...)
					break
				}
			}

		} else {
			// operasi insert
			deliveryReturnDetailModel := model.DeliveryReturnDetail{Pb: inventories.DeliveryReturnDetail{
				DeliveryReturnId: deliveryReturnModel.Pb.GetId(),
				Product:          detail.GetProduct(),
				Shelve:           detail.GetShelve(),
			}}
			deliveryReturnDetailModel.PbDeliveryReturn = inventories.DeliveryReturn{
				Id:         deliveryReturnModel.Pb.Id,
				BranchId:   deliveryReturnModel.Pb.BranchId,
				BranchName: deliveryReturnModel.Pb.BranchName,
				Delivery:   deliveryReturnModel.Pb.GetDelivery(),
				Code:       deliveryReturnModel.Pb.Code,
				ReturnDate: deliveryReturnModel.Pb.ReturnDate,
				Remark:     deliveryReturnModel.Pb.Remark,
				CreatedAt:  deliveryReturnModel.Pb.CreatedAt,
				CreatedBy:  deliveryReturnModel.Pb.CreatedBy,
				UpdatedAt:  deliveryReturnModel.Pb.UpdatedAt,
				UpdatedBy:  deliveryReturnModel.Pb.UpdatedBy,
				Details:    deliveryReturnModel.Pb.Details,
			}
			err = deliveryReturnDetailModel.Create(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &deliveryReturnModel.Pb, err
			}

			newDetails = append(newDetails, &deliveryReturnDetailModel.Pb)
		}
	}

	// delete existing detail
	for _, data := range deliveryReturnModel.Pb.GetDetails() {
		deliveryReturnDetailModel := model.DeliveryReturnDetail{Pb: inventories.DeliveryReturnDetail{
			DeliveryReturnId: deliveryReturnModel.Pb.GetId(),
			Id:               data.GetId(),
		}}
		err = deliveryReturnDetailModel.Delete(ctx, tx)
		if err != nil {
			tx.Rollback()
			return &deliveryReturnModel.Pb, err
		}
	}

	tx.Commit()

	return &deliveryReturnModel.Pb, nil
}
