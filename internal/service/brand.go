package service

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jacky-htg/erp-pkg/app"
	"github.com/jacky-htg/erp-proto/go/pb/inventories"
	"github.com/jacky-htg/inventory-service/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Brand struct
type Brand struct {
	Db  *sql.DB
	Log map[string]*log.Logger
	inventories.UnimplementedBrandServiceServer
}

// Create Brand
func (u *Brand) Create(ctx context.Context, in *inventories.Brand) (*inventories.Brand, error) {
	var brandModel model.Brand
	var err error

	// basic validation
	{
		if len(in.GetName()) == 0 {
			err = status.Error(codes.InvalidArgument, "Please supply valid name")
			u.Log["error"].Println(err)
			return &brandModel.Pb, err
		}
	}

	// code validation
	{
		if len(in.GetCode()) == 0 {
			err = status.Error(codes.InvalidArgument, "Please supply valid code")
			u.Log["error"].Println(err)
			return &brandModel.Pb, err
		}

		brandModel = model.Brand{}
		brandModel.Pb.Code = in.GetCode()
		err = brandModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); !ok && st.Code() != codes.NotFound {
				return &brandModel.Pb, err
			}
		}

		if len(brandModel.Pb.GetId()) > 0 {
			err = status.Error(codes.AlreadyExists, "code must be unique")
			u.Log["error"].Println(err)
			return &brandModel.Pb, err
		}
	}

	brandModel.Pb = inventories.Brand{
		Code: in.GetCode(),
		Name: in.GetName(),
	}
	err = brandModel.Create(ctx, u.Db)
	if err != nil {
		return &brandModel.Pb, err
	}

	return &brandModel.Pb, nil
}

// Update Brand
func (u *Brand) Update(ctx context.Context, in *inventories.Brand) (*inventories.Brand, error) {
	var brandModel model.Brand
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			err = status.Error(codes.InvalidArgument, "Please supply valid id")
			u.Log["error"].Println(err)
			return &brandModel.Pb, err
		}
		brandModel.Pb.Id = in.GetId()
	}

	err = brandModel.Get(ctx, u.Db)
	if err != nil {
		return &brandModel.Pb, err
	}

	if len(in.GetName()) > 0 {
		brandModel.Pb.Name = in.GetName()
	}

	err = brandModel.Update(ctx, u.Db)
	if err != nil {
		return &brandModel.Pb, err
	}

	return &brandModel.Pb, nil
}

// View Brand
func (u *Brand) View(ctx context.Context, in *inventories.Id) (*inventories.Brand, error) {
	var brandModel model.Brand
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			err = status.Error(codes.InvalidArgument, "Please supply valid id")
			u.Log["error"].Println(err)
			return &brandModel.Pb, err
		}
		brandModel.Pb.Id = in.GetId()
	}

	err = brandModel.Get(ctx, u.Db)
	if err != nil {
		return &brandModel.Pb, err
	}

	return &brandModel.Pb, nil
}

// Delete Brand
func (u *Brand) Delete(ctx context.Context, in *inventories.Id) (*inventories.MyBoolean, error) {
	var output inventories.MyBoolean
	output.Boolean = false

	var brandModel model.Brand
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			err = status.Error(codes.InvalidArgument, "Please supply valid id")
			u.Log["error"].Println(err)
			return &output, err
		}
		brandModel.Pb.Id = in.GetId()
	}

	err = brandModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = brandModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Brand
func (u *Brand) List(in *inventories.Pagination, stream inventories.BrandService_ListServer) error {
	ctx := stream.Context()
	var brandModel model.Brand
	query, paramQueries, paginationResponse, err := brandModel.ListQuery(ctx, u.Db, in)
	if err != nil {
		return err
	}

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		u.Log["error"].Println(err)
		return err
	}
	defer rows.Close()
	paginationResponse.Pagination = in

	for rows.Next() {
		err := app.ContextError(ctx)
		if err != nil {
			u.Log["error"].Println(err)
			return err
		}

		var pbBrand inventories.Brand
		var companyID string
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbBrand.Id, &companyID, &pbBrand.Code, &pbBrand.Name, &createdAt, &pbBrand.CreatedBy, &updatedAt, &pbBrand.UpdatedBy)
		if err != nil {
			status.Errorf(codes.Internal, "scan data: %v", err)
			u.Log["error"].Println(err)
			return err
		}

		pbBrand.CreatedAt = createdAt.String()
		pbBrand.UpdatedAt = updatedAt.String()

		res := &inventories.ListBrandResponse{
			Pagination: paginationResponse,
			Brand:      &pbBrand,
		}

		err = stream.Send(res)
		if err != nil {
			err = status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
			u.Log["error"].Println(err)
			return err
		}
	}
	return nil
}
