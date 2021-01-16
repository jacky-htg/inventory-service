gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:. 

init:
	go mod init inventory-service

migrate:
	go run cmd/cli.go migrate
	
seed:
	go run cmd/cli.go seed

server:
	go run server.go

.PHONY: gen init migrate seed server