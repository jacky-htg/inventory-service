package model

import (
	"context"
	"database/sql"

	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveDetail struct
type ReceiveDetail struct {
	Pb inventories.ReceiveDetail
}

// Get func
func (u *ReceiveDetail) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
		SELECT receive_details.id, receives.company_id, receive_details.receive_id, receive_details.product_id, 
		receive_details.shelve_id, receive_details.expired_date 
		FROM receive_details 
		JOIN receives ON receive_details.receive_id = receives.id
		WHERE receive_details.id = $1 AND receive_details.receive_id = $2
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get receive detail: %v", err)
	}
	defer stmt.Close()

	var pbProduct inventories.Product
	var pbShelve inventories.Shelve
	var companyID string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId(), u.Pb.GetReceiveId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.ReceiveId, &pbProduct.Id, &pbShelve.Id, &u.Pb.ExpiredDate,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code receive detail: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code receive detail: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.Product = &pbProduct
	u.Pb.Shelve = &pbShelve

	return nil
}

// Create ReceiveDetail
func (u *ReceiveDetail) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	expirdDate, err := ptypes.Timestamp(u.Pb.GetExpiredDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert expired date: %v", err)
	}

	query := `
		INSERT INTO receive_details (id, receive_id, product_id, shelve_id, expired_date) 
		VALUES ($1, $2, $3, $4, $5)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert receive detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetReceiveId(),
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
		expirdDate,
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert receive detail: %v", err)
	}

	return nil
}

// Update ReceiveDetail
func (u *ReceiveDetail) Update(ctx context.Context, tx *sql.Tx) error {
	query := `
		UPDATE receive_details SET
		product_id = $1, 
		shelve_id = $2, 
		expired_date= $3
		WHERE id = $4
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update receive detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
		u.Pb.GetExpiredDate(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update receive detail: %v", err)
	}

	return nil
}

// Delete ReceiveDetail
func (u *ReceiveDetail) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM receive_details WHERE id = $1 AND receive_id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete receive detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId(), u.Pb.GetReceiveId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete receive detail: %v", err)
	}

	return nil
}
