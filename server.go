package main

import (
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"inventory-service/internal/config"
	"inventory-service/internal/middleware"
	"inventory-service/internal/pkg/db/postgres"
	"inventory-service/internal/route"
)

const defaultPort = "8001"

func main() {
	// lookup and setup env
	if _, ok := os.LookupEnv("PORT"); !ok {
		config.Setup(".env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// init log
	log := map[string]*log.Logger{
		"error":   log.New(os.Stdout, "ERROR: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		"info":    log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		"warning": log.New(os.Stdout, "WARNING: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}

	// create postgres database connection
	db, err := postgres.Open()
	if err != nil {
		log["error"].Fatalf("connecting to db: %v", err)
		return
	}
	log["info"].Println("connecting to postgresql database")

	defer db.Close()

	// listen tcp port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log["error"].Fatalf("failed to listen: %v", err)
	}

	mdInterceptor := middleware.Metadata{}
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			mdInterceptor.Unary(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			mdInterceptor.Stream(),
		)),
	}

	grpcServer := grpc.NewServer(serverOptions...)

	userConn, err := grpc.NewClient(os.Getenv("USER_SERVICE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log["info"].Printf("create user service connection: %v", err)
	}
	defer userConn.Close()

	purchaseConn, err := grpc.NewClient(os.Getenv("PURCHASE_SERVICE"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log["info"].Printf("create purchase service connection: %v", err)
	}
	defer purchaseConn.Close()

	// routing grpc services
	route.GrpcRoute(grpcServer, db, log, userConn, purchaseConn)

	if err := grpcServer.Serve(lis); err != nil {
		log["error"].Fatalf("failed to serve: %s", err)
		return
	}
	log["info"].Println("serve grpc on port: " + port)
}
