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

// ReceiveReturnDetail struct
type ReceiveReturnDetail struct {
	Pb              inventories.ReceiveReturnDetail
	PbReceiveReturn inventories.ReceiveReturn
}

// Create ReceiveReturnDetail
func (u *ReceiveReturnDetail) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	query := `
		INSERT INTO receive_return_details (id, receive_id, product_id, shelve_id) 
		VALUES ($1, $2, $3, $4)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare insert receive detail: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		u.Pb.GetId(),
		u.Pb.GetReceiveReturnId(),
		u.Pb.GetProduct().GetId(),
		u.Pb.GetShelve().GetId(),
	)
	if err != nil {
		return status.Errorf(codes.Internal, "Exec insert receive detail: %v", err)
	}

	transactionDate, err := ptypes.Timestamp(u.PbReceiveReturn.GetReturnDate())
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
