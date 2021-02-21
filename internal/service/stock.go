package service

import (
	"context"
	"database/sql"
	"inventory-service/internal/model"
	"inventory-service/pb/inventories"
	"inventory-service/pb/users"
)

// Stock struct
type Stock struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}

// Closing Stock
func (u *Stock) Closing(ctx context.Context, in *inventories.ClosingStockRequest) (*inventories.MyBoolean, error) {
	var stockModel model.Stock
	var err error

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &inventories.MyBoolean{}, err
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &inventories.MyBoolean{}, err
	}

	err = stockModel.Closing(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &inventories.MyBoolean{}, err
	}

	tx.Commit()

	return &inventories.MyBoolean{Boolean: true}, err
}
