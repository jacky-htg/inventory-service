package model

import (
	"context"
	"database/sql"
	"encoding/json"
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

// Delivery struct
type Delivery struct {
	Pb inventories.Delivery
}

// Get func
func (u *Delivery) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT deliveries.id, deliveries.company_id, deliveries.branch_id, deliveries.branch_name, deliveries.sales_order_id, deliveries.code, 
		deliveries.delivery_date, deliveries.remark, deliveries.created_at, deliveries.created_by, deliveries.updated_at, deliveries.updated_by,
		json_agg(DISTINCT jsonb_build_object(
			'id', delivery_details.id,
			'delivery_id', delivery_details.delivery_id,
			'product_id', delivery_details.product_id,
			'product_name', products.name,
			'product_code', products.code,
			'shelve_id', delivery_details.shelve_id,
			'shelve_code', shelves.code,
			'barcode', delivery_details.barcode
		)) as details
		FROM deliveries 
		JOIN delivery_details ON deliveries.id = delivery_details.delivery_id
		JOIN products ON delivery_details.product_id = products.id
		JOIN shelves ON delivery_details.shelve_id = shelves.id
		WHERE deliveries.id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get delivery: %v", err)
	}
	defer stmt.Close()

	var dateDelivery, createdAt, updatedAt time.Time
	var companyID, details string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.SalesOrderId, &u.Pb.Code, &dateDelivery, &u.Pb.Remark,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy, &details,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code delivery: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code delivery: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.DeliveryDate, err = ptypes.TimestampProto(dateDelivery)
	if err != nil {
		return status.Errorf(codes.Internal, "convert date: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(createdAt)
	if err != nil {
		return status.Errorf(codes.Internal, "convert createdAt: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updateddAt: %v", err)
	}

	detailDeliverys := []struct {
		ID          string
		DeliveryID  string
		ProductID   string
		ProductName string
		ProductCode string
		ShelveID    string
		ShelveCode  string
		Barcode     string
	}{}
	err = json.Unmarshal([]byte(details), &detailDeliverys)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal access: %v", err)
	}

	for _, detail := range detailDeliverys {
		u.Pb.Details = append(u.Pb.Details, &inventories.DeliveryDetail{
			Barcode: detail.Barcode,
			Id:      detail.ID,
			Product: &inventories.Product{
				Id:   detail.ProductID,
				Code: detail.ProductCode,
				Name: detail.ProductName,
			},
			DeliveryId: detail.DeliveryID,
			Shelve: &inventories.Shelve{
				Id:   detail.ShelveID,
				Code: detail.ShelveCode,
			},
		})
	}

	return nil
}

// GetByCode func
func (u *Delivery) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, branch_id, branch_name, sales_order_id, code, delivery_date, remark, created_at, created_by, updated_at, updated_by 
		FROM deliveries WHERE deliveries.code = $1 AND deliveries.company_id = $2
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get by code delivery: %v", err)
	}
	defer stmt.Close()

	var dateDelivery, createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode(), ctx.Value(app.Ctx("companyID")).(string)).Scan(
		&u.Pb.Id, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.SalesOrderId, &u.Pb.Code, &dateDelivery, &u.Pb.Remark,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code delivery: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code delivery: %v", err)
	}

	u.Pb.DeliveryDate, err = ptypes.TimestampProto(dateDelivery)
	if err != nil {
		return status.Errorf(codes.Internal, "convert date: %v", err)
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

func (u *Delivery) getCode(ctx context.Context, tx *sql.Tx) (string, error) {
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM deliveries 
			WHERE company_id = $1 AND to_char(created_at, 'YYYY-mm') = to_char(now(), 'YYYY-mm')`,
		ctx.Value(app.Ctx("companyID")).(string)).Scan(&count)

	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return fmt.Sprintf("DO%d%d%d",
		time.Now().UTC().Year(),
		int(time.Now().UTC().Month()),
		(count + 1)), nil
}

// Create Delivery
func (u *Delivery) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateDelivery, err := ptypes.Timestamp(u.Pb.GetDeliveryDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert Date: %v", err)
	}

	u.Pb.Code, err = u.getCode(ctx, tx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO deliveries (id, company_id, branch_id, branch_name, sales_order_id, code, delivery_date, remark, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert delivery: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetBranchId(),
		u.Pb.GetBranchName(),
		u.Pb.GetSalesOrderId(),
		u.Pb.GetCode(),
		dateDelivery,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert delivery: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	for _, detail := range u.Pb.GetDetails() {
		deliveryDetailModel := DeliveryDetail{}
		deliveryDetailModel.Pb = inventories.DeliveryDetail{
			DeliveryId: u.Pb.GetId(),
			Barcode:    detail.GetBarcode(),
			Product:    detail.GetProduct(),
			Shelve:     detail.GetShelve(),
		}
		deliveryDetailModel.PbDelivery = inventories.Delivery{
			Id:           u.Pb.Id,
			BranchId:     u.Pb.BranchId,
			BranchName:   u.Pb.BranchName,
			SalesOrderId: u.Pb.SalesOrderId,
			Code:         u.Pb.Code,
			DeliveryDate: u.Pb.DeliveryDate,
			Remark:       u.Pb.Remark,
			CreatedAt:    u.Pb.CreatedAt,
			CreatedBy:    u.Pb.CreatedBy,
			UpdatedAt:    u.Pb.UpdatedAt,
			UpdatedBy:    u.Pb.UpdatedBy,
		}
		err = deliveryDetailModel.Create(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update Delivery
func (u *Delivery) Update(ctx context.Context, tx *sql.Tx) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateDelivery, err := ptypes.Timestamp(u.Pb.GetDeliveryDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert delivery date: %v", err)
	}

	query := `
		UPDATE deliveries SET
		sales_order_id = $1,
		delivery_date = $2,
		remark = $3, 
		updated_at = $4, 
		updated_by= $5
		WHERE id = $6
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update delivery: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetSalesOrderId(),
		dateDelivery,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update delivery: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updated by: %v", err)
	}

	return nil
}

// Delete Delivery
func (u *Delivery) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM deliveries WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete delivery: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete delivery: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Delivery) ListQuery(ctx context.Context, db *sql.DB, in *inventories.ListDeliveryRequest) (string, []interface{}, *inventories.DeliveryPaginationResponse, error) {
	var paginationResponse inventories.DeliveryPaginationResponse
	query := `SELECT id, company_id, branch_id, branch_name, sales_order_id, code, delivery_date, remark, created_at, created_by, updated_at, updated_by FROM deliveries`

	where := []string{"company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetBranchId()) > 0 {
		paramQueries = append(paramQueries, in.GetBranchId())
		where = append(where, fmt.Sprintf(`branch_id = $%d`, len(paramQueries)))
	}

	if len(in.GetSalesOrderId()) > 0 {
		paramQueries = append(paramQueries, in.GetSalesOrderId())
		where = append(where, fmt.Sprintf(`sales_order_id = $%d`, len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, in.GetPagination().GetSearch())
		where = append(where, fmt.Sprintf(`(code ILIKE $%d OR remark ILIKE $%d)`, len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM deliveries`
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

	if len(in.GetPagination().GetOrderBy()) == 0 || !(in.GetPagination().GetOrderBy() == "code") {
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
