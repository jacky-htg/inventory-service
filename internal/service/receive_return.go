package service

import (
	"context"
	"database/sql"
	"inventory-service/internal/model"
	"inventory-service/pb/inventories"
	"inventory-service/pb/users"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ReceiveReturn struct
type ReceiveReturn struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
	inventories.UnimplementedReceiveReturnServiceServer
}

// Create ReceiveReturn
func (u *ReceiveReturn) Create(ctx context.Context, in *inventories.ReceiveReturn) (*inventories.ReceiveReturn, error) {
	var receiveReturnModel model.ReceiveReturn
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetBranchId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid branch")
		}

		if len(in.GetReceive().GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid receiving")
		}

		if _, err := time.Parse("2006-01-02T15:04:05.000Z", in.GetReturnDate()); err != nil {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid date")
		}
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			return &receiveReturnModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			return &receiveReturnModel.Pb, err
		}
	}

	err = isYourBranch(ctx, u.UserClient, u.RegionClient, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	branch, err := getBranch(ctx, u.BranchClient, in.GetBranchId())
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	receiveReturnModel.Pb = inventories.ReceiveReturn{
		BranchId:   in.GetBranchId(),
		BranchName: branch.GetName(),
		Code:       in.GetCode(),
		ReturnDate: in.GetReturnDate(),
		Receive:    in.GetReceive(),
		Remark:     in.GetRemark(),
		Details:    in.GetDetails(),
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	err = receiveReturnModel.Create(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &receiveReturnModel.Pb, err
	}

	tx.Commit()

	return &receiveReturnModel.Pb, nil
}

// View ReceiveReturn
func (u *ReceiveReturn) View(ctx context.Context, in *inventories.Id) (*inventories.ReceiveReturn, error) {
	var receiveReturnModel model.ReceiveReturn
	var err error

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		receiveReturnModel.Pb.Id = in.GetId()
	}

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	err = receiveReturnModel.Get(ctx, u.Db)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	return &receiveReturnModel.Pb, nil
}

