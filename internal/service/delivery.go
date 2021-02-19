package service

import (
	"context"
	"database/sql"
	"inventory-service/pb/users"
	"time"

	"inventory-service/internal/model"
	"inventory-service/pb/inventories"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Delivery struct
type Delivery struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}

// Create Delivery
func (u *Delivery) Create(ctx context.Context, in *inventories.Delivery) (*inventories.Delivery, error) {
	var deliveryModel model.Delivery
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetBranchId()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid branch")
		}

		if len(in.GetSalesOrderId()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid sales order")
		}

		if in.GetDeliveryDate().IsValid() {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid date")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &deliveryModel.Pb, err
	}

	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			return &deliveryModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			return &deliveryModel.Pb, err
		}

		if len(detail.GetBarcode()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid barcode")
		}

		inventory := model.Inventory{
			BranchID: in.GetBranchId(),
			Barcode:  detail.GetBarcode(),
		}
		err = inventory.CheckBarcode(ctx, u.Db)
		if err != nil {
			return &deliveryModel.Pb, err
		}

	}

	err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &deliveryModel.Pb, err
	}

	branch, err := getBranch(ctx, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &deliveryModel.Pb, err
	}

	deliveryModel.Pb = inventories.Delivery{
		BranchId:     in.GetBranchId(),
		BranchName:   branch.GetName(),
		Code:         in.GetCode(),
		DeliveryDate: in.GetDeliveryDate(),
		SalesOrderId: in.GetSalesOrderId(),
		Remark:       in.GetRemark(),
		Details:      in.GetDetails(),
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &deliveryModel.Pb, err
	}

	err = deliveryModel.Create(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &deliveryModel.Pb, err
	}

	tx.Commit()

	return &deliveryModel.Pb, nil
}

// Update Delivery
func (u *Delivery) Update(ctx context.Context, in *inventories.Delivery) (*inventories.Delivery, error) {
	var deliveryModel model.Delivery
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		deliveryModel.Pb.Id = in.GetId()
	}

	// TODO : if any return do update will be blocked

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &deliveryModel.Pb, err
	}

	err = deliveryModel.Get(ctx, u.Db)
	if err != nil {
		return &deliveryModel.Pb, err
	}

	if len(in.GetSalesOrderId()) > 0 {
		deliveryModel.Pb.SalesOrderId = in.GetSalesOrderId()
	}

	if in.GetDeliveryDate().IsValid() {
		deliveryModel.Pb.DeliveryDate = in.GetDeliveryDate()
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &deliveryModel.Pb, status.Errorf(codes.Internal, "begin transaction: %v", err)
	}

	err = deliveryModel.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &deliveryModel.Pb, err
	}

	var newDetails []*inventories.DeliveryDetail
	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			tx.Rollback()
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &deliveryModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			tx.Rollback()
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &deliveryModel.Pb, err
		}

		if len(detail.GetId()) > 0 {
			for index, data := range deliveryModel.Pb.GetDetails() {
				if data.GetId() == detail.GetId() {
					deliveryModel.Pb.Details = append(deliveryModel.Pb.Details[:index], deliveryModel.Pb.Details[index+1:]...)
					break
				}
			}
		} else {
			if len(detail.GetBarcode()) == 0 {
				return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid barcode")
			}

			inventory := model.Inventory{
				BranchID: in.GetBranchId(),
				Barcode:  detail.GetBarcode(),
			}
			err = inventory.CheckBarcode(ctx, u.Db)
			if err != nil {
				return &deliveryModel.Pb, err
			}

			// operasi insert
			deliveryDetailModel := model.DeliveryDetail{Pb: inventories.DeliveryDetail{
				DeliveryId: deliveryModel.Pb.GetId(),
				Barcode:    detail.GetBarcode(),
				Product:    detail.GetProduct(),
				Shelve:     detail.GetShelve(),
			}}
			deliveryDetailModel.PbDelivery = inventories.Delivery{
				Id:           deliveryModel.Pb.Id,
				BranchId:     deliveryModel.Pb.BranchId,
				BranchName:   deliveryModel.Pb.BranchName,
				SalesOrderId: deliveryModel.Pb.SalesOrderId,
				Code:         deliveryModel.Pb.Code,
				DeliveryDate: deliveryModel.Pb.DeliveryDate,
				Remark:       deliveryModel.Pb.Remark,
				CreatedAt:    deliveryModel.Pb.CreatedAt,
				CreatedBy:    deliveryModel.Pb.CreatedBy,
				UpdatedAt:    deliveryModel.Pb.UpdatedAt,
				UpdatedBy:    deliveryModel.Pb.UpdatedBy,
				Details:      deliveryModel.Pb.Details,
			}
			err = deliveryDetailModel.Create(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &deliveryModel.Pb, err
			}

			newDetails = append(newDetails, &deliveryDetailModel.Pb)
		}
	}

	// delete existing detail
	for _, data := range deliveryModel.Pb.GetDetails() {
		pbDeliveryDetail := inventories.DeliveryDetail{
			DeliveryId: deliveryModel.Pb.GetId(),
			Id:         data.GetId(),
		}
		deliveryDetailModel := model.DeliveryDetail{Pb: pbDeliveryDetail}
		err = deliveryDetailModel.Delete(ctx, tx)
		if err != nil {
			tx.Rollback()
			return &deliveryModel.Pb, err
		}
	}

	tx.Commit()

	return &deliveryModel.Pb, nil
}

// View Delivery
func (u *Delivery) View(ctx context.Context, in *inventories.Id) (*inventories.Delivery, error) {
	var deliveryModel model.Delivery
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &deliveryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		deliveryModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &deliveryModel.Pb, err
	}

	err = deliveryModel.Get(ctx, u.Db)
	if err != nil {
		return &deliveryModel.Pb, err
	}

	return &deliveryModel.Pb, nil
}

// List Delivery
func (u *Delivery) List(in *inventories.ListDeliveryRequest, stream inventories.DeliveryService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	var deliveryModel model.Delivery
	query, paramQueries, paginationResponse, err := deliveryModel.ListQuery(ctx, u.Db, in)

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()
	paginationResponse.Pagination = in.GetPagination()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbDelivery inventories.Delivery
		var companyID string
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbDelivery.Id, &companyID, &pbDelivery.BranchId, &pbDelivery.BranchName,
			&pbDelivery.Code, &pbDelivery.DeliveryDate, &pbDelivery.Remark,
			&createdAt, &pbDelivery.CreatedBy, &updatedAt, &pbDelivery.UpdatedBy)
		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbDelivery.CreatedAt, err = ptypes.TimestampProto(createdAt)
		if err != nil {
			return status.Errorf(codes.Internal, "convert createdAt: %v", err)
		}

		pbDelivery.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
		if err != nil {
			return status.Errorf(codes.Internal, "convert updateddAt: %v", err)
		}

		res := &inventories.ListDeliveryResponse{
			Pagination: paginationResponse,
			Delivery:   &pbDelivery,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
