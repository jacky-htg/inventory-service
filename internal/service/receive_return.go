package service

import (
	"database/sql"
	"inventory-service/pb/users"
)

// ReceiveReturn struct
type ReceiveReturn struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}
