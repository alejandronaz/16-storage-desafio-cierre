package repository_test

import (
	"app/internal/repository"
	"database/sql"
	"testing"

	txdb "github.com/DATA-DOG/go-txdb"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func init() {
	cfg := &mysql.Config{
		User:   "root",
		Passwd: "bootcamp",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "fantasy_products",
	}
	txdb.Register("testDriver", "mysql", cfg.FormatDSN())
}

func TestGetTotalPriceByCondition_MySQL(t *testing.T) {
	t.Run("Monto total redondeado a 2 decimales por condition 1 del customer", func(t *testing.T) {
		// arrange
		db, err := sql.Open("testDriver", "")
		require.NoError(t, err)
		defer db.Close()

		r := repository.NewInvoicesMySQL(db)

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

			query := "INSERT INTO customers (id, first_name, last_name, `condition`) VALUES	(1, 'John', 'Doe', 1),	(2, 'Jane', 'Doe', 1)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = " INSERT INTO products (id, `description`, price) VALUES (1, 'Product 1', 100.00), (2, 'Product 2', 200.00), (3, 'Product 3', 300.00), (4, 'Product 4', 400.00) "
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = "INSERT INTO invoices (id, `datetime`, total, customer_id) VALUES (1, '2021-01-01 00:00:00', 1400, 1), (2, '2021-01-01 00:00:00', 2000, 2)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = `
				INSERT INTO sales (id, quantity, invoice_id, product_id) VALUES
					(1, 1, 1, 1),
					(2, 2, 1, 2),
					(3, 3, 1, 3),
					(4, 4, 2, 4),
					(5, 2, 2, 2)
			`
			_, err = db.Exec(query)
			require.NoError(t, err)

			return nil
		}()

		require.NoError(t, err)

		// act
		total, err := r.GetTotalPriceByCondition(1)

		// assert
		require.NoError(t, err)
		totalExpected := 3400.0
		require.Equal(t, totalExpected, total)

	})
	t.Run("Monto total redondeado a 2 decimales por condition 0 del customer", func(t *testing.T) {
		// arrange
		db, err := sql.Open("testDriver", "")
		require.NoError(t, err)
		defer db.Close()

		r := repository.NewInvoicesMySQL(db)

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

			query := "INSERT INTO customers (id, first_name, last_name, `condition`) VALUES	(1, 'John', 'Doe', 0),	(2, 'Jane', 'Doe', 0)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = " INSERT INTO products (id, `description`, price) VALUES (1, 'Product 1', 100.00), (2, 'Product 2', 200.00), (3, 'Product 3', 300.00), (4, 'Product 4', 400.00) "
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = "INSERT INTO invoices (id, `datetime`, total, customer_id) VALUES (1, '2021-01-01 00:00:00', 1400, 1), (2, '2021-01-01 00:00:00', 2000, 2)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = `
					INSERT INTO sales (id, quantity, invoice_id, product_id) VALUES
						(1, 1, 1, 1),
						(2, 2, 1, 2),
						(3, 3, 1, 3),
						(4, 4, 2, 4),
						(5, 2, 2, 2)
					`
			_, err = db.Exec(query)
			require.NoError(t, err)

			return nil
		}()

		require.NoError(t, err)

		// act
		total, err := r.GetTotalPriceByCondition(0)

		// assert
		require.NoError(t, err)
		totalExpected := 3400.0
		require.Equal(t, totalExpected, total)

	})
}
