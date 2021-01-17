package model

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ProductCategory struct
type ProductCategory struct {
	Pb inventories.ProductCategory
}

// Get func
func (u *ProductCategory) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, company_id, category_id, name, created_at, created_by, updated_at, updated_by 
		FROM product_categories WHERE id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get product category: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var pbCategory inventories.Category
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &pbCategory.Id, &u.Pb.Name, &createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get category: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get category: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company data")
	}

	u.Pb.Category = &pbCategory
	u.Pb.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		return status.Errorf(codes.Internal, "convert createdAt: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updateddAt: %v", err)
	}

	return nil
}

// Create ProductCategory
func (u *ProductCategory) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO product_categories (id, company_id, category_id, name, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert product category: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetCategory().GetId(),
		u.Pb.GetName(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert product category: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}

// Update ProductCategory
func (u *ProductCategory) Update(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		UPDATE product_categories SET
		category_id = $1, 
		name = $2, 
		updated_at = $3, 
		updated_by= $4
		WHERE id = $5
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update product category: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetCategory().GetId(),
		u.Pb.GetName(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update product category: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updated by: %v", err)
	}

	return nil
}

// Delete ProductCategory
func (u *ProductCategory) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM product_categories WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete product category: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete product category: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *ProductCategory) ListQuery(ctx context.Context, db *sql.DB, in *inventories.ListProductCategoryRequest) (string, []interface{}, *inventories.ProductCategoryPaginationResponse, error) {
	var paginationResponse inventories.ProductCategoryPaginationResponse
	query := `SELECT id, company_id, category_id, name, created_at, created_by, updated_at, updated_by 
	FROM product_categories`
	where := []string{"company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetCategoryId()) > 0 {
		paramQueries = append(paramQueries, in.GetCategoryId())
		where = append(where, fmt.Sprintf("category_id = $%d", len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`name ILIKE $%d`, len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM product_categories`
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

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "name") {
		if in.GetPagination() == nil {
			in.Pagination = &inventories.Pagination{OrderBy: "created_at"}
		} else {
			in.GetPagination().OrderBy = "created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
