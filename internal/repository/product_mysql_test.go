package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTopProducts_MySQL(t *testing.T) {
	t.Run("succes - get top 5 products by quantity sold", func(t *testing.T) {
		// arrange
		db, err := sql.Open("testDriver", "")
		require.NoError(t, err)
		defer db.Close()

		r := repository.NewProductsMySQL(db)

		// - set up the database
		err = func() (err error) {

			_, err = db.Exec("DELETE FROM customers")
			require.NoError(t, err)

			_, err = db.Exec("DELETE FROM invoices")
			require.NoError(t, err)

			_, err = db.Exec("DELETE FROM products")
			require.NoError(t, err)

			_, err = db.Exec("DELETE FROM sales")
			require.NoError(t, err)

			query := "INSERT INTO customers (id, first_name, last_name, `condition`) VALUES	(1, 'John', 'Doe', 1)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = "INSERT INTO products (id, `description`, price) VALUES (1, 'Product 1', 100.00), (2, 'Product 2', 200.00), (3, 'Product 3', 300.00), (4, 'Product 4', 400.00), (5, 'Product 5', 500.00), (6, 'Product 6', 600.00)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = "INSERT INTO invoices (id, `datetime`, total, customer_id) VALUES (1, '2021-01-01 00:00:00', 1400, 1)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = `
				INSERT INTO sales (id, quantity, invoice_id, product_id) VALUES
					(1, 9, 1, 1),
					(2, 10, 1, 2),
					(3, 12, 1, 3),
					(4, 15, 1, 4),
					(5, 6, 1, 5),
					(6, 3, 1, 6)
			`
			_, err = db.Exec(query)
			require.NoError(t, err)

			return nil
		}()

		require.NoError(t, err)

		// act
		prod, err := r.GetTopProducts()

		// assert
		totalProd := 5
		productsExpected := []internal.Product{
			{Id: 4, ProductAttributes: internal.ProductAttributes{Description: "Product 4", Price: 400.00}},
			{Id: 3, ProductAttributes: internal.ProductAttributes{Description: "Product 3", Price: 300.00}},
			{Id: 2, ProductAttributes: internal.ProductAttributes{Description: "Product 2", Price: 200.00}},
			{Id: 1, ProductAttributes: internal.ProductAttributes{Description: "Product 1", Price: 100.00}},
			{Id: 5, ProductAttributes: internal.ProductAttributes{Description: "Product 5", Price: 500.00}},
		}

		require.NoError(t, err)
		require.Len(t, prod, totalProd)
		require.Equal(t, productsExpected, prod)
	})
}
