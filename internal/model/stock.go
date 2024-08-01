package model

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jacky-htg/erp-pkg/app"
	"github.com/jacky-htg/erp-proto/go/pb/inventories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stock struct
type Stock struct {
	ListInput inventories.StockListInput
	InfoInput inventories.StockInfoInput
	StockInfo inventories.StockInfo
	StockList inventories.StockList
}

// Closing Stock
func (u *Stock) Closing(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `CALL closing_stocks($1, 0, 0)`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare closing stock: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(app.Ctx("companyID")).(string))
	if err != nil {
		return status.Errorf(codes.Internal, "Exec closing stock: %v", err)
	}

	return nil
}

// List Stock
func (u *Stock) List(ctx context.Context, db *sql.DB) error {
	const productSelect string = `
	products.id, products.company_id, 
	brands.id, brands.code, brands.name,
	product_categories.id, product_categories.name,
	products.code, products.name, products.minimum_stock, 
	products.created_at, products.created_by, products.updated_at, products.updated_by,
	`

	var stockQuery string = `stock (` + ctx.Value(app.Ctx("companyID")).(string) + `, products.id)`
	if len(u.ListInput.GetBranchId()) > 0 {
		stockQuery = `stock_branch (` + ctx.Value(app.Ctx("companyID")).(string) + `, ` + u.ListInput.GetBranchId() + `, products.id)`
	}

	query := `SELECT ` + productSelect + stockQuery + ` 
		FROM products 
		JOIN brands ON products.brand_id = brands.id AND products.company_id = brands.company_id
		JOIN product_categories ON products.product_category_id = product_categories.id AND products.company_id = product_categories.company_id 
	`
	where := []string{"products.company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}
	rows, err := db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var pbStockInfo inventories.StockInfo
		var pbProduct inventories.Product
		var companyID string
		var createdAt, updatedAt time.Time
		var pbBrand inventories.Brand
		var pbProductCategory inventories.ProductCategory
		var stock int32
		err = rows.Scan(
			&pbProduct.Id, &companyID,
			&pbBrand.Id, &pbBrand.Code, &pbBrand.Name,
			&pbProductCategory.Id, &pbProductCategory.Name,
			&pbProduct.Code, &pbProduct.Name, &pbProduct.MinimumStock,
			&createdAt, &pbProduct.CreatedBy, &updatedAt, &pbProduct.UpdatedBy,
			&stock,
		)

		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbProduct.Brand = &pbBrand
		pbProduct.ProductCategory = &pbProductCategory

		pbProduct.CreatedAt = createdAt.String()
		pbProduct.UpdatedAt = updatedAt.String()
		pbStockInfo.Product = &pbProduct
		pbStockInfo.Qty = stock
		u.StockList.StockInfos = append(u.StockList.StockInfos, &pbStockInfo)
	}

	if rows.Err() != nil {
		return status.Errorf(codes.Internal, "rows error: %v", err)
	}

	return nil
}

// Info Stock
func (u *Stock) Info(ctx context.Context, db *sql.DB) error {
	const productSelect string = `
	products.id, products.company_id, 
	brands.id, brands.code, brands.name,
	product_categories.id, product_categories.name,
	products.code, products.name, products.minimum_stock, 
	products.created_at, products.created_by, products.updated_at, products.updated_by,
	`

	var stockQuery string = `stock (` + ctx.Value(app.Ctx("companyID")).(string) + `, products.id)`
	if len(u.ListInput.GetBranchId()) > 0 {
		stockQuery = `stock_branch (` + ctx.Value(app.Ctx("companyID")).(string) + `, ` + u.ListInput.GetBranchId() + `, products.id)`
	}

	query := `SELECT ` + productSelect + stockQuery + ` 
		FROM products 
		JOIN brands ON products.brand_id = brands.id AND products.company_id = brands.company_id
		JOIN product_categories ON products.product_category_id = product_categories.id AND products.company_id = product_categories.company_id 
	`
	where := []string{"products.company_id = $1 AND products.id = $2"}

	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get stock info: %v", err)
	}
	defer stmt.Close()

	var pbProduct inventories.Product
	var companyID string
	var createdAt, updatedAt time.Time
	var pbBrand inventories.Brand
	var pbProductCategory inventories.ProductCategory
	var stock int32
	err = stmt.QueryRowContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.InfoInput.ProductId).Scan(
		&pbProduct.Id, &companyID,
		&pbBrand.Id, &pbBrand.Code, &pbBrand.Name,
		&pbProductCategory.Id, &pbProductCategory.Name,
		&pbProduct.Code, &pbProduct.Name, &pbProduct.MinimumStock,
		&createdAt, &pbProduct.CreatedBy, &updatedAt, &pbProduct.UpdatedBy,
		&stock,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get stock info: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get stock info: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company data")
	}

	pbProduct.Brand = &pbBrand
	pbProduct.ProductCategory = &pbProductCategory

	pbProduct.CreatedAt = createdAt.String()
	pbProduct.UpdatedAt = updatedAt.String()

	u.StockInfo = inventories.StockInfo{
		Product: &pbProduct,
		Qty:     stock,
	}

	return nil
}
