package model

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jacky-htg/erp-pkg/app"
	"github.com/jacky-htg/erp-pkg/util"
	"github.com/jacky-htg/erp-proto/go/pb/inventories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveDetail struct
type ReceiveDetail struct {
	Pb        inventories.ReceiveDetail
	PbReceive inventories.Receive
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

// Get func
func (u *ReceiveDetail) ListByPurchaseId(ctx context.Context, tx *sql.Tx, purchaseId string, ids []string) ([]*inventories.OutstandingDetail, error) {
	var output []*inventories.OutstandingDetail
	query := `
	with receive_details as (
		select receive_details.product_id, COUNT(receive_details.product_id) quantity
		from receive_details 
		join receives ON receive_details.receive_id = receives.id
		where receives.purchase_id = $1
		group by receive_details.product_id
	)
	select products.id, products.code, products.name, coalesce(receive_details.quantity, 0)
	from products
	left join receive_details on products.id = receive_details.product_id
	`
	where := []string{}
	paramQueries := []interface{}{purchaseId}

	if len(ids) > 0 {
		productIds := make([]interface{}, len(ids))
		for i, productId := range ids {
			productIds[i] = productId
		}
		var iCond string
		paramQueries, iCond = util.ConvertWhereIn("products.id", paramQueries, productIds)
		where = append(where, iCond)
	}

	if len(where) > 0 {
		query += ` WHERE ` + strings.Join(where, " AND ")
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Prepare statement ListByPurchaseId: %v", err)
	}
	defer stmt.Close()

	row, err := stmt.QueryContext(ctx, paramQueries...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "QueryContext ListByPurchaseId: %v", err)
	}

	for row.Next() {
		var detail inventories.OutstandingDetail
		err = row.Scan(&detail.ProductId, &detail.ProductCode, &detail.ProductName, &detail.Quantity)
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Query Raw ListByPurchaseId: %v", err)
		}

		if err != nil {
			return nil, status.Errorf(codes.Internal, "Query Raw ListByPurchaseId: %v", err)
		}

		output = append(output, &detail)
	}

	if row.Err() != nil {
		return nil, status.Error(codes.Internal, row.Err().Error())
	}

	return output, nil
}

// Create ReceiveDetail
func (u *ReceiveDetail) Create(ctx context.Context, tx *sql.Tx) error {
	u.Pb.Id = uuid.New().String()
	expirdDate, err := time.Parse("2006-01-02T15:04:05.000Z", u.Pb.GetExpiredDate())
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

	transactionDate, err := time.Parse("2006-01-02T15:04:05.000Z", u.PbReceive.GetReceiveDate())
	if err != nil {
		return status.Errorf(codes.Internal, "convert transactiondate inventory: %v", err)
	}
	inventory := Inventory{
		Barcode:         u.Pb.GetId(),
		BranchID:        u.PbReceive.GetBranchId(),
		CompanyID:       ctx.Value(app.Ctx("companyID")).(string),
		IsIn:            true,
		ProductID:       u.Pb.GetProduct().GetId(),
		ShelveID:        u.Pb.GetShelve().GetId(),
		TransactionDate: transactionDate,
		TransactionCode: u.PbReceive.GetCode(),
		TransactionID:   u.PbReceive.GetId(),
		Type:            "GR",
	}
	err = inventory.Create(ctx, tx)
	if err != nil {
		return err
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

	inventory := Inventory{
		Barcode:       u.Pb.GetId(),
		TransactionID: u.PbReceive.GetId(),
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

	inventory := Inventory{
		Barcode:       u.Pb.GetId(),
		TransactionID: u.Pb.GetReceiveId(),
	}
	err = inventory.Get(ctx, tx)
	if err != nil {
		return err
	}

	return inventory.Delete(ctx, tx)
}
