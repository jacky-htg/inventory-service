package service

import (
	"database/sql"
	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Category struct
type Category struct {
	Db  *sql.DB
	Log map[string]*log.Logger
	inventories.UnimplementedCategoryServiceServer
}

// List Category
func (u *Category) List(in *inventories.MyEmpty, stream inventories.CategoryService_ListServer) error {
	ctx := stream.Context()
	rows, err := u.Db.QueryContext(ctx, `SELECT id, name FROM categories`)
	if err != nil {
		err = status.Error(codes.Internal, err.Error())
		u.Log["error"].Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := app.ContextError(ctx)
		if err != nil {
			u.Log["error"].Println(err)
			return err
		}

		var pbCategory inventories.Category
		err = rows.Scan(&pbCategory.Id, &pbCategory.Name)
		if err != nil {
			err = status.Error(codes.Internal, err.Error())
			u.Log["error"].Println(err)
			return err
		}

		err = stream.Send(&pbCategory)
		if err != nil {
			err = status.Errorf(codes.Unknown, "cannot send stream response: %v", err)
			u.Log["error"].Println(err)
			return err
		}
	}
	return nil
}
