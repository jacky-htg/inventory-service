package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveReturn struct
type ReceiveReturn struct {
	Pb inventories.ReceiveReturn
}

// Get func
func (u *ReceiveReturn) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT receive_returns.id, receive_returns.company_id, receive_returns.branch_id, receive_returns.branch_name, receive_returns.receiving_id, receive_returns.code, 
		receive_returns.return_date, receive_returns.remark, receive_returns.created_at, receive_returns.created_by, receive_returns.updated_at, receive_returns.updated_by,
		json_agg(DISTINCT jsonb_build_object(
			'id', receive_return_details.id,
			'receive_return_id', receive_return_details.receive_return_id,
			'product_id', receive_return_details.product_id,
			'product_name', products.name,
			'product_code', products.code,
			'shelve_id', receive_return_details.shelve_id,
			'shelve_code', shelves.code
		)) as details
		FROM receive_returns 
		JOIN receive_return_details ON receive_returns.id = receive_return_details.receive_return_id
		JOIN products ON receive_return_details.product_id = products.id
		JOIN shelves ON receive_return_details.shelve_id = shelves.id
		WHERE receive_returns.id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get receive return: %v", err)
	}
	defer stmt.Close()

	var dateReturn, createdAt, updatedAt time.Time
	var companyID, details string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.Receive.Id, &u.Pb.Code, &dateReturn, &u.Pb.Remark,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy, &details,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code receive return: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code receive return: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.ReturnDate, err = ptypes.TimestampProto(dateReturn)
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

	detailReceiveReturns := []struct {
		ID              string
		ReceiveReturnID string
		ProductID       string
		ProductName     string
		ProductCode     string
		ShelveID        string
		ShelveCode      string
	}{}
	err = json.Unmarshal([]byte(details), &detailReceiveReturns)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal detailReceiveReturns: %v", err)
	}

	for _, detail := range detailReceiveReturns {
		u.Pb.Details = append(u.Pb.Details, &inventories.ReceiveReturnDetail{
			Id: detail.ID,
			Product: &inventories.Product{
				Id:   detail.ProductID,
				Code: detail.ProductCode,
				Name: detail.ProductName,
			},
			ReceiveReturnId: detail.ReceiveReturnID,
			Shelve: &inventories.Shelve{
				Id:   detail.ShelveID,
				Code: detail.ShelveCode,
			},
		})
	}

	return nil
}

// Create ReceiveReturn
func (u *ReceiveReturn) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateReturn, err := ptypes.Timestamp(u.Pb.GetReturnDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert Date: %v", err)
	}

	u.Pb.Code, err = u.getCode(ctx, tx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO receive_returns (id, company_id, branch_id, branch_name, receive_id, code, return_date, remark, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert receive return: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetBranchId(),
		u.Pb.GetBranchName(),
		u.Pb.GetReceive().GetId(),
		u.Pb.GetCode(),
		dateReturn,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert receive: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	for _, detail := range u.Pb.GetDetails() {
		receiveReturnDetailModel := ReceiveReturnDetail{}
		receiveReturnDetailModel.Pb = inventories.ReceiveReturnDetail{
			ReceiveReturnId: u.Pb.GetId(),
			Product:         detail.GetProduct(),
			Shelve:          detail.GetShelve(),
		}
		receiveReturnDetailModel.PbReceiveReturn = inventories.ReceiveReturn{
			Id:         u.Pb.Id,
			BranchId:   u.Pb.BranchId,
			BranchName: u.Pb.BranchName,
			Receive:    u.Pb.Receive,
			Code:       u.Pb.Code,
			ReturnDate: u.Pb.ReturnDate,
			Remark:     u.Pb.Remark,
			CreatedAt:  u.Pb.CreatedAt,
			CreatedBy:  u.Pb.CreatedBy,
			UpdatedAt:  u.Pb.UpdatedAt,
			UpdatedBy:  u.Pb.UpdatedBy,
		}
		err = receiveReturnDetailModel.Create(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *ReceiveReturn) getCode(ctx context.Context, tx *sql.Tx) (string, error) {
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM receive_returns 
			WHERE company_id = $1 AND to_char(created_at, 'YYYY-mm') = to_char(now(), 'YYYY-mm')`,
		ctx.Value(app.Ctx("companyID")).(string)).Scan(&count)

	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return fmt.Sprintf("RR%d%d%d",
		time.Now().UTC().Year(),
		int(time.Now().UTC().Month()),
		(count + 1)), nil
}
