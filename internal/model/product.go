package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"inventory-service/internal/pkg/app"
	"inventory-service/internal/pkg/util"
	"inventory-service/pb/inventories"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Product struct
type Product struct {
	Pb             inventories.Product
	PbTransactions inventories.Transactions
}

// Get func
func (u *Product) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT products.id, products.company_id, 
			brands.id, brands.code, brands.name,
			product_categories.id, product_categories.name,
			products.code, products.name, products.minimum_stock, 
			products.created_at, products.created_by, products.updated_at, products.updated_by 
		FROM products 
		JOIN brands ON products.brand_id = brands.id AND products.company_id = brands.company_id
		JOIN product_categories ON products.product_category_id = product_categories.id AND products.company_id = product_categories.company_id 
		WHERE products.id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get product: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var createdAt, updatedAt time.Time
	var pbBrand inventories.Brand
	var pbProductCategory inventories.ProductCategory
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID,
		&pbBrand.Id, &pbBrand.Code, &pbBrand.Name,
		&pbProductCategory.Id, &pbProductCategory.Name,
		&u.Pb.Code, &u.Pb.Name, &u.Pb.MinimumStock,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get product: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get product: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company data")
	}

	u.Pb.Brand = &pbBrand
	u.Pb.ProductCategory = &pbProductCategory

	u.Pb.CreatedAt = createdAt.String()
	u.Pb.UpdatedAt = updatedAt.String()

	return nil
}

// GetByCode func
func (u *Product) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT products.id, products.company_id, 
			brands.id, brands.code, brands.name,
			product_categories.id, product_categories.name,
			products.code, products.name, products.minimum_stock, 
			products.created_at, products.created_by, products.updated_at, products.updated_by 
		FROM products 
		JOIN brands ON products.brand_id = brands.id AND products.company_id = brands.company_id
		JOIN product_categories ON products.product_category_id = product_categories.id AND products.company_id = product_categories.company_id 
		WHERE products.id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get product by code: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var createdAt, updatedAt time.Time
	var pbBrand inventories.Brand
	var pbProductCategory inventories.ProductCategory
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID,
		&pbBrand.Id, &pbBrand.Code, &pbBrand.Name,
		&pbProductCategory.Id, &pbProductCategory.Name,
		&u.Pb.Code, &u.Pb.Name, &u.Pb.MinimumStock,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get product by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get product by code: %v", err)
	}

	u.Pb.Brand = &pbBrand
	u.Pb.ProductCategory = &pbProductCategory

	u.Pb.CreatedAt = createdAt.String()
	u.Pb.UpdatedAt = updatedAt.String()

	return nil
}

// Create Product
func (u *Product) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO products (id, company_id, brand_id, product_category_id, code, name, minimum_stock, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert product: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetBrand().GetId(),
		u.Pb.GetProductCategory().GetId(),
		u.Pb.GetCode(),
		u.Pb.GetName(),
		u.Pb.GetMinimumStock(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert product: %v", err)
	}

	u.Pb.CreatedAt = now.String()
	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}

// Update Product
func (u *Product) Update(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		UPDATE products SET
		brand_id = $1,
		product_category_id = $2,
		name = $3,
		minimum_stock = $4, 
		updated_at = $5, 
		updated_by= $6
		WHERE id = $7
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update product: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetBrand().GetId(),
		u.Pb.GetProductCategory().GetId(),
		u.Pb.GetName(),
		u.Pb.GetMinimumStock(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update product: %v", err)
	}

	u.Pb.UpdatedAt = now.String()

	return nil
}

// Delete Product
func (u *Product) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM products WHERE company_id = $1 AND id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete product: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete product: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Product) ListQuery(ctx context.Context, db *sql.DB, in *inventories.ListProductRequest) (string, []interface{}, *inventories.ProductPaginationResponse, error) {
	var paginationResponse inventories.ProductPaginationResponse
	query := `
		SELECT products.id, products.company_id, 
			brands.id, brands.code, brands.name,
			product_categories.id, product_categories.name,
			products.code, products.name, products.minimum_stock, 
			products.created_at, products.created_by, products.updated_at, products.updated_by 
		FROM products 
		JOIN brands ON products.brand_id = brands.id AND products.company_id = brands.company_id
		JOIN product_categories ON products.product_category_id = product_categories.id AND products.company_id = product_categories.company_id 
	`
	where := []string{"products.company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetIds()) > 0 {
		productIds := make([]interface{}, len(in.GetIds()))
		for i, productId := range in.GetIds() {
			productIds[i] = productId
		}
		var iCond string
		paramQueries, iCond = util.ConvertWhereIn("id", paramQueries, productIds)
		where = append(where, iCond)
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, "%"+in.GetPagination().GetSearch()+"%")
		where = append(where, fmt.Sprintf(`(
			products.name ILIKE $%d OR 
			products.code ILIKE $%d OR 
			brands.name ILIKE $%d OR 
			brands.code ILIKE $%d OR 
			product_categories.name ILIKE $%d)`,
			len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM products
		JOIN brands ON products.brand_id = brands.id AND products.company_id = brands.company_id
		JOIN product_categories ON products.product_category_id = product_categories.id AND products.company_id = product_categories.company_id
		`
		if len(where) > 0 {
			qCount += " WHERE " + strings.Join(where, " AND ")
		}
		var count int
		err := db.QueryRowContext(ctx, qCount, paramQueries...).Scan(&count)
		if err != nil && err != sql.ErrNoRows {
			return query, paramQueries, &paginationResponse, status.Error(codes.Internal, err.Error())
		}

		paginationResponse.Count = uint32(count)
	}

	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "products.name" || in.GetPagination().GetOrderBy() == "products.code") {
		if in.GetPagination() == nil {
			in.Pagination = &inventories.Pagination{OrderBy: "products.created_at"}
		} else {
			in.GetPagination().OrderBy = "products.created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}

// Track Product History
func (u *Product) Track(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT
			inventories.branch_id, warehouses.branch_name, shelves.warehouse_id, warehouses.name, inventories.shelve_id, 
			shelves.code, inventories.product_id, inventories.barcode, inventories.transaction_code,
			inventories.type, inventories.transaction_date, inventories.in_out
		FROM inventories
		JOIN shelves ON inventories.shelve_id = shelves.id
		JOIN warehouses ON shelves.warehouse_id = warehouses.id 
	`
	where := []string{"inventories.company_id = $1", "inventories.product_id = $2"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string), u.Pb.Id}

	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}
	rows, err := db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var pbTransaction inventories.Transaction
		err = rows.Scan(
			&pbTransaction.BranchId, &pbTransaction.BranchName, &pbTransaction.WarehouseId, &pbTransaction.WarehouseName, &pbTransaction.ShelveId,
			&pbTransaction.ShelveCode, &pbTransaction.ProductId, *&pbTransaction.Barcode, &pbTransaction.TransactionCode,
			&pbTransaction.TransactionType, &pbTransaction.TransactionDate, &pbTransaction.IsIn,
		)

		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		u.PbTransactions.Transactions = append(u.PbTransactions.Transactions, &pbTransaction)
	}

	if rows.Err() != nil {
		return status.Errorf(codes.Internal, "rows error: %v", err)
	}

	return nil
}
