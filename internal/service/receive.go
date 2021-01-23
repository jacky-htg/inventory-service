package service

import (
	"context"
	"database/sql"
	"inventory-service/pb/users"
	"time"

	"inventory-service/internal/model"
	"inventory-service/pb/inventories"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Receive struct
type Receive struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}

// Create Receive
func (u *Receive) Create(ctx context.Context, in *inventories.Receive) (*inventories.Receive, error) {
	var receiveModel model.Receive
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetBranchId()) == 0 {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid branch")
		}

		if len(in.GetPurchaseId()) == 0 {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid purchase")
		}

		if in.GetReceiveDate().IsValid() {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid date")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveModel.Pb, err
	}

	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			return &receiveModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			return &receiveModel.Pb, err
		}

		if !detail.GetExpiredDate().IsValid() {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid expired date")
		}
	}

	err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &receiveModel.Pb, err
	}

	branch, err := getBranch(ctx, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &receiveModel.Pb, err
	}

	receiveModel.Pb = inventories.Receive{
		BranchId:    in.GetBranchId(),
		BranchName:  branch.GetName(),
		Code:        in.GetCode(),
		ReceiveDate: in.GetReceiveDate(),
		PurchaseId:  in.GetPurchaseId(),
		Remark:      in.GetRemark(),
		Details:     in.GetDetails(),
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &receiveModel.Pb, err
	}

	err = receiveModel.Create(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &receiveModel.Pb, err
	}

	tx.Commit()

	return &receiveModel.Pb, nil
}

// Update Receive
func (u *Receive) Update(ctx context.Context, in *inventories.Receive) (*inventories.Receive, error) {
	var receiveModel model.Receive
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		receiveModel.Pb.Id = in.GetId()
	}

	// TODO : if any mutation_unit, return receive, or delivery order, update will be blocked

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveModel.Pb, err
	}

	err = receiveModel.Get(ctx, u.Db)
	if err != nil {
		return &receiveModel.Pb, err
	}

	if len(in.GetPurchaseId()) > 0 {
		receiveModel.Pb.PurchaseId = in.GetPurchaseId()
	}

	if in.GetReceiveDate().IsValid() {
		receiveModel.Pb.ReceiveDate = in.GetReceiveDate()
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &receiveModel.Pb, status.Errorf(codes.Internal, "begin transaction: %v", err)
	}

	err = receiveModel.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &receiveModel.Pb, err
	}

	var newDetails []*inventories.ReceiveDetail
	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			tx.Rollback()
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &receiveModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			tx.Rollback()
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &receiveModel.Pb, err
		}

		if !detail.GetExpiredDate().IsValid() {
			tx.Rollback()
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "please supllay valid expired date")
		}

		if len(detail.GetId()) > 0 {
			// operasi update
			receiveDetailModel := model.ReceiveDetail{}
			receiveDetailModel.Pb.Id = detail.GetId()
			receiveDetailModel.Pb.ReceiveId = receiveModel.Pb.GetId()
			err = receiveDetailModel.Get(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &receiveModel.Pb, err
			}

			receiveDetailModel.Pb.ExpiredDate = detail.GetExpiredDate()
			receiveDetailModel.Pb.Product = detail.GetProduct()
			receiveDetailModel.Pb.Shelve = detail.GetShelve()
			receiveDetailModel.PbReceive = receiveModel.Pb
			err = receiveDetailModel.Update(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &receiveModel.Pb, err
			}

			newDetails = append(newDetails, &receiveDetailModel.Pb)
			for index, data := range receiveModel.Pb.GetDetails() {
				if data.GetId() == detail.GetId() {
					receiveModel.Pb.Details = append(receiveModel.Pb.Details[:index], receiveModel.Pb.Details[index+1:]...)
					break
				}
			}

		} else {
			// operasi insert
			pbReceiveDetail := inventories.ReceiveDetail{
				ReceiveId:   receiveModel.Pb.GetId(),
				ExpiredDate: detail.GetExpiredDate(),
				Product:     detail.GetProduct(),
				Shelve:      detail.GetShelve(),
			}
			receiveDetailModel := model.ReceiveDetail{Pb: pbReceiveDetail}
			receiveDetailModel.PbReceive = receiveModel.Pb
			err = receiveDetailModel.Create(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &receiveModel.Pb, err
			}

			newDetails = append(newDetails, &receiveDetailModel.Pb)
		}
	}

	// delete existing detail
	for _, data := range receiveModel.Pb.GetDetails() {
		pbReceiveDetail := inventories.ReceiveDetail{
			ReceiveId: receiveModel.Pb.GetId(),
			Id:        data.GetId(),
		}
		receiveDetailModel := model.ReceiveDetail{Pb: pbReceiveDetail}
		err = receiveDetailModel.Delete(ctx, tx)
		if err != nil {
			tx.Rollback()
			return &receiveModel.Pb, err
		}
	}

	tx.Commit()

	return &receiveModel.Pb, nil
}

// View Receive
func (u *Receive) View(ctx context.Context, in *inventories.Id) (*inventories.Receive, error) {
	var receiveModel model.Receive
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &receiveModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		receiveModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveModel.Pb, err
	}

	err = receiveModel.Get(ctx, u.Db)
	if err != nil {
		return &receiveModel.Pb, err
	}

	return &receiveModel.Pb, nil
}

// List Receive
func (u *Receive) List(in *inventories.ListReceiveRequest, stream inventories.ReceiveService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	var receiveModel model.Receive
	query, paramQueries, paginationResponse, err := receiveModel.ListQuery(ctx, u.Db, in)

	rows, err := u.Db.QueryContext(ctx, query, paramQueries...)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()
	paginationResponse.Pagination = in.GetPagination()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbReceive inventories.Receive
		var companyID string
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbReceive.Id, &companyID, &pbReceive.BranchId, &pbReceive.BranchName,
			&pbReceive.Code, &pbReceive.ReceiveDate, &pbReceive.Remark,
			&createdAt, &pbReceive.CreatedBy, &updatedAt, &pbReceive.UpdatedBy)
		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbReceive.CreatedAt, err = ptypes.TimestampProto(createdAt)
		if err != nil {
			return status.Errorf(codes.Internal, "convert createdAt: %v", err)
		}

		pbReceive.UpdatedAt, err = ptypes.TimestampProto(updatedAt)
		if err != nil {
			return status.Errorf(codes.Internal, "convert updateddAt: %v", err)
		}

		res := &inventories.ListReceiveResponse{
			Pagination: paginationResponse,
			Receive:    &pbReceive,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
