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

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Receive struct
type Receive struct {
	Pb inventories.Receive
}

// Get func
func (u *Receive) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT receives.id, receives.company_id, receives.branch_id, receives.branch_name, receives.purchase_id, receives.code, 
		receives.receive_date, receives.remark, receives.created_at, receives.created_by, receives.updated_at, receives.updated_by,
		json_agg(DISTINCT jsonb_build_object(
			'id', receive_details.id,
			'receive_id', receive_details.receive_id,
			'product_id', receive_details.product_id,
			'product_name', products.name,
			'product_code', products.code,
			'shelve_id', receive_details.shelve_id,
			'shelve_code', shelves.code,
			'expired_date', receive_details.expired_date
		)) as details
		FROM receives 
		JOIN receive_details ON receives.id = receive_details.receive_id
		JOIN products ON receive_details.product_id = products.id
		JOIN shelves ON receive_details.shelve_id = shelves.id
		WHERE receives.id = $1
		GROUP BY receives.id
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get receive: %v", err)
	}
	defer stmt.Close()

	var dateReceive, createdAt, updatedAt time.Time
	var companyID, details string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.PurchaseId, &u.Pb.Code, &dateReceive, &u.Pb.Remark,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy, &details,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code receive: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code receive: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.ReceiveDate = dateReceive.String()
	u.Pb.CreatedAt = createdAt.String()
	u.Pb.UpdatedAt = updatedAt.String()

	detailReceives := []struct {
		ID          string
		ReceiveID   string `json:"receive_id"`
		ProductID   string `json:"product_id"`
		ProductName string `json:"product_name"`
		ProductCode string `json:"product_code"`
		ShelveID    string `json:"shelve_id"`
		ShelveCode  string `json:"shelve_code"`
		ExpiredDate string `json:"expired_date"`
	}{}
	err = json.Unmarshal([]byte(details), &detailReceives)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal access: %v", err)
	}

	for _, detail := range detailReceives {
		u.Pb.Details = append(u.Pb.Details, &inventories.ReceiveDetail{
			ExpiredDate: detail.ExpiredDate,
			Id:          detail.ID,
			Product: &inventories.Product{
				Id:   detail.ProductID,
				Code: detail.ProductCode,
				Name: detail.ProductName,
			},
			ReceiveId: detail.ReceiveID,
			Shelve: &inventories.Shelve{
				Id:   detail.ShelveID,
				Code: detail.ShelveCode,
			},
		})
	}

	return nil
}

// GetByCode func
func (u *Receive) GetByCode(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT id, branch_id, branch_name, purchase_id, code, receive_date, remark, created_at, created_by, updated_at, updated_by 
		FROM receives WHERE receives.code = $1 AND receives.company_id = $2
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get by code receive: %v", err)
	}
	defer stmt.Close()

	var dateReceive, createdAt, updatedAt time.Time
	err = stmt.QueryRowContext(ctx, u.Pb.GetCode(), ctx.Value(app.Ctx("companyID")).(string)).Scan(
		&u.Pb.Id, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.PurchaseId, &u.Pb.Code, &dateReceive, &u.Pb.Remark,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code receive: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code receive: %v", err)
	}

	u.Pb.ReceiveDate = dateReceive.String()
	u.Pb.CreatedAt = createdAt.String()
	u.Pb.UpdatedAt = updatedAt.String()

	return nil
}

