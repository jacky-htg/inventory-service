package model

import (
	"context"
	"database/sql"
	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"

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
