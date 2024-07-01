package route

import (
	"database/sql"
	"inventory-service/internal/service"
	"inventory-service/pb/inventories"
	"inventory-service/pb/purchases"
	"inventory-service/pb/users"
	"log"

	"google.golang.org/grpc"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log map[string]*log.Logger,
	userConn, purchaseConn *grpc.ClientConn) {
	categoryServer := service.Category{Db: db, Log: log}
	inventories.RegisterCategoryServiceServer(grpcServer, &categoryServer)

	productCategoryServer := service.ProductCategory{Db: db, Log: log}
	inventories.RegisterProductCategoryServiceServer(grpcServer, &productCategoryServer)

	brandServer := service.Brand{Db: db, Log: log}
	inventories.RegisterBrandServiceServer(grpcServer, &brandServer)

	productServer := service.Product{Db: db, Log: log}
	inventories.RegisterProductServiceServer(grpcServer, &productServer)

	shelveServer := service.Shelve{Db: db, Log: log}
	inventories.RegisterShelveServiceServer(grpcServer, &shelveServer)

	warehouseServer := service.Warehouse{
		Db:           db,
		UserClient:   users.NewUserServiceClient(userConn),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
		Log:          log,
	}
	inventories.RegisterWarehouseServiceServer(grpcServer, &warehouseServer)

	receiveServer := service.Receive{
		Db:             db,
		UserClient:     users.NewUserServiceClient((userConn)),
		RegionClient:   users.NewRegionServiceClient(userConn),
		BranchClient:   users.NewBranchServiceClient(userConn),
		PurchaseClient: purchases.NewPurchaseServiceClient(purchaseConn),
		Log:            log,
	}
	inventories.RegisterReceiveServiceServer(grpcServer, &receiveServer)

	deliveryServer := service.Delivery{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
		Log:          log,
	}
	inventories.RegisterDeliveryServiceServer(grpcServer, &deliveryServer)

	receiveReturnServer := service.ReceiveReturn{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
		Log:          log,
	}
	inventories.RegisterReceiveReturnServiceServer(grpcServer, &receiveReturnServer)

	deliveryReturnServer := service.DeliveryReturn{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
		Log:          log,
	}
	inventories.RegisterDeliveryReturnServiceServer(grpcServer, &deliveryReturnServer)

	stockServer := service.Stock{
		Db:           db,
		UserClient:   users.NewUserServiceClient((userConn)),
		RegionClient: users.NewRegionServiceClient(userConn),
		BranchClient: users.NewBranchServiceClient(userConn),
		Log:          log,
	}
	inventories.RegisterStockServiceServer(grpcServer, &stockServer)
}
