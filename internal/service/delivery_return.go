package service

import (
	"database/sql"
	"inventory-service/pb/users"
)

// DeliveryReturn struct
type DeliveryReturn struct {
	Db           *sql.DB
	UserClient   users.UserServiceClient
	RegionClient users.RegionServiceClient
	BranchClient users.BranchServiceClient
}
