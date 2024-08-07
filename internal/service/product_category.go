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

// ProductCategory struct
type ProductCategory struct {
	Db  *sql.DB
	Log map[string]*log.Logger
	inventories.UnimplementedProductCategoryServiceServer
}

// Create ProductCategory
func (u *ProductCategory) Create(ctx context.Context, in *inventories.ProductCategory) (*inventories.ProductCategory, error) {
	var productCategoryModel model.ProductCategory
	var err error

	// basic validation
	{
		if len(in.GetName()) == 0 {
			return &productCategoryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid name")
		}

		if len(in.GetCategory().GetId()) == 0 {
			return &productCategoryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid category")
		}
	}

	// category validation
	{
		categoryModel := model.Category{}
		categoryModel.Pb.Id = in.GetCategory().GetId()
		err = categoryModel.Get(ctx, u.Db)
		if err != nil {
			return &productCategoryModel.Pb, err
		}
	}

	productCategoryModel.Pb = inventories.ProductCategory{
		Category: in.GetCategory(),
		Name:     in.GetName(),
	}
	err = productCategoryModel.Create(ctx, u.Db)
	if err != nil {
		return &productCategoryModel.Pb, err
	}

	return &productCategoryModel.Pb, nil
}

// Update ProductCategory
func (u *ProductCategory) Update(ctx context.Context, in *inventories.ProductCategory) (*inventories.ProductCategory, error) {
	var productCategoryModel model.ProductCategory
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &productCategoryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		productCategoryModel.Pb.Id = in.GetId()
	}

	err = productCategoryModel.Get(ctx, u.Db)
	if err != nil {
		return &productCategoryModel.Pb, err
	}

	if len(in.GetName()) > 0 {
		productCategoryModel.Pb.Name = in.GetName()
	}

	if len(in.GetCategory().GetId()) > 0 && in.GetCategory().GetId() != productCategoryModel.Pb.GetCategory().GetId() {
		categoryModel := model.Category{}
		categoryModel.Pb.Id = in.GetCategory().GetId()
		err = categoryModel.Get(ctx, u.Db)
		if err != nil {
			return &productCategoryModel.Pb, err
		}

		productCategoryModel.Pb.Category = in.GetCategory()
	}

	err = productCategoryModel.Update(ctx, u.Db)
	if err != nil {
		return &productCategoryModel.Pb, err
	}

	return &productCategoryModel.Pb, nil
}

// View ProductCategory
func (u *ProductCategory) View(ctx context.Context, in *inventories.Id) (*inventories.ProductCategory, error) {
	var productCategoryModel model.ProductCategory
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &productCategoryModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		productCategoryModel.Pb.Id = in.GetId()
	}

	err = productCategoryModel.Get(ctx, u.Db)
	if err != nil {
		return &productCategoryModel.Pb, err
	}

	return &productCategoryModel.Pb, nil
}

// Delete ProductCategory
func (u *ProductCategory) Delete(ctx context.Context, in *inventories.Id) (*inventories.MyBoolean, error) {
	var output inventories.MyBoolean
	output.Boolean = false

	var productCategoryModel model.ProductCategory
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &output, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		productCategoryModel.Pb.Id = in.GetId()
	}

	err = productCategoryModel.Get(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	err = productCategoryModel.Delete(ctx, u.Db)
	if err != nil {
		return &output, err
	}

	output.Boolean = true
	return &output, nil
}

// List ProductCategory
func (u *ProductCategory) List(in *inventories.ListProductCategoryRequest, stream inventories.ProductCategoryService_ListServer) error {
	ctx := stream.Context()
	var productCategoryModel model.ProductCategory
	query, paramQueries, paginationResponse, err := productCategoryModel.ListQuery(ctx, u.Db, in)
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

		var pbProductCategory inventories.ProductCategory
		var companyID string
		var pbCategory inventories.Category
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbProductCategory.Id, &companyID, &pbCategory.Id, &pbCategory.Name, &pbProductCategory.Name, &createdAt, &pbProductCategory.CreatedBy, &updatedAt, &pbProductCategory.UpdatedBy)
		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbProductCategory.Category = &pbCategory
		pbProductCategory.CreatedAt = createdAt.String()
		pbProductCategory.UpdatedAt = updatedAt.String()

		res := &inventories.ListProductCategoryResponse{
			Pagination:      paginationResponse,
			ProductCategory: &pbProductCategory,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
