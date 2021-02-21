package model

import (
	"context"
	"database/sql"
	"encoding/json"
	"inventory-service/internal/pkg/app"
	"inventory-service/internal/pkg/util"
	"inventory-service/pb/inventories"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeliveryReturn struct
type DeliveryReturn struct {
	Pb inventories.DeliveryReturn
}

// Get func
func (u *DeliveryReturn) Get(ctx context.Context, db *sql.DB) error {
	query := `
		SELECT delivery_returns.id, delivery_returns.company_id, delivery_returns.branch_id, delivery_returns.branch_name, delivery_returns.delivery_id, delivery_returns.code, 
		delivery_returns.return_date, delivery_returns.remark, delivery_returns.created_at, delivery_returns.created_by, delivery_returns.updated_at, delivery_returns.updated_by,
		json_agg(DISTINCT jsonb_build_object(
			'id', delivery_return_details.id,
			'delivery_return_id', delivery_return_details.delivery_return_id,
			'product_id', delivery_return_details.product_id,
			'product_name', products.name,
			'product_code', products.code,
			'shelve_id', delivery_return_details.shelve_id,
			'shelve_code', shelves.code
		)) as details
		FROM delivery_returns 
		JOIN delivery_return_details ON delivery_returns.id = delivery_return_details.delivery_return_id
		JOIN products ON delivery_return_details.product_id = products.id
		JOIN shelves ON delivery_return_details.shelve_id = shelves.id
		WHERE delivery_returns.id = $1
	`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get delivery return: %v", err)
	}
	defer stmt.Close()

	var dateReturn, createdAt, updatedAt time.Time
	var companyID, details string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.BranchId, &u.Pb.BranchName, &u.Pb.Delivery.Id, &u.Pb.Code, &dateReturn, &u.Pb.Remark,
		&createdAt, &u.Pb.CreatedBy, &updatedAt, &u.Pb.UpdatedBy, &details,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code delivery return: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code delivery return: %v", err)
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

	detailDeliveryReturns := []struct {
		ID               string
		DeliveryReturnID string
		ProductID        string
		ProductName      string
		ProductCode      string
		ShelveID         string
		ShelveCode       string
	}{}
	err = json.Unmarshal([]byte(details), &detailDeliveryReturns)
	if err != nil {
		return status.Errorf(codes.Internal, "unmarshal detailReceiveReturns: %v", err)
	}

	for _, detail := range detailDeliveryReturns {
		u.Pb.Details = append(u.Pb.Details, &inventories.DeliveryReturnDetail{
			Id: detail.ID,
			Product: &inventories.Product{
				Id:   detail.ProductID,
				Code: detail.ProductCode,
				Name: detail.ProductName,
			},
			DeliveryReturnId: detail.DeliveryReturnID,
			Shelve: &inventories.Shelve{
				Id:   detail.ShelveID,
				Code: detail.ShelveCode,
			},
		})
	}

	return nil
}

// Create DeliveryReturn
func (u *DeliveryReturn) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	now := time.Now().UTC()
	u.Pb.CreatedBy = ctx.Value(app.Ctx("userID")).(string)
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateReturn, err := ptypes.Timestamp(u.Pb.GetReturnDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert Date: %v", err)
	}

	u.Pb.Code, err = util.GetCode(ctx, tx, "delivery_returns", "DR")
	if err != nil {
		return err
	}

	query := `
		INSERT INTO delivery_returns (id, company_id, branch_id, branch_name, delivery_id, code, return_date, remark, created_at, created_by, updated_at, updated_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert delivery return: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		ctx.Value(app.Ctx("companyID")).(string),
		u.Pb.GetBranchId(),
		u.Pb.GetBranchName(),
		u.Pb.GetDelivery().GetId(),
		u.Pb.GetCode(),
		dateReturn,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetCreatedBy(),
		now,
		u.Pb.GetUpdatedBy(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert delivery return: %v", err)
	}

	u.Pb.CreatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert created by: %v", err)
	}

	u.Pb.UpdatedAt = u.Pb.CreatedAt

	for _, detail := range u.Pb.GetDetails() {
		deliveryReturnDetailModel := DeliveryReturnDetail{}
		deliveryReturnDetailModel.Pb = inventories.DeliveryReturnDetail{
			DeliveryReturnId: u.Pb.GetId(),
			Product:          detail.GetProduct(),
			Shelve:           detail.GetShelve(),
		}
		deliveryReturnDetailModel.PbDeliveryReturn = inventories.DeliveryReturn{
			Id:         u.Pb.Id,
			BranchId:   u.Pb.BranchId,
			BranchName: u.Pb.BranchName,
			Delivery:   u.Pb.Delivery,
			Code:       u.Pb.Code,
			ReturnDate: u.Pb.ReturnDate,
			Remark:     u.Pb.Remark,
			CreatedAt:  u.Pb.CreatedAt,
			CreatedBy:  u.Pb.CreatedBy,
			UpdatedAt:  u.Pb.UpdatedAt,
			UpdatedBy:  u.Pb.UpdatedBy,
		}
		err = deliveryReturnDetailModel.Create(ctx, tx)
		if err != nil {
			return err
		}
	}

	return nil
}

// Update DeliveryReturn
func (u *DeliveryReturn) Update(ctx context.Context, tx *sql.Tx) error {
	now := time.Now().UTC()
	u.Pb.UpdatedBy = ctx.Value(app.Ctx("userID")).(string)
	dateReturn, err := ptypes.Timestamp(u.Pb.GetReturnDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert delivery return date: %v", err)
	}

	query := `
		UPDATE delivery_returns SET
		delivery_id = $1,
		return_date = $2,
		remark = $3, 
		updated_at = $4, 
		updated_by= $5
		WHERE id = $6
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update delivery return: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetDelivery().GetId(),
		dateReturn,
		u.Pb.GetRemark(),
		now,
		u.Pb.GetUpdatedBy(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update delivery return: %v", err)
	}

	u.Pb.UpdatedAt, err = ptypes.TimestampProto(now)
	if err != nil {
		return status.Errorf(codes.Internal, "convert updated by: %v", err)
	}

	return nil
}
