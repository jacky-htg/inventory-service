package service

import (
	"database/sql"
	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Category struct
type Category struct {
	Db *sql.DB
	inventories.UnimplementedCategoryServiceServer
}

// List Category
func (u *Category) List(in *inventories.MyEmpty, stream inventories.CategoryService_ListServer) error {
	ctx := stream.Context()
	rows, err := u.Db.QueryContext(ctx, `SELECT id, name FROM categories`)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := app.ContextError(ctx)
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
