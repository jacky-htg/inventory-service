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

// Warehouse struct
type Warehouse struct {
	Db           *sql.DB
	Log          map[string]*log.Logger
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
	inventories.UnimplementedWarehouseServiceServer
}

// Create Warehouse
func (u *Warehouse) Create(ctx context.Context, in *inventories.Warehouse) (*inventories.Warehouse, error) {
	var warehouseModel model.Warehouse
	var err error

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &warehouseModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid name")
		}

		if len(in.GetBranchId()) == 0 {
			return &warehouseModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid branch")
		}

		if len(in.GetPicName()) == 0 {
			return &warehouseModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid PIC name")
		}

		if len(in.GetPicPhone()) == 0 {
			return &warehouseModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid PIC phone")
		}
	}

	// code validation
	{
		warehouseModel = model.Warehouse{}
		warehouseModel.Pb.Code = in.GetCode()
		err = warehouseModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &warehouseModel.Pb, err
			}
		}

		if len(warehouseModel.Pb.GetId()) > 0 {
			return &warehouseModel.Pb, status.Error(codes.AlreadyExists, "code must be unique")
		}
	}

	err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &warehouseModel.Pb, err
	}

	branch, err := getBranch(ctx, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &warehouseModel.Pb, err
	}

	warehouseModel.Pb = inventories.Warehouse{
		BranchId:   in.GetBranchId(),
		BranchName: branch.GetName(),
		Code:       in.GetCode(),
		Name:       in.GetName(),
		PicName:    in.GetPicName(),
		PicPhone:   in.GetPicPhone(),
	}
	err = warehouseModel.Create(ctx, u.Db)
	if err != nil {
		return &warehouseModel.Pb, err
	}

	return &warehouseModel.Pb, nil
}

// Update Warehouse
func (u *Warehouse) Update(ctx context.Context, in *inventories.Warehouse) (*inventories.Warehouse, error) {
	var warehouseModel model.Warehouse
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &warehouseModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		warehouseModel.Pb.Id = in.GetId()
	}

	err = warehouseModel.Get(ctx, u.Db)
	if err != nil {
		return &warehouseModel.Pb, err
	}

	if len(in.GetName()) > 0 {
		warehouseModel.Pb.Name = in.GetName()
	}

	if len(in.GetPicName()) > 0 {
		warehouseModel.Pb.PicName = in.GetPicName()
	}

	if len(in.GetPicPhone()) > 0 {
		warehouseModel.Pb.PicPhone = in.GetPicPhone()
	}

	err = warehouseModel.Update(ctx, u.Db)
	if err != nil {
		return &warehouseModel.Pb, err
	}

	return &warehouseModel.Pb, nil
}

// View Warehouse
func (u *Warehouse) View(ctx context.Context, in *inventories.Id) (*inventories.Warehouse, error) {
	var warehouseModel model.Warehouse
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &warehouseModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		warehouseModel.Pb.Id = in.GetId()
	}

	err = warehouseModel.Get(ctx, u.Db)
	if err != nil {
		return &warehouseModel.Pb, err
	}

	return &warehouseModel.Pb, nil
}

// Delete Warehouse
func (u *Warehouse) Delete(ctx context.Context, in *inventories.Id) (*inventories.MyBoolean, error) {
	var output inventories.MyBoolean
	output.Boolean = false

	var warehouseModel model.Warehouse
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		warehouseModel.Pb.Id = in.GetId()
	}

	err = warehouseModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = warehouseModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Warehouse
func (u *Warehouse) List(in *inventories.ListWarehouseRequest, stream inventories.WarehouseService_ListServer) error {
	ctx := stream.Context()
	var warehouseModel model.Warehouse
	query, paramQueries, paginationResponse, err := warehouseModel.ListQuery(ctx, u.Db, in)
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

		var pbWarehouse inventories.Warehouse
		var companyID string
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbWarehouse.Id, &companyID, &pbWarehouse.BranchId, &pbWarehouse.BranchName,
			&pbWarehouse.Code, &pbWarehouse.Name, &pbWarehouse.PicName, &pbWarehouse.PicPhone,
			&createdAt, &pbWarehouse.CreatedBy, &updatedAt, &pbWarehouse.UpdatedBy)
		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbWarehouse.CreatedAt = createdAt.String()
		pbWarehouse.UpdatedAt = updatedAt.String()

		res := &inventories.ListWarehouseResponse{
			Pagination: paginationResponse,
			Warehouse:  &pbWarehouse,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
