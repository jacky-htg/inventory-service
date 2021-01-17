package service

import (
	"database/sql"
	"inventory-service/pb/inventories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Category struct
type Category struct {
	Db *sql.DB
}

// List Group
func (u *Group) List(in *inventories.Empty, stream inventories.CategoryService_ListServer) error {
	ctx := stream.Context()
	rows, err := u.Db.QueryContext(ctx, `SELECT id, name FROM categories`)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := contextError(ctx)
		if err != nil {
			return err
		}

		var pbCategory inventories.Category
		err = rows.Scan(&pbCategory.Id, &pbCategory.Name)
		if err != nil {
			return err
		}

		err = stream.Send(&pbCategory)
		if err != nil {
			return status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
		}
	}
	return nil
}
