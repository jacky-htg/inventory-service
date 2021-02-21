package service

import (
	"database/sql"
	"inventory-service/pb/users"
)

// Stock struct
type Stock struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}
