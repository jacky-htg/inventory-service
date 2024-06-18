package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Warehouse struct
type Warehouse struct {
	Pb inventories.Warehouse
}

// Get func
func (u *Warehouse) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, company_id, branch_id, branch_name, code, name, pic_name, pic_phone, 
			created_at, created_by, updated_at, updated_by 
		FROM warehouses WHERE id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get warehouse: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.Code, &u.Pb.Name,
		&u.Pb.PicName, &u.Pb.PicPhone, &createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get warehouse: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get warehouse: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company data")
	}

	u.Pb.CreatedAt = createdAt.String()
	u.Pb.UpdatedAt = updatedAt.String()

	return nil
}

// GetByCode func
func (u *Warehouse) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, company_id, branch_id, branch_name, code, name, pic_name, pic_phone, 
			created_at, created_by, updated_at, updated_by 
		FROM warehouses WHERE company_id = $1 AND code = $2
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get warehouse by code: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetCode()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.Code, &u.Pb.Name,
		&u.Pb.PicName, &u.Pb.PicPhone, &createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get warehouse by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get warehouse by code: %v", err)
	}

	u.Pb.CreatedAt = createdAt.String()
	u.Pb.UpdatedAt = updatedAt.String()

	return nil
}

// Create Warehouse
func (u *Warehouse) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO warehouses (id, company_id, branch_id, branch_name, code, name, pic_name, pic_phone, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert warehouse: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetBranchId(),
		u.Pb.GetBranchName(),
		u.Pb.GetCode(),
		u.Pb.GetName(),
		u.Pb.GetPicName(),
		u.Pb.GetPicPhone(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert warehouse: %v", err)
	}

	u.Pb.CreatedAt = now.String()
	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}

// Update Warehouse
func (u *Warehouse) Update(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		UPDATE warehouses SET
		name = $1,
		pic_name = $2,
		pic_phone = $3, 
		updated_at = $4, 
		updated_by= $5
		WHERE id = $6
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update warehouse: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetName(),
		u.Pb.GetPicName(),
		u.Pb.GetPicPhone(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update warehouse: %v", err)
	}

	u.Pb.UpdatedAt = now.String()

	return nil
}

// Delete Warehouse
func (u *Warehouse) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM warehouses WHERE company_id = $1 AND id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete warehouse: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete warehouse: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Warehouse) ListQuery(ctx context.Context, db *sql.DB, in *inventories.ListWarehouseRequest) (string, []interface{}, *inventories.WarehousePaginationResponse, error) {
	var paginationResponse inventories.WarehousePaginationResponse
	query := `SELECT id, company_id, branch_id, branch_name, code, name, pic_name, pic_phone, 
			created_at, created_by, updated_at, updated_by 
		FROM warehouses`
	where := []string{"company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetBranchId()) > 0 {
		paramQueries = append(paramQueries, in.GetBranchId())
		where = append(where, fmt.Sprintf("branch_id = $%d", len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`(name ILIKE $%d OR code ILIKE $%d OR pic_name ILIKE $%d OR pic_phone ILIKE $%d)`,
			len(paramQueries), len(paramQueries), len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM warehouses`
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

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "name" ||
		in.GetPagination().GetOrderBy() == "code" || in.GetPagination().GetOrderBy() == "pic_name" ||
		in.GetPagination().GetOrderBy() == "pic_phone") {
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
