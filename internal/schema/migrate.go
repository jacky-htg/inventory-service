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
			created TIMESTAMP NOT NULL DEFAULT NOW(),
			updated TIMESTAMP NOT NULL DEFAULT NOW(),
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
			created TIMESTAMP NOT NULL DEFAULT NOW(),
			updated TIMESTAMP NOT NULL DEFAULT NOW(),
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
			created TIMESTAMP NOT NULL DEFAULT NOW(),
			updated TIMESTAMP NOT NULL DEFAULT NOW(),
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
			created TIMESTAMP NOT NULL DEFAULT NOW(),
			updated TIMESTAMP NOT NULL DEFAULT NOW(),
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
			created TIMESTAMP NOT NULL DEFAULT NOW(),
			updated TIMESTAMP NOT NULL DEFAULT NOW(),
			UNIQUE(warehouse_id, code),
			CONSTRAINT fk_shelves_to_warehouses FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
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
