package model

import (
	"context"
	"database/sql"

	"inventory-service/internal/pkg/app"
	"inventory-service/pb/inventories"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Category struct
type Category struct {
	Pb inventories.Category
}

// Get func
func (u *Category) Get(ctx context.Context, db *sql.DB) error {
	query := `SELECT id, company_id, name FROM categories WHERE id = $1`
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare statement Get category: %v", err)
	}
	defer stmt.Close()

	var companyID string
	err = stmt.QueryRowContext(ctx, u.Pb.GetId()).Scan(&u.Pb.Id, &companyID, &u.Pb.Name)

	if err == sql.ErrNoRows {
		return status.Errorf(codes.NotFound, "Query Raw get category: %v", err)
	}

	if err != nil {
		return status.Errorf(codes.Internal, "Query Raw get category: %v", err)
	}

	if companyID != ctx.Value(app.Ctx("companyID")).(string) {
		return status.Error(codes.Unauthenticated, "its not your company data")
	}

	return nil
}
