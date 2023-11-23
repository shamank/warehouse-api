package postgres

import (
	"database/sql"
	"github.com/shamank/warehouse-service/internal/domain/models"
	"github.com/shamank/warehouse-service/internal/domain/schemas"
	"github.com/shamank/warehouse-service/internal/repository"
	"github.com/shamank/warehouse-service/internal/service"
	"log/slog"
)

var _ service.Repository = (*PostgresRepo)(nil)

type PostgresRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewPostgresRepo(db *sql.DB, logger *slog.Logger) *PostgresRepo {
	return &PostgresRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresRepo) GetProductsQuantity(productArticle string) ([]models.WarehouseProduct, error) {

	query := `select wp.warehouse_uuid, wp.quantity, wp.reserved_quantity from warehouse_products wp
    			inner join products p on wp.product_uuid = p.uuid
                inner join warehouses w on wp.warehouse_uuid = w.uuid
        			where p.article = $1 and w.is_available = true`

	rows, err := r.db.Query(query, productArticle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productByWarehouses := make([]models.WarehouseProduct, 0)

	for rows.Next() {
		var warehouseProduct models.WarehouseProduct
		err := rows.Scan(&warehouseProduct.WarehouseUUID, &warehouseProduct.Quantity, &warehouseProduct.ReservedQuantity)
		if err != nil {
			r.logger.Error("error scanning warehouse products", err)
			return nil, err
		}
		productByWarehouses = append(productByWarehouses, warehouseProduct)
	}
	return productByWarehouses, nil
}

func (r *PostgresRepo) GetRemainingProductsByWarehouse(warehouseUUID string) ([]models.Product, error) {

	products := make([]models.Product, 0)

	query := `SELECT p.name, p.size, p.article, sum(wp.quantity)
					FROM warehouse_products wp
					    INNER JOIN products p on p.uuid = wp.product_uuid
                           WHERE wp.warehouse_uuid = $1
                            GROUP BY p.name, p.size, p.article`

	rows, err := r.db.Query(query, warehouseUUID)
	if err != nil {
		r.logger.Error("error scanning warehouse products", err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var product models.Product

		err := rows.Scan(&product.Name, &product.Size, &product.Code, &product.Quantity)
		if err != nil {
			r.logger.Error("error scanning warehouse products", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil

}

func (r *PostgresRepo) ReserveProducts(products []schemas.ProductWarehouseSplitted) error {
	tx, err := r.db.Begin()
	if err != nil {
		// TODO: обработка ошибки
		return err
	}

	for _, product := range products {
		for _, warehouseData := range product.WarehouseData {
			err := r.updateProductQuantities(tx, product.ProductArticle, warehouseData.WarehouseUUID, -warehouseData.Count, warehouseData.Count)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()

}

func (r *PostgresRepo) ReleaseProducts(productsWithSplit []schemas.ProductWarehouseSplitted) error {
	tx, err := r.db.Begin()
	if err != nil {
		// TODO: обработка ошибки
		return err
	}

	for _, product := range productsWithSplit {
		for _, warehouseData := range product.WarehouseData {
			err := r.updateProductQuantities(tx, product.ProductArticle, warehouseData.WarehouseUUID, 0, -warehouseData.Count)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (r *PostgresRepo) updateProductQuantities(tx *sql.Tx, productArticle string, warehouseUUID string, quantityDelta int, reservedQuantityDelta int) error {
	query := `UPDATE warehouse_products wp
				SET quantity = wp.quantity + $1, reserved_quantity = wp.reserved_quantity + $2
				FROM products
				WHERE wp.product_uuid = products.uuid AND products.article = $3 AND wp.warehouse_uuid = $4`

	result, err := tx.Exec(query, quantityDelta, reservedQuantityDelta, productArticle, warehouseUUID)
	if err != nil {
		r.logger.Error("error occurred while updating products", err)
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("error rows affected: ", err)
		return err
	}

	if rows == 0 {
		return repository.ErrNoUpdatedProducts
	}

	return nil
}

func (r *PostgresRepo) GenerateTestData() error {

	query1 := `INSERT INTO warehouses (uuid, name, is_available) VALUES
					('af5fc7cd-afb0-43f8-a9d2-ce532512b2ac', 'warehouse1', true),
					('f1dd9277-a8af-49ee-a5b1-3f8ee9e74cfd', 'warehouse2', false),
					('c1bf338d-1953-4b9f-8dd7-71dfca0a29cc', 'warehouse3', true);`

	query2 := `INSERT INTO products (uuid, name, size, article) VALUES
					('854427c7-c53c-40be-935f-a97df1c89a13', 'product1', 10, '123'),
					('a8bff1ba-125a-45cb-b779-a1f5b813f0c3', 'product2', 20, '456'),
					('c175b84e-a62c-4094-871f-4c03c34aa37e', 'product3', 30, '789'),
					('5cb17c38-aa38-4797-a295-475244bb2e53', 'product4', 40, '987'),
					('d19031d1-eb57-4e2b-9c0b-db80fd694a51', 'product5', 50, '654'),
					('463d8a77-7916-4c1f-94b3-2408017f22da', 'product6', 60, '321');`

	query3 := `INSERT INTO warehouse_products (warehouse_uuid, product_uuid, quantity, reserved_quantity) VALUES
					('af5fc7cd-afb0-43f8-a9d2-ce532512b2ac', '854427c7-c53c-40be-935f-a97df1c89a13', 15, 0),
					('af5fc7cd-afb0-43f8-a9d2-ce532512b2ac', 'a8bff1ba-125a-45cb-b779-a1f5b813f0c3', 25, 2),
					('af5fc7cd-afb0-43f8-a9d2-ce532512b2ac', 'c175b84e-a62c-4094-871f-4c03c34aa37e', 35, 5),
					('af5fc7cd-afb0-43f8-a9d2-ce532512b2ac', '5cb17c38-aa38-4797-a295-475244bb2e53', 45, 7),
					('f1dd9277-a8af-49ee-a5b1-3f8ee9e74cfd', '463d8a77-7916-4c1f-94b3-2408017f22da', 12, 0),
					('f1dd9277-a8af-49ee-a5b1-3f8ee9e74cfd', '5cb17c38-aa38-4797-a295-475244bb2e53', 32, 23),
					('c1bf338d-1953-4b9f-8dd7-71dfca0a29cc', '463d8a77-7916-4c1f-94b3-2408017f22da', 10, 30),
					('c1bf338d-1953-4b9f-8dd7-71dfca0a29cc', '5cb17c38-aa38-4797-a295-475244bb2e53', 0, 20),
					('c1bf338d-1953-4b9f-8dd7-71dfca0a29cc', 'd19031d1-eb57-4e2b-9c0b-db80fd694a51', 2, 5);`

	tx, err := r.db.Begin()
	if err != nil {
		r.logger.Error("failed to start transaction", err)
		return err
	}

	for _, query := range []string{query1, query2, query3} {
		if _, err := tx.Exec(query); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
