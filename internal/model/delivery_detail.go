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

// DeliveryDetail struct
type DeliveryDetail struct {
	Pb         inventories.DeliveryDetail
	PbDelivery inventories.Delivery
}

// Get func
func (u *DeliveryDetail) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
		SELECT delivery_details.id, deliveries.company_id, delivery_details.delivery_id, delivery_details.product_id, 
		delivery_details.shelve_id, delivery_details.barcode 
		FROM delivery_details 
		JOIN deliveries ON delivery_details.delivery_id = deliveries.id
		WHERE delivery_details.id = $1 AND delivery_details.delivery_id = $2
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get delivery detail: %v", err)
	}
	defer stmt.Close()

	var pbProduct inventories.Product
	var pbShelve inventories.Shelve
	var companyID string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId(), u.Pb.GetDeliveryId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.DeliveryId, &pbProduct.Id, &pbShelve.Id, &u.Pb.Barcode,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code delivery detail: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code delivery detail: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.Product = &pbProduct
	u.Pb.Shelve = &pbShelve

	return nil
}

// Create DeliveryDetail
func (u *DeliveryDetail) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	query := `
		INSERT INTO delivery_details (id, delivery_id, product_id, shelve_id, barcode) 
		VALUES ($1, $2, $3, $4, $5)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert delivery detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetDeliveryId(),
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
		u.Pb.GetBarcode(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert delivery detail: %v", err)
	}

	transactionDate, err := ptypes.Timestamp(u.PbDelivery.GetDeliveryDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert transactiondate inventory: %v", err)
	}
	inventory := Inventory{
		Barcode:         u.Pb.GetBarcode(),
		BranchID:        u.PbDelivery.GetBranchId(),
		CompanyID:       ctx.Value(app.Ctx("companyID")).(string),
		IsIn:            false,
		ProductID:       u.Pb.GetProduct().GetId(),
		ShelveID:        u.Pb.GetShelve().GetId(),
		TransactionDate: transactionDate,
		TransactionCode: u.PbDelivery.GetCode(),
		TransactionID:   u.PbDelivery.GetId(),
		Type:            "DO",
	}
	err = inventory.Create(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

// Delete DeliveryDetail
func (u *DeliveryDetail) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM delivery_details WHERE id = $1 AND delivery_id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete delivery detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId(), u.Pb.GetDeliveryId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete delivery detail: %v", err)
	}

	inventory := Inventory{
		Barcode:       u.Pb.GetBarcode(),
		TransactionID: u.Pb.GetDeliveryId(),
	}
	err = inventory.Get(ctx, tx)
	if err != nil {
		return err
	}

	return inventory.Delete(ctx, tx)
}