// Update ReceiveReturn
func (u *ReceiveReturn) Update(ctx context.Context, in *inventories.ReceiveReturn) (*inventories.ReceiveReturn, error) {
	var receiveReturnModel model.ReceiveReturn
	var err error

	// TODO : if this month any closing stock, create transaction for thus month will be blocked

	// basic validation
	{
		if len(in.GetId()) == 0 {
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid id")
		}
		receiveReturnModel.Pb.Id = in.GetId()
	}

	// TODO : if any mutation_unit update will be blocked

	ctx, err = getMetadata(ctx)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	err = receiveReturnModel.Get(ctx, u.Db)
	if err != nil {
		return &receiveReturnModel.Pb, err
	}

	if len(in.GetReceive().GetId()) > 0 {
		receiveReturnModel.Pb.Receive = in.GetReceive()
	}

	if _, err := time.Parse("2006-01-02T15:04:05.000Z", in.GetReturnDate()); err == nil {
		receiveReturnModel.Pb.ReturnDate = in.GetReturnDate()
	}

	tx, err := u.Db.BeginTx(ctx, nil)
	if err != nil {
		return &receiveReturnModel.Pb, status.Errorf(codes.Internal, "begin transaction: %v", err)
	}

	err = receiveReturnModel.Update(ctx, tx)
	if err != nil {
		tx.Rollback()
		return &receiveReturnModel.Pb, err
	}

	var newDetails []*inventories.ReceiveReturnDetail
	for _, detail := range in.GetDetails() {
		// product validation
		if len(detail.GetProduct().GetId()) == 0 {
			tx.Rollback()
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid product")
		}

		productModel := model.Product{}
		productModel.Pb = inventories.Product{Id: detail.GetProduct().GetId()}
		err = productModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &receiveReturnModel.Pb, err
		}

		// shelve validation
		if len(detail.GetShelve().GetId()) == 0 {
			tx.Rollback()
			return &receiveReturnModel.Pb, status.Error(codes.InvalidArgument, "Please supply valid shelve")
		}

		shelveModel := model.Shelve{}
		shelveModel.Pb = inventories.Shelve{Id: detail.GetShelve().GetId()}
		err = shelveModel.Get(ctx, u.Db)
		if err != nil {
			tx.Rollback()
			return &receiveReturnModel.Pb, err
		}

		if len(detail.GetId()) > 0 {
			// operasi update
			receiveReturnDetailModel := model.ReceiveReturnDetail{}
			receiveReturnDetailModel.Pb.Id = detail.GetId()
			receiveReturnDetailModel.Pb.ReceiveReturnId = receiveReturnModel.Pb.GetId()
			err = receiveReturnDetailModel.Get(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &receiveReturnModel.Pb, err
			}

			receiveReturnDetailModel.Pb.Product = detail.GetProduct()
			receiveReturnDetailModel.Pb.Shelve = detail.GetShelve()
			receiveReturnDetailModel.PbReceiveReturn = inventories.ReceiveReturn{
				Id:         receiveReturnModel.Pb.Id,
				BranchId:   receiveReturnModel.Pb.BranchId,
				BranchName: receiveReturnModel.Pb.BranchName,
				Receive:    receiveReturnModel.Pb.GetReceive(),
				Code:       receiveReturnModel.Pb.Code,
				ReturnDate: receiveReturnModel.Pb.ReturnDate,
				Remark:     receiveReturnModel.Pb.Remark,
				CreatedAt:  receiveReturnModel.Pb.CreatedAt,
				CreatedBy:  receiveReturnModel.Pb.CreatedBy,
				UpdatedAt:  receiveReturnModel.Pb.UpdatedAt,
				UpdatedBy:  receiveReturnModel.Pb.UpdatedBy,
				Details:    receiveReturnModel.Pb.Details,
			}
			err = receiveReturnDetailModel.Update(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &receiveReturnModel.Pb, err
			}

			newDetails = append(newDetails, &receiveReturnDetailModel.Pb)
			for index, data := range receiveReturnModel.Pb.GetDetails() {
				if data.GetId() == detail.GetId() {
					receiveReturnModel.Pb.Details = append(receiveReturnModel.Pb.Details[:index], receiveReturnModel.Pb.Details[index+1:]...)
					break
				}
			}

		} else {
			// operasi insert
			receiveReturnDetailModel := model.ReceiveReturnDetail{Pb: inventories.ReceiveReturnDetail{
				ReceiveReturnId: receiveReturnModel.Pb.GetId(),
				Product:         detail.GetProduct(),
				Shelve:          detail.GetShelve(),
			}}
			receiveReturnDetailModel.PbReceiveReturn = inventories.ReceiveReturn{
				Id:         receiveReturnModel.Pb.Id,
				BranchId:   receiveReturnModel.Pb.BranchId,
				BranchName: receiveReturnModel.Pb.BranchName,
				Receive:    receiveReturnModel.Pb.GetReceive(),
				Code:       receiveReturnModel.Pb.Code,
				ReturnDate: receiveReturnModel.Pb.ReturnDate,
				Remark:     receiveReturnModel.Pb.Remark,
				CreatedAt:  receiveReturnModel.Pb.CreatedAt,
				CreatedBy:  receiveReturnModel.Pb.CreatedBy,
				UpdatedAt:  receiveReturnModel.Pb.UpdatedAt,
				UpdatedBy:  receiveReturnModel.Pb.UpdatedBy,
				Details:    receiveReturnModel.Pb.Details,
			}
			err = receiveReturnDetailModel.Create(ctx, tx)
			if err != nil {
				tx.Rollback()
				return &receiveReturnModel.Pb, err
			}

			newDetails = append(newDetails, &receiveReturnDetailModel.Pb)
		}
	}

	// delete existing detail
	for _, data := range receiveReturnModel.Pb.GetDetails() {
		receiveReturnDetailModel := model.ReceiveReturnDetail{Pb: inventories.ReceiveReturnDetail{
			ReceiveReturnId: receiveReturnModel.Pb.GetId(),
			Id:              data.GetId(),
		}}
		err = receiveReturnDetailModel.Delete(ctx, tx)
		if err != nil {
			tx.Rollback()
			return &receiveReturnModel.Pb, err
		}
	}

	tx.Commit()

	return &receiveReturnModel.Pb, nil
}

// List ReceiveReturn
func (u *ReceiveReturn) List(in *inventories.ListReceiveReturnRequest, stream inventories.ReceiveReturnService_ListServer) error {
	ctx := stream.Context()
	ctx, err := getMetadata(ctx)
	if err != nil {
		return err
	}

	var receiveReturnModel model.ReceiveReturn
	query, paramQueries, paginationResponse, err := receiveReturnModel.ListQuery(ctx, u.Db, in)

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

		var pbReceiveReturn inventories.ReceiveReturn
		var companyID string
		var createdAt, updatedAt time.Time
		err = rows.Scan(&pbReceiveReturn.Id, &companyID, &pbReceiveReturn.BranchId, &pbReceiveReturn.BranchName,
			&pbReceiveReturn.Code, &pbReceiveReturn.ReturnDate, &pbReceiveReturn.Remark,
			&createdAt, &pbReceiveReturn.CreatedBy, &updatedAt, &pbReceiveReturn.UpdatedBy)
		if err != nil {
			return status.Errorf(codes.Internal, "scan data: %v", err)
		}

		pbReceiveReturn.CreatedAt = createdAt.String()
		pbReceiveReturn.UpdatedAt = updatedAt.String()

		res := &inventories.ListReceiveReturnResponse{
			Pagination:    paginationResponse,
			ReceiveReturn: &pbReceiveReturn,
		}

		err = stream.Send(res)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
