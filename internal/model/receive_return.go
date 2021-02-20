package model

import (
	"context"
	"database/sql"
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