func (u *Receive) getCode(ctx context.Context, tx *sql.Tx) (string, error) {
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM receives 
			WHERE company_id = $1 AND to_char(created_at, 'YYYY-mm') = to_char(now(), 'YYYY-mm')`,
		ctx.Value(app.Ctx("companyID")).(string)).Scan(&count)

	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return fmt.Sprintf("GR%d%02d%05d",
		time.Now().UTC().Year(),
		int(time.Now().UTC().Month()),
		(count + 1)), nil
}

// Create Receive
func (u *Receive) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateReceive, err := time.Parse("2006-01-02T15:04:05.000Z", u.Pb.GetReceiveDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert Date: %v", err)
	}

	u.Pb.Code, err = u.getCode(ctx, tx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO receives (id, company_id, branch_id, branch_name, purchase_id, code, receive_date, remark, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert receive: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetBranchId(),
		u.Pb.GetBranchName(),
		u.Pb.GetPurchaseId(),
		u.Pb.GetCode(),
		dateReceive,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert receive: %v", err)
	}

	u.Pb.CreatedAt = now.String()
	u.Pb.UpdatedAt = u.Pb.CreatedAt

	for _, detail := range u.Pb.GetDetails() {
		receiveDetailModel := ReceiveDetail{}
		receiveDetailModel.Pb = inventories.ReceiveDetail{
			ReceiveId:   u.Pb.GetId(),
			ExpiredDate: detail.GetExpiredDate(),
			Product:     detail.GetProduct(),
			Shelve:      detail.GetShelve(),
		}
		receiveDetailModel.PbReceive = inventories.Receive{
			Id:          u.Pb.Id,
			BranchId:    u.Pb.BranchId,
			BranchName:  u.Pb.BranchName,
			PurchaseId:  u.Pb.PurchaseId,
			Code:        u.Pb.Code,
			ReceiveDate: u.Pb.ReceiveDate,
			Remark:      u.Pb.Remark,
			CreatedAt:   u.Pb.CreatedAt,
			CreatedBy:   u.Pb.CreatedBy,
			UpdatedAt:   u.Pb.UpdatedAt,
			UpdatedBy:   u.Pb.UpdatedBy,
		}
		err = receiveDetailModel.Create(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update Receive
func (u *Receive) Update(ctx context.Context, tx *sql.Tx) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateReceive, err := time.Parse("2006-01-02T15:04:05.000Z", u.Pb.GetReceiveDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert receive date: %v", err)
	}

	query := `
		UPDATE receives SET
		purchase_id = $1,
		receive_date = $2,
		remark = $3, 
		updated_at = $4, 
		updated_by= $5
		WHERE id = $6
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update receive: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetPurchaseId(),
		dateReceive,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update receive: %v", err)
	}

	u.Pb.UpdatedAt = now.String()

	return nil
}

// Delete Receive
func (u *Receive) Delete(ctx context.Context, db *sql.DB) error {
	stmt, err := db.PrepareContext(ctx, `DELETE FROM receives WHERE id = $1`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete receive: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete receive: %v", err)
	}

	return nil
}

// ListQuery builder
func (u *Receive) ListQuery(ctx context.Context, db *sql.DB, in *inventories.ListReceiveRequest) (string, []interface{}, *inventories.ReceivePaginationResponse, error) {
	var paginationResponse inventories.ReceivePaginationResponse
	query := `
		SELECT 
			id, company_id, branch_id, branch_name, purchase_id, code, receive_date, remark, created_at, created_by, updated_at, updated_by FROM receives`

	where := []string{"company_id = $1"}
	paramQueries := []interface{}{ctx.Value(app.Ctx("companyID")).(string)}

	if len(in.GetBranchId()) > 0 {
		paramQueries = append(paramQueries, in.GetBranchId())
		where = append(where, fmt.Sprintf(`branch_id = $%d`, len(paramQueries)))
	}

	if len(in.GetPurchaseId()) > 0 {
		paramQueries = append(paramQueries, in.GetPurchaseId())
		where = append(where, fmt.Sprintf(`purchase_id = $%d`, len(paramQueries)))
	}

	if len(in.GetPagination().GetSearch()) > 0 {
		paramQueries = append(paramQueries, "%"+in.GetPagination().GetSearch()+"%")
		where = append(where, fmt.Sprintf(`(code ILIKE $%d OR remark ILIKE $%d)`, len(paramQueries), len(paramQueries)))
	}

	{
		qCount := `SELECT COUNT(*) FROM receives`
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
