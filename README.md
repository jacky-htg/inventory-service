# inventory-service
Inventory service using golang grpc and postgresql. 
- The service is part of ERP microservices.
- The service will be call in local network.
- Using grpc insecure connection

## Get Started
- git clone git@github.com:jacky-htg/inventory-service.git
- make init
- cp .env.example .env (and edit with your environment)
- make migrate
- make seed
- make server
- You can test the service using `go run client/main.go` and select the test case on file client/main.go

## Features
- [ ] Products
- [X] Product Categories
- [X] Brands
- [ ] Warehouses
- [ ] Shelves
- [ ] Good Receivings
- [ ] Receiving Returns
- [ ] Delivery Orders
- [ ] Delivery Returns
- [ ] Internal Warehouse Mutations
- [ ] External Warehouse Mutations
- [ ] Stock Opname
- [ ] Stock Information
- [ ] Product Track History
- [ ] Closing Stocks

## How To Contribute
- Give star or clone and fork the repository
- Report the bug
- Submit issue for request of enhancement
- Pull Request for fixing bug or enhancement module 

## License
[The license of application is GPL-3.0](https://github.com/jacky-htg/user-service/blob/main/LICENSE), You can use this apllication for commercial use, distribution or modification. But there is no liability and warranty. Please read the license details carefully.

## Link Repository
- [API Gateway for ERP](https://github.com/jacky-htg/api-gateway-service)
- [User Service](https://github.com/jacky-htg/user-service)
- [Sales Service](https://github.com/jacky-htg/sales-service)
- [Purchase Service](https://github.com/jacky-htg/purchase-service)
- [General Ledger Service](https://github.com/jacky-htg/ledger-service)
- [Simple gRPC Skeleton](https://github.com/jacky-htg/grpc-skeleton)
