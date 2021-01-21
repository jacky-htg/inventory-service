package service

import (
	"context"
	"database/sql"
	"time"

	"inventory-service/internal/model"
	"inventory-service/pb/inventories"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Product struct
type Product struct {
	Db *sql.DB
}

// Create Product
func (u *Product) Create(ctx context.Context, in *inventories.Product) (*inventories.Product, error) {
	var productModel model.Product
	var err error

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &productModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid name")
		}

		if len(in.GetBrand().GetId()) == 0 {
			return &productModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid brand")
		}

		if len(in.GetProductCategory().GetId()) == 0 {
			return &productModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product category")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &productModel.Pb, err
	}

	// brand validation
	{
		brandModel := model.Brand{}
		brandModel.Pb.Id = in.GetBrand().GetId()
		err = brandModel.Get(ctx, u.Db)
		if err != nil {
			return &productModel.Pb, err
		}
	}

	// product category validation
	{
		productCategoryModel := model.ProductCategory{}
		productCategoryModel.Pb.Id = in.GetProductCategory().GetId()
		err = productCategoryModel.Get(ctx, u.Db)
		if err != nil {
			return &productModel.Pb, err
		}
	}

	// code validation
	{
		productModel = model.Product{}
		productModel.Pb.Code = in.GetCode()
		err = productModel.GetByCode(ctx, u.Db)
		if err != nil {
			if st, ok := status.FromError(err); ok && st.Code() != codes.NotFound {
				return &productModel.Pb, err
			}
		}

		if len(productModel.Pb.GetId()) > 0 {
			return &productModel.Pb, status.Error(codes.AlreadyExists, "code must be unique")
		}
	}

	productModel.Pb = inventories.Product{
		Brand:           in.GetBrand(),
		ProductCategory: in.GetProductCategory(),
		Name:            in.GetName(),
		Code:            in.GetCode(),
		MinimunStock:    in.GetMinimunStock(),
	}
	err = productModel.Create(ctx, u.Db)
	if err != nil {
		return &productModel.Pb, err
	}

	return &productModel.Pb, nil
}

// Update Product
func (u *Product) Update(ctx context.Context, in *inventories.Product) (*inventories.Product, error) {
	var productModel model.Product
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &productModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		productModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &productModel.Pb, err
	}

	err = productModel.Get(ctx, u.Db)
	if err != nil {
		return &productModel.Pb, err
	}

	if len(in.GetName()) > 0 {
		productModel.Pb.Name = in.GetName()
	}

	productModel.Pb.MinimunStock = in.GetMinimunStock()

	if len(in.GetBrand().GetId()) > 0 && in.GetBrand().GetId() != productModel.Pb.GetBrand().GetId() {
		brandModel := model.Brand{}
		brandModel.Pb.Id = in.GetBrand().GetId()
		err = brandModel.Get(ctx, u.Db)
		if err != nil {
			return &productModel.Pb, err
		}

		productModel.Pb.Brand = in.GetBrand()
	}

	if len(in.GetProductCategory().GetId()) > 0 && in.GetProductCategory().GetId() != productModel.Pb.GetProductCategory().GetId() {
		productCategoryModel := model.ProductCategory{}
		productCategoryModel.Pb.Id = in.GetProductCategory().GetId()
		err = productCategoryModel.Get(ctx, u.Db)
		if err != nil {
			return &productModel.Pb, err
		}

		productModel.Pb.ProductCategory = in.GetProductCategory()
	}

	err = productModel.Update(ctx, u.Db)
	if err != nil {
		return &productModel.Pb, err
	}

	return &productModel.Pb, nil
}

// View Product
func (u *Product) View(ctx context.Context, in *inventories.Id) (*inventories.Product, error) {
	var productModel model.Product
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &productModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		productModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &productModel.Pb, err
	}

	err = productModel.Get(ctx, u.Db)
	if err != nil {
		return &productModel.Pb, err
	}

	return &productModel.Pb, nil
}

// Delete Product
func (u *Product) Delete(ctx context.Context, in *inventories.Id) (*inventories.Boolean, error) {
	var output inventories.Boolean
	output.Boolean = false

	var productModel model.Product
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		productModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &output, err
	}

	err = productModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = productModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List Product
func (u *Product) List(in *inventories.ListProductRequest, stream inventories.ProductService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	var productModel model.Product
	query, paramQueries, paginationResponse, err := productModel.ListQuery(ctx, u.Db, in)

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

		var pbProduct inventories.Product
		var companyID string
		var createdAt, updatedAt time.Time
		var pbBrand inventories.Brand
		var pbProductCategory inventories.ProductCategory
		err = rows.Scan(
			&pbProduct.Id, &companyID,
			&pbBrand.Id, &pbBrand.Code, &pbBrand.Name,
			&pbProductCategory.Id, &pbProductCategory.Name,
			&pbProduct.Code, &pbProduct.Name, &pbProduct.MinimunStock,
			&createdAt, &pbProduct.CreatedBy, &updatedAt, &pbProduct.UpdatedBy,
		)

		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbProduct.Brand = &pbBrand
		pbProduct.ProductCategory = &pbProductCategory

		pbProduct.CreatedAt, err = ptypes.TimestampProto(createdAt)
		if err != nil {
			return status.Errorf(codes.Internal, "convert createdAt: %v", err)
		}

		pbProduct.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
		if err != nil {
			return status.Errorf(codes.Internal, "convert updateddAt: %v", err)
		}

		res := &inventories.ListProductResponse{
			Pagination: paginationResponse,
			Product:    &pbProduct,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}