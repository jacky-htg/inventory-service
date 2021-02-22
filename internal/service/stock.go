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

// List Stock
func (u *Stock) List(ctx context.Context, in *inventories.StockListInput) (*inventories.StockList, error) {
	var stockModel model.Stock
	var err error

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &inventories.StockList{}, err
	}

	if len(in.GetBranchId()) > 0 {
		err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
		if err != nil {
			return &inventories.StockList{}, err
		}
	}

	stockModel.ListInput = inventories.StockListInput{
		BranchId: in.BranchId,
	}

	err = stockModel.List(ctx, u.Db)
	if err != nil {
		return &inventories.StockList{}, err
	}

	return &stockModel.StockList, nil
}

// Info Stock
func (u *Stock) Info(ctx context.Context, in *inventories.StockInfoInput) (*inventories.StockInfo, error) {
	var stockModel model.Stock
	var err error

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &inventories.StockInfo{}, err
	}

	if len(in.GetBranchId()) > 0 {
		err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
		if err != nil {
			return &inventories.StockInfo{}, err
		}
	}

	stockModel.InfoInput = inventories.StockInfoInput{
		BranchId:  in.BranchId,
		ProductId: in.ProductId,
	}

	err = stockModel.Info(ctx, u.Db)
	if err != nil {
		return &inventories.StockInfo{}, err
	}

	return &stockModel.StockInfo, nil
}
