package model

import (
	"context"
	"database/sql"

	"github.com/jacky-htg/erp-proto/go/pb/inventories"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Category struct
type Category struct {
	Pb inventories.Category
}

// Get func
func (u *Category) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, name FROM categories WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get category: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &u.Pb.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get category: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get category: %v", err)
	}

	return nil
}
