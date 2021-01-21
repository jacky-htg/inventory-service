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

// Shelve struct
type Shelve struct {
	Pb inventories.Shelve
}

// Get func
func (u *Shelve) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT shelves.id, warehouses.id, shelves.code, shelves.capacity, 
		shelves.created_at, shelves.created_by, shelves.updated_at, shelves.updated_by 
		FROM shelves 
		JOIN warehouses ON shelves.warehouse_id = warehouses.id 
		WHERE shelves.id = $1 AND warohouses.company_id = $2
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get shelve: %v", err)
	}
	defer stmt.Close()

	var pbWarehouse inventories.Warehouse
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetId(), ctx.Value(app.Ctx("companyID")).(string)).Scan(
		&u.Pb.Id, &pbWarehouse.Id, &u.Pb.Code, &u.Pb.Capacity,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get shelve: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get shelve: %v", err)
	}

	u.Pb.Warehouse = &pbWarehouse
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
func (u *Shelve) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT shelves.id, warehouses.id, shelves.code, shelves.capacity, 
		shelves.created_at, shelves.created_by, shelves.updated_at, shelves.updated_by 
		FROM shelves 
		JOIN warehouses ON shelves.warehouse_id = warehouses.id 
		WHERE shelves.code = $1 AND warohouses.company_id = $2
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get by code shelve: %v", err)
	}
	defer stmt.Close()

	var pbWarehouse inventories.Warehouse
	var createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode(), ctx.Value(app.Ctx("companyID")).(string)).Scan(
		&u.Pb.Id, &pbWarehouse.Id, &u.Pb.Code, &u.Pb.Capacity,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code shelve: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code shelve: %v", err)
	}

	u.Pb.Warehouse = &pbWarehouse
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

// Create Shelve
func (u *Shelve) Create(ctx context.Context, db *sql.DB) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		INSERT INTO shelves (id, warehouse_id, code, capacity, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert shelve: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetWarehouse().GetId(),
		u.Pb.GetCode(),
		u.Pb.GetCapacity(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert shelve: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	return nil
}

// Update Shelve
func (u *Shelve) Update(ctx context.Context, db *sql.DB) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)

	query := `
		UPDATE shelves SET
		capacity = $1, 
		updated_at = $2, 
		updated_by= $3
		WHERE id = $4
	`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update shelve: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetCapacity(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update shelve: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updated by: %v", err)
	}

	return nil
}

// Delete Shelve
func (u *Shelve) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM shelves WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete shelve: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete shelve: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Shelve) ListQuery(ctx context.Context, db *sql.DB, in *inventories.ListShelveRequest) (string, []interface{}, *inventories.ShelvePaginationResponse, error) {
	var paginationResponse inventories.ShelvePaginationResponse
	query := `SELECT shelves.id, warehouses.id, shelves.code, shelves.capacity, 
	shelves.created_at, shelves.created_by, shelves.updated_at, shelves.updated_by 
	FROM shelves 
	JOIN warehouses ON shelves.warehouse_id = warehouses.id `

	where := []string{"shelves.warehouse_id = $1"}
	paramQueries := []interface{}{in.GetWarehouseId()}

	if len(in.GetWarehouseId()) == 0 {
		return query, paramQueries, &paginationResponse, status.Error(codes.InvalidArgument, "Please suplay valid warehouse")
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`shelves.code ILIKE $%d`, len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM shelves JOIN warehouses ON shelves.warehouse_id = warehouses.id`
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

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "shelves.code" || in.GetPagination().GetOrderBy() == "shelves.capacity") {
		if in.GetPagination() == nil {
			in.Pagination = &inventories.Pagination{OrderBy: "shelves.created_at"}
		} else {
			in.GetPagination().OrderBy = "shelves.created_at"
		}
	}

	query += ` ORDER BY ` + in.GetPagination().GetOrderBy() + ` ` + in.GetPagination().GetSort().String()

	if in.GetPagination().GetLimit() > 0 {
		query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, (len(paramQueries) + 1), (len(paramQueries) + 2))
		paramQueries = append(paramQueries, in.GetPagination().GetLimit(), in.GetPagination().GetOffset())
	}

	return query, paramQueries, &paginationResponse, nil
}
