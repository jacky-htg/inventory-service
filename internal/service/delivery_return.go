package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jacky-htg/erp-pkg/app"
	"github.com/jacky-htg/erp-proto/go/pb/inventories"
	"github.com/jacky-htg/erp-proto/go/pb/users"
	"github.com/jacky-htg/inventory-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeliveryReturn struct
type DeliveryReturn struct {
	Db           *sql.DB
	Log          map[string]*log.Logger
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
	inventories.UnimplementedDeliveryReturnServiceServer
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

		if _, err := time.Parse("2006-01-02T15:04:05.000Z", in.GetReturnDate()); err != nil {
			return &deliveryReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid date")
		}
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

	err = deliveryReturnModel.Get(ctx, u.Db)
	if err != nil {
		return &deliveryReturnModel.Pb, err
	}

	if len(in.GetDelivery().GetId()) > 0 {
		deliveryReturnModel.Pb.Delivery = in.GetDelivery()
	}

	if _, err := time.Parse("2006-01-02T15:04:05.000Z", in.GetReturnDate()); err == nil {
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

	// var newDetails []*inventories.DeliveryReturnDetail
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

			// newDetails = append(newDetails, &deliveryReturnDetailModel.Pb)
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

			// newDetails = append(newDetails, &deliveryReturnDetailModel.Pb)
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

// List DeliveryReturn
func (u *DeliveryReturn) List(in *inventories.ListDeliveryReturnRequest, stream inventories.DeliveryReturnService_ListServer) error {
	ctx := stream.Context()
	var deliveryReturnModel model.DeliveryReturn
	query, paramQueries, paginationResponse, err := deliveryReturnModel.ListQuery(ctx, u.Db, in)
	if err != nil {
		return err
	}

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()
	paginationResponse.Pagination = in.GetPagination()

	for rows.Next() {
		err := app.ContextError(ctx)
		if err != nil {
			return err
		}

		var pbDeliveryReturn inventories.DeliveryReturn
		var companyID string
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbDeliveryReturn.Id, &companyID, &pbDeliveryReturn.BranchId, &pbDeliveryReturn.BranchName,
			&pbDeliveryReturn.Code, &pbDeliveryReturn.ReturnDate, &pbDeliveryReturn.Remark,
			&createdAt, &pbDeliveryReturn.CreatedBy, &updatedAt, &pbDeliveryReturn.UpdatedBy)
		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbDeliveryReturn.CreatedAt = createdAt.String()
		pbDeliveryReturn.UpdatedAt = updatedAt.String()

		res := &inventories.ListDeliveryReturnResponse{
			Pagination:     paginationResponse,
			DeliveryReturn: &pbDeliveryReturn,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
