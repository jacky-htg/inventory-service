package service

import (
	"context"
	"database/sql"
	"time"

	"inventory-service/internal/model"
	"inventory-service/pb/inventories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Shelve struct
type Shelve struct {
	Db *sql.DB
	inventories.UnimplementedShelveServiceServer
}

// Create Shelve
func (u *Shelve) Create(ctx context.Context, in *inventories.Shelve) (*inventories.Shelve, error) {
	var shelveModel model.Shelve
	var err error

	// basic validation
	{
		if len(in.GetCapacity()) == 0 {
			return &shelveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid capacity")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &shelveModel.Pb, err
	}

	// warehouse validation
	{
		warehouseModel := model.Warehouse{}
		warehouseModel.Pb.Id = in.GetWarehouse().GetId()
		err = warehouseModel.Get(ctx, u.Db)
		if err != nil {
			return &shelveModel.Pb, err
		}
	}

	// code validation
	{
		shelveModel = model.Shelve{}
		shelveModel.Pb.Code = in.GetCode()
		err = shelveModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &shelveModel.Pb, err
			}
		}

		if len(shelveModel.Pb.GetId()) > 0 {
			return &shelveModel.Pb, status.Error(codes.AlreadyExists, "code must be unique")
		}
	}

	shelveModel.Pb = inventories.Shelve{
		Capacity:  in.GetCapacity(),
		Code:      in.GetCode(),
		Warehouse: in.GetWarehouse(),
	}
	err = shelveModel.Create(ctx, u.Db)
	if err != nil {
		return &shelveModel.Pb, err
	}

	return &shelveModel.Pb, nil
}

// Update Shelve
func (u *Shelve) Update(ctx context.Context, in *inventories.Shelve) (*inventories.Shelve, error) {
	var shelveModel model.Shelve
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &shelveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		shelveModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &shelveModel.Pb, err
	}

	err = shelveModel.Get(ctx, u.Db)
	if err != nil {
		return &shelveModel.Pb, err
	}

	if len(in.GetCapacity()) > 0 {
		shelveModel.Pb.Capacity = in.GetCapacity()
	}

	err = shelveModel.Update(ctx, u.Db)
	if err != nil {
		return &shelveModel.Pb, err
	}

	return &shelveModel.Pb, nil
}

// View Shelve
func (u *Shelve) View(ctx context.Context, in *inventories.Id) (*inventories.Shelve, error) {
	var shelveModel model.Shelve
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &shelveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		shelveModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &shelveModel.Pb, err
	}

	err = shelveModel.Get(ctx, u.Db)
	if err != nil {
		return &shelveModel.Pb, err
	}

	return &shelveModel.Pb, nil
}

// Delete Shelve
func (u *Shelve) Delete(ctx context.Context, in *inventories.Id) (*inventories.MyBoolean, error) {
	var output inventories.MyBoolean
	output.Boolean = false

	var shelveModel model.Shelve
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		shelveModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	err = shelveModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = shelveModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Shelve
func (u *Shelve) List(in *inventories.ListShelveRequest, stream inventories.ShelveService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	var shelveModel model.Shelve
	query, paramQueries, paginationResponse, err := shelveModel.ListQuery(ctx, u.Db, in)

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

		var pbShelve inventories.Shelve
		var pbWarehouse inventories.Warehouse
		var createdAt, updatedAt time.Time
		err = rows.Scan(
			&pbShelve.Id, &pbWarehouse.Id, &pbShelve.Code, &pbShelve.Capacity,
			&createdAt, &pbShelve.CreatedBy, &updatedAt, &pbShelve.UpdatedBy,
		)

		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbShelve.Warehouse = &pbWarehouse
		pbShelve.CreatedAt = createdAt.String()
		pbShelve.UpdatedAt = updatedAt.String()

		res := &inventories.ListShelveResponse{
			Pagination: paginationResponse,
			Shelve:     &pbShelve,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
