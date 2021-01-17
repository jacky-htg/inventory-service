package route

import (
	"database/sql"
	"inventory-service/internal/service"
	"inventory-service/pb/inventories"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// GrpcRoute func
func GrpcRoute(grpcServer *grpc.Server, db *sql.DB, log *logrus.Entry) {
	categoryServer := service.Category{Db: db}
	inventories.RegisterCategoryServiceServer(grpcServer, &categoryServer)

}
