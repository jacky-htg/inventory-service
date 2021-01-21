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
			company_id	char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			product_id char(36) NOT NULL,
			barcode CHAR(20) NOT NULL,
			transaction_id char(36) NOT NULL,
			transaction_code CHAR(13) NOT NULL,
			transaction_date DATE NOT NULL,
			type CHAR(2) NOT NULL,
			in_out BOOLEAN NOT NULL,
			qty INTEGER NOT NULL, 
			shelve_id char(36) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			CONSTRAINT fk_inventories_to_products FOREIGN KEY (product_id) REFERENCES products(id)
		);`,
	},
	{
		Version:     8,
		Description: "Add Receiving",
		Script: `
		CREATE TABLE receivings (
			id char(36) NOT NULL PRIMARY KEY,
			company_id	char(36) NOT NULL,
			branch_id char(36) NOT NULL,
			purchase_id char(36) NOT NULL,
			code	CHAR(13) NOT NULL,
			date	DATE NOT NULL,
			remark VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
			created_by BIGINT(20) UNSIGNED NOT NULL,
			updated_by BIGINT(20) UNSIGNED NOT NULL,
			UNIQUE(company_id, code)
		);`,
	},
	{
		Version:     9,
		Description: "Add Receiving Details",
		Script: `
		CREATE TABLE receiving_details (
			id char(36) NOT NULL PRIMARY KEY,
			receiving_id	char(36) NOT NULL,
			product_id char(36) NOT NULL,
			qty INTEGER NOT NULL,
			barcode CHAR(20) NOT NULL,
			shelve_id char(36) NOT NULL,
			expired_date TIMESTAMP,
			UNIQUE(barcode, receiving_id),
			UNIQUE(product_id, receiving_id),
			CONSTRAINT fk_receiving_details_to_receivings FOREIGN KEY (receiving_id) REFERENCES receivings(id) ON DELETE CASCADE ON UPDATE CASCADE,
			CONSTRAINT fk_receiving_details_to_products FOREIGN KEY (product_id) REFERENCES products(id),
			CONSTRAINT fk_receiving_details_to_shelves FOREIGN KEY (shelve_id) REFERENCES shelves(id)
		);`,
	},
}

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(db *sql.DB) error {
	driver := darwin.NewGenericDriver(db, darwin.PostgresDialect{})

	d := darwin.New(driver, migrations, nil)

	return d.Migrate()
}
