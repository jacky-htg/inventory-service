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

	productServer := service.Product{Db: db}
	inventories.RegisterProductServiceServer(grpcServer, &productServer)

	warehouseServer := service.Warehouse{
		Db:           db,
		UserClient:   users.NewUserServiceClient(userConn),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterWarehouseServiceServer(grpcServer, &warehouseServer)

	receiveServer := service.Receive{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterReceiveServiceServer(grpcServer, &receiveServer)

	deliveryServer := service.Delivery{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterDeliveryServiceServer(grpcServer, &deliveryServer)

	receiveReturnServer := service.ReceiveReturn{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterReceiveReturnServiceServer(grpcServer, &receiveReturnServer)

	deliveryReturnServer := service.DeliveryReturn{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterDeliveryReturnServiceServer(grpcServer, &deliveryReturnServer)

	stockServer := service.Stock{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
	}
	inventories.RegisterStockServiceServer(grpcServer, &stockServer)
}
