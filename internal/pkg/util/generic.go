package util

import (
	"context"
	"database/sql"
	"fmt"
	"inventory-service/internal/pkg/app"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetCode : func to generate transaction code
func GetCode(ctx context.Context, tx *sql.Tx, tableName string, code string) (string, error) {
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(*) FROM `+tableName+` 
			WHERE company_id = $1 AND to_char(created_at, 'YYYY-mm') = to_char(now(), 'YYYY-mm')`,
		ctx.Value(app.Ctx("companyID")).(string)).Scan(&count)

	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}

	return fmt.Sprintf("%s%d%d%d",
		code,
		time.Now().UTC().Year(),
		int(time.Now().UTC().Month()),
		(count + 1)), nil
}

func ConvertWhereIn(field string, paramQueries []interface{}, data []interface{}) ([]interface{}, string) {
	ph := make([]string, len(data))
	tmpLen := make([]interface{}, len(data))
	for i, d := range data {
		paramQueries = append(paramQueries, d)
		ph[i] = "$%d"
		tmpLen[i] = len(paramQueries)
	}

	if len(tmpLen) > 0 {
		where := fmt.Sprintf("%s IN ", field) + "(" + strings.Join(ph, ", ") + ")"
		return paramQueries, fmt.Sprintf(where, tmpLen...)
	}

	return paramQueries, ""
}
