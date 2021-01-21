package route

import (
	"database/sql"
	"inventory-service/internal/service"
	"inventory-service/pb/inventories"
	"inventory-service/pb/users"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log *logrus.Entry, userConn *grpc.ClientConn) {
	categoryServer := service.Category{Db: db}
	inventories.RegisterCategoryServiceServer(grpcServer, &categoryServer)

	productCategoryServer := service.ProductCategory{Db: db}
	inventories.RegisterProductCategoryServiceServer(grpcServer, &productCategoryServer)

	brandServer := service.Brand{Db: db}
	inventories.RegisterBrandServiceServer(grpcServer, &brandServer)

	warehouseServer := service.Warehouse{
		Db:           db,
		UserClient:   users.NewUserServiceClient(userConn),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterWarehouseServiceServer(grpcServer, &warehouseServer)

}
