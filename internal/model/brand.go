package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Brand struct
type Brand struct {
	Pb inventories.Brand
}

// Get func
func (u *Brand) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, company_id, code, name, created_at, created_by, updated_at, updated_by 
		FROM brands WHERE id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get brand: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.Code, &u.Pb.Name, &createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get brand: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get brand: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company data")
	}

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

// GetByCode func
func (u *Brand) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, company_id, code, name, created_at, created_by, updated_at, updated_by 
		FROM brands WHERE company_id = $1 AND code = $2
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get brand by code: %v", err)
	}
	defer stmt.Close()

	var companyID string
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetCode()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.Code, &u.Pb.Name, &createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get brand by code: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get brand by code: %v", err)
	}

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

// Create Brand
func (u *Brand) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO brands (id, company_id, code, name, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert brand: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetCode(),
		u.Pb.GetName(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert brand: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}

// Update Brand
func (u *Brand) Update(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		UPDATE brands SET
		name = $1, 
		updated_at = $2, 
		updated_by= $3
		WHERE id = $4
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update brand: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetName(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update brand: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updated by: %v", err)
	}

	return nil
}

// Delete Brand
func (u *Brand) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM brands WHERE company_id = $1 AND id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete brand: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(app.Ctx("companyID")).(string), u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete brand: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Brand) ListQuery(ctx context.Context, db *sql.DB, in *inventories.Pagination) (string, []interface{}, *inventories.PaginationResponse, error) {
	var paginationResponse inventories.PaginationResponse
	query := `SELECT id, company_id, code, name, created_at, created_by, updated_at, updated_by FROM brands`
	where := []string{"company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetSearch())
		where = append(where, fmt.Sprintf(`(name ILIKE $%d OR code ILIKE $%d)`, len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM brands`
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

	if len(in.GetOrderBy()) == 0 || !(in.GetOrderBy() == "name" || in.GetOrderBy() == "code") {
		if in == nil {
			in = &inventories.Pagination{OrderBy: "created_at"}
		} else {
			in.OrderBy = "created_at"
		}
	}

	query += ` ORDER BY ` + in.GetOrderBy() + ` ` + in.GetSort().String()

	if in.GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetLimit(), in.GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
