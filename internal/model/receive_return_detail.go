package model

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jacky-htg/erp-pkg/app"
	"github.com/jacky-htg/erp-proto/go/pb/inventories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveReturnDetail struct
type ReceiveReturnDetail struct {
	Pb              inventories.ReceiveReturnDetail
	PbReceiveReturn inventories.ReceiveReturn
}

// Get func
func (u *ReceiveReturnDetail) Get(ctx context.Context, tx *sql.Tx) error {
	query := `
		SELECT receive_return_details.id, receive_returns.company_id, receive_return_details.receive_return_id, receive_return_details.product_id, 
		receive_return_details.shelve_id 
		FROM receive_return_details 
		JOIN receive_returns ON receive_return_details.receive_return_id = receive_returns.id
		WHERE receive_return_details.id = $1 AND receive_return_details.receive_return_id = $2
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get receive return detail: %v", err)
	}
	defer stmt.Close()

	var pbProduct inventories.Product
	var pbShelve inventories.Shelve
	var companyID string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId(), u.Pb.GetReceiveReturnId()).Scan(
		&u.Pb.Id, &companyID, &u.Pb.ReceiveReturnId, &pbProduct.Id, &pbShelve.Id,
	)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get by code receive return detail: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get by code receive return detail: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company")
	}

	u.Pb.Product = &pbProduct
	u.Pb.Shelve = &pbShelve

	return nil
}

// Create ReceiveReturnDetail
func (u *ReceiveReturnDetail) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	query := `
		INSERT INTO receive_return_details (id, receive_return_id, product_id, shelve_id) 
		VALUES ($1, $2, $3, $4)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert receive return detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetReceiveReturnId(),
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert receive return detail: %v", err)
	}

	transactionDate, err := time.Parse("2006-01-02T15:04:05.000Z", u.PbReceiveReturn.GetReturnDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert transactiondate inventory: %v", err)
	}
	inventory := Inventory{
		Barcode:         u.Pb.GetId(),
		BranchID:        u.PbReceiveReturn.GetBranchId(),
		CompanyID:       ctx.Value(app.Ctx("companyID")).(string),
		IsIn:            true,
		ProductID:       u.Pb.GetProduct().GetId(),
		ShelveID:        u.Pb.GetShelve().GetId(),
		TransactionDate: transactionDate,
		TransactionCode: u.PbReceiveReturn.GetCode(),
		TransactionID:   u.PbReceiveReturn.GetId(),
		Type:            "RR",
	}
	err = inventory.Create(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

// Update ReceiveReturnDetail
func (u *ReceiveReturnDetail) Update(ctx context.Context, tx *sql.Tx) error {
	query := `
		UPDATE receive_return_details SET
		product_id = $1, 
		shelve_id = $2
		WHERE id = $3
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare update receive return detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
		u.Pb.GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec update receive return detail: %v", err)
	}

	inventory := Inventory{
		Barcode:       u.Pb.GetId(),
		TransactionID: u.PbReceiveReturn.GetId(),
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

// Delete ReceiveReturnDetail
func (u *ReceiveReturnDetail) Delete(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `DELETE FROM receive_return_details WHERE id = $1 AND receive_return_id = $2`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare delete receive return detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, u.Pb.GetId(), u.Pb.GetReceiveReturnId())
	if err != nil {
		return status.Errorf(codes.Internal, "Exec delete receive return detail: %v", err)
	}

	inventory := Inventory{
		Barcode:       u.Pb.GetId(),
		TransactionID: u.Pb.GetReceiveReturnId(),
	}
	err = inventory.Get(ctx, tx)
	if err != nil {
		return err
	}

	return inventory.Delete(ctx, tx)
}
