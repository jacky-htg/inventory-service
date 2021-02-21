package schema

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var migrations = []darwin.Migration{
	{
		Version:     1,
		Description: "Add Categories",
		Script: `
		CREATE TABLE categories (
			id char(36) NOT NULL PRIMARY KEY,
			name VARCHAR(45) NOT NULL UNIQUE
		);`,
	},
	{
		Version:     2,
		Description: "Add Product Categories",
		Script: `
		CREATE TABLE product_categories (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			category_id char(36) NOT NULL,
			name VARCHAR(45) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_by char(36) NOT NULL,
			CONSTRAINT fk_product_categories_to_categories FOREIGN KEY (category_id) REFERENCES categories(id)
		);`,
	},
	{
		Version:     3,
		Description: "Add Brands",
		Script: `
		CREATE TABLE brands (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			code CHAR(10) NOT NULL,
			name VARCHAR(45) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     4,
		Description: "Add Products",
		Script: `
		CREATE TABLE products (
			id char(36) NOT NULL PRIMARY KEY,
			company_id char(36) NOT NULL,
			brand_id char(36) NOT NULL,
			product_category_id char(36) NOT NULL,
			code CHAR(10) NOT NULL,
			name VARCHAR(255) NOT NULL,
			minimum_stock INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code),
			CONSTRAINT fk_products_to_brands FOREIGN KEY (brand_id) REFERENCES brands(id),
			CONSTRAINT fk_products_to_product_categories FOREIGN KEY (product_category_id) REFERENCES product_categories(id)
		);`,
	},
	{
		Version:     5,
		Description: "Add Warehouses",
		Script: `
		CREATE TABLE warehouses (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			branch_id char(36) NULL,
			branch_name varchar(100) NULL,
			code	CHAR(10) NOT NULL,
			name	VARCHAR(255) NOT NULL,
			pic_name	VARCHAR(255) NOT NULL,
			pic_phone VARCHAR(20) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     6,
		Description: "Add Shelves",
		Script: `
		CREATE TABLE shelves (
			id char(36) NOT NULL PRIMARY KEY,
			warehouse_id	char(36) NOT NULL,
			code CHAR(10) NOT NULL,
			capacity INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_by char(36) NOT NULL,
			UNIQUE(warehouse_id, code),
			CONSTRAINT fk_shelves_to_warehouses FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
		);`,
	},
	{
		Version:     7,
		Description: "Add Inventories",
		Script: `
		CREATE TABLE inventories (
			id char(36) NOT NULL PRIMARY KEY,
			company_id char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			product_id char(36) NOT NULL,
			barcode char(36) NOT NULL,
			transaction_id char(36) NOT NULL,
			transaction_code CHAR(13) NOT NULL,
			transaction_date DATE NOT NULL,
			type CHAR(2) NOT NULL,
			in_out BOOLEAN NOT NULL,
			shelve_id char(36) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_inventories_to_products FOREIGN KEY (product_id) REFERENCES products(id)
		);`,
	},
	{
		Version:     8,
		Description: "Add Receives",
		Script: `
		CREATE TABLE receives (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			branch_name varchar(100) NOT NULL,
			purchase_id char(36) NOT NULL,
			code	CHAR(13) NOT NULL,
			receive_date	DATE NOT NULL,
			remark VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     9,
		Description: "Add Receive Details",
		Script: `
		CREATE TABLE receive_details (
			id char(36) NOT NULL PRIMARY KEY,
			receive_id	char(36) NOT NULL,
			product_id char(36) NOT NULL,
			shelve_id char(36) NOT NULL,
			expired_date TIMESTAMP,
			CONSTRAINT fk_receive_details_to_receives FOREIGN KEY (receive_id) REFERENCES receives(id) ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT fk_receive_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
			CONSTRAINT fk_receive_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
		);`,
	},
	{
		Version:     10,
		Description: "Add Deliveries",
		Script: `
		CREATE TABLE deliveries (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			branch_name varchar(100) NOT NULL,
			sales_order_id char(36) NOT NULL,
			code	CHAR(13) NOT NULL,
			delivery_date	DATE NOT NULL,
			remark VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     11,
		Description: "Add Delivery Details",
		Script: `
		CREATE TABLE delivery_details (
			id char(36) NOT NULL PRIMARY KEY,
			delivery_id	char(36) NOT NULL,
			product_id char(36) NOT NULL,
			shelve_id char(36) NOT NULL,
			CONSTRAINT fk_delivery_details_to_deliveries FOREIGN KEY (delivery_id) REFERENCES deliveries(id) ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT fk_delivery_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
			CONSTRAINT fk_delivery_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
		);`,
	},
	{
		Version:     12,
		Description: "Add Receive Return",
		Script: `
		CREATE TABLE receive_returns (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			branch_name varchar(100) NOT NULL,
			receive_id char(36) NOT NULL,
			code	CHAR(13) NOT NULL,
			return_date	DATE NOT NULL,
			remark VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     13,
		Description: "Add Receive Return Details",
		Script: `
		CREATE TABLE receive_return_details (
			id char(36) NOT NULL PRIMARY KEY,
			receive_return_id	char(36) NOT NULL,
			product_id char(36) NOT NULL,
			shelve_id char(36) NOT NULL,
			CONSTRAINT fk_receive_return_details_to_receive_returns FOREIGN KEY (receive_return_id) REFERENCES receive_returns(id) ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT fk_receive_return_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
			CONSTRAINT fk_receive_return_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
		);`,
	},
	{
		Version:     14,
		Description: "Add Delivery Return",
		Script: `
		CREATE TABLE delivery_returns (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			branch_name varchar(100) NOT NULL,
			delivery_id char(36) NOT NULL,
			code	CHAR(13) NOT NULL,
			return_date	DATE NOT NULL,
			remark VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by char(36) NOT NULL,
			updated_by char(36) NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     15,
		Description: "Add Delivery Return Details",
		Script: `
		CREATE TABLE delivery_return_details (
			id char(36) NOT NULL PRIMARY KEY,
			delivery_return_id	char(36) NOT NULL,
			product_id char(36) NOT NULL,
			shelve_id char(36) NOT NULL,
			CONSTRAINT fk_delivery_return_details_to_delivery_returns FOREIGN KEY (delivery_return_id) REFERENCES delivery_returns(id) ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT fk_delivery_return_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
			CONSTRAINT fk_delivery_return_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
		);`,
	},
	{
		Version:     16,
		Description: "Add saldo stocks",
		Script: `
	CREATE TABLE saldo_stocks (
		id bigserial PRIMARY KEY,
		company_id	char(36) NOT NULL,
		product_id char(36) NOT NULL,
		qty integer NOT NULL,
		year integer NOT NULL,
		month smallint NOT NULL,
		CONSTRAINT fk_saldo_stocks_to_products FOREIGN KEY (product_id) REFERENCES products(id)
	);`,
	},
	{
		Version:     17,
		Description: "Add Saldo Stock Details",
		Script: `
		CREATE TABLE saldo_stock_details (
			id   bigserial PRIMARY KEY,
			saldo_stock_id	bigint NOT NULL,
			branch_id char(36) NOT NULL,
			code CHAR(36) NOT NULL,
			UNIQUE(code, saldo_stock_id),
			CONSTRAINT fk_saldo_stock_details_to_saldo_stocks FOREIGN KEY (saldo_stock_id) REFERENCES saldo_stocks(id) ON DELETE CASCADE ON UPDATE CASCADE
		);`,
	},
	{
		Version:     18,
		Description: "Add Closing Stock Details",
		Script: `
		create or replace procedure closing_stock_details(
			companyID char, 
			curYear int, 
			curMonth int
		) 
		as $$
		declare 
				nextYear int;
				nextMonth int;
		begin	
			
			IF curYear = 0 THEN
				select date_part('year', CURRENT_DATE) into curYear;
			END IF;
		
			IF curMonth = 0 THEN 
				select date_part('month', CURRENT_DATE) into curMonth;
			END IF;
			
			select curYear into nextYear;
			select curMonth + 1 into nextMonth;
			
			IF curMonth = 12 THEN 
				select curYear+1 into nextYear;
				select 1 into nextMonth;
			END IF;
		
			INSERT INTO saldo_stock_details (saldo_stock_id, branch_id, code)
			SELECT 
				saldo_stocks.id saldo_stock_id,
				COALESCE(inventories.branch_id, group_inventories.branch_id) branch_id,
				COALESCE(inventories.barcode, group_inventories.barcode) code
			FROM (
				SELECT 
					MAX(union_inventories.id) id, 
					MAX(union_inventories.company_id) company_id, 
					MAX(union_inventories.branch_id) branch_id, 
					MAX(union_inventories.product_id) product_id, 
					MAX(union_inventories.barcode) barcode, 
					SUM(union_inventories.qty) qty
				FROM (
					(SELECT 
						'0' id,
						saldo_stocks.company_id, 
						saldo_stock_details.branch_id,
						saldo_stocks.product_id, 
						saldo_stock_details.code barcode,
						1 qty  
					FROM saldo_stocks
					JOIN saldo_stock_details ON saldo_stocks.id = saldo_stock_details.saldo_stock_id
					WHERE saldo_stocks.year = curYear AND saldo_stocks.month = curMonth and saldo_stocks.company_id = companyID)
					union all
					(SELECT 
						inventories.id,
						inventories.company_id,
						inventories.branch_id,
						inventories.product_id,
						inventories.barcode,
						case 
							when inventories.in_out then 1
							else -1
						end 
						as qty
					FROM inventories
					where date_part('month', inventories.transaction_date)=curMonth and date_part('year', inventories.transaction_date)=curYear and inventories.company_id = companyID)
				) union_inventories
				GROUP BY union_inventories.company_id, union_inventories.product_id, union_inventories.barcode
			) group_inventories
			left join inventories ON group_inventories.id = inventories.id and inventories.company_id = companyID
			join saldo_stocks ON COALESCE(inventories.company_id, group_inventories.company_id)=saldo_stocks.company_id and COALESCE(inventories.product_id, group_inventories.product_id)=saldo_stocks.product_id and saldo_stocks.year = nextYear and saldo_stocks.month = nextMonth and saldo_stocks.company_id = companyID
			WHERE group_inventories.qty > 0;
			
		end;
		$$ language plpgsql `,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
