package model

import (
	"context"
	"database/sql"
	"inventory-service/internal/pkg/app"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Stock struct
type Stock struct{}

// Closing Stock
func (u *Stock) Closing(ctx context.Context, tx *sql.Tx) error {
	stmt, err := tx.PrepareContext(ctx, `CALL closing_stocks($1, 0, 0)`)
	if err != nil {
		return status.Errorf(codes.Internal, "Prepare closing stock: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, ctx.Value(app.Ctx("companyID")).(string))
	if err != nil {
		return status.Errorf(codes.Internal, "Exec closing stock: %v", err)
	}

	return nil
}
