package model

import (
	"context"
	"database/sql"
	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeliveryReturnDetail struct
type DeliveryReturnDetail struct {
	Pb               inventories.DeliveryReturnDetail
	PbDeliveryReturn inventories.DeliveryReturn
}

// Get func
func (u *DeliveryReturnDetail) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
		SELECT delivery_return_details.id, delivery_returns.company_id, delivery_return_details.delivery_return_id, delivery_return_details.product_id, 
		delivery_return_details.shelve_id 
		FROM delivery_return_details 
		JOIN delivery_returns ON delivery_return_details.delivery_return_id = delivery_returns.id
		WHERE delivery_return_details.id = $1 AND delivery_return_details.delivery_return_id = $2
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get delivery return detail: %v", err)
	}
	defer stmt.Close()

	var pbProduct inventories.Product
	var pbShelve inventories.Shelve
	var companyID string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId(), u.Pb.GetDeliveryReturnId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.DeliveryReturnId, &pbProduct.Id, &pbShelve.Id,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code delivery return detail: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code delivery return detail: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.Product = &pbProduct
	u.Pb.Shelve = &pbShelve

	return nil
}

// Create DeliveryReturnDetail
func (u *DeliveryReturnDetail) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	query := `
		INSERT INTO delivery_return_details (id, delivery_return_id, product_id, shelve_id) 
		VALUES ($1, $2, $3, $4)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert delivery return detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetDeliveryReturnId(),
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert delivery return detail: %v", err)
	}

	transactionDate, err := time.Parse("2006-01-02T15:04:05.000Z", u.PbDeliveryReturn.GetReturnDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert transactiondate inventory: %v", err)
	}
	inventory := Inventory{
		Barcode:         u.Pb.GetId(),
		BranchID:        u.PbDeliveryReturn.GetBranchId(),
		CompanyID:       ctx.Value(app.Ctx("companyID")).(string),
		IsIn:            true,
		ProductID:       u.Pb.GetProduct().GetId(),
		ShelveID:        u.Pb.GetShelve().GetId(),
		TransactionDate: transactionDate,
		TransactionCode: u.PbDeliveryReturn.GetCode(),
		TransactionID:   u.PbDeliveryReturn.GetId(),
		Type:            "DR",
	}
	err = inventory.Create(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

// Update DeliveryReturnDetail
func (u *DeliveryReturnDetail) Update(ctx context.Context, tx *sql.Tx) error {
	query := `
		UPDATE delivery_return_details SET
		product_id = $1, 
		shelve_id = $2
		WHERE id = $3
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update delivery return detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update delivery return detail: %v", err)
	}

	inventory := Inventory{
		Barcode:       u.Pb.GetId(),
		TransactionID: u.PbDeliveryReturn.GetId(),
	}
	err = inventory.Get(ctx, tx)
	if err != nil {
		return err
	}

	inventory.ProductID = u.Pb.GetProduct().GetId()
	inventory.ShelveID = u.Pb.GetShelve().GetId()
	err = inventory.Update(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

// Delete DeliveryReturnDetail
func (u *DeliveryReturnDetail) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM delivery_return_details WHERE id = $1 AND delivery_return_id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete delivery return detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId(), u.Pb.GetDeliveryReturnId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete delivery return detail: %v", err)
	}

	inventory := Inventory{
		Barcode:       u.Pb.GetId(),
		TransactionID: u.Pb.GetDeliveryReturnId(),
	}
	err = inventory.Get(ctx, tx)
	if err != nil {
		return err
	}

	return inventory.Delete(ctx, tx)
}
