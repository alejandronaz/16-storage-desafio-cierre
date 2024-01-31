package repository_test

import (
	"app/internal"
	"app/internal/repository"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetTopCustomers(t *testing.T) {
	t.Run("success - get top 5 customers by quantity of purchases", func(t *testing.T) {
		// arrange
		db, err := sql.Open("testDriver", "")
		require.NoError(t, err)
		defer db.Close()

		r := repository.NewCustomersMySQL(db)

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

			query := "INSERT INTO customers (id, first_name, last_name, `condition`) VALUES	(1, 'John', 'Doe', 1),	(2, 'Jane', 'Doe', 1), (3, 'John', 'Doe', 1), (4, 'John', 'Doe', 1), (5, 'John', 'Doe', 1), (6, 'John', 'Doe', 1)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = " INSERT INTO products (id, `description`, price) VALUES (1, 'Product 1', 100.00), (2, 'Product 2', 200.00), (3, 'Product 3', 300.00), (4, 'Product 4', 400.00) "
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = "INSERT INTO invoices (id, `datetime`, total, customer_id) VALUES (1, '2021-01-01 00:00:00', 1400, 1), (2, '2021-01-01 00:00:00', 2000, 2), (3, '2021-01-01 00:00:00', 3000, 3), (4, '2021-01-01 00:00:00', 4000, 4), (5, '2021-01-01 00:00:00', 5000, 5), (6, '2021-01-01 00:00:00', 6000, 6)"
			_, err = db.Exec(query)
			require.NoError(t, err)

			query = `
				INSERT INTO sales (id, quantity, invoice_id, product_id) VALUES
					(1, 14, 1, 1),
					(2, 10, 2, 2),
					(3, 10, 3, 3),
					(4, 10, 4, 4),
					(5, 50, 5, 1),
					(6, 30, 6, 2)
			`
			_, err = db.Exec(query)
			require.NoError(t, err)

			return nil
		}()

		require.NoError(t, err)

		// act
		customers, err := r.GetTopCustomers()

		// assert
		totalCustomers := 5
		customersExpected := []internal.Customer{
			{Id: 6, CustomerAttributes: internal.CustomerAttributes{FirstName: "John", LastName: "Doe", Condition: 1}},
			{Id: 5, CustomerAttributes: internal.CustomerAttributes{FirstName: "John", LastName: "Doe", Condition: 1}},
			{Id: 4, CustomerAttributes: internal.CustomerAttributes{FirstName: "John", LastName: "Doe", Condition: 1}},
			{Id: 3, CustomerAttributes: internal.CustomerAttributes{FirstName: "John", LastName: "Doe", Condition: 1}},
			{Id: 2, CustomerAttributes: internal.CustomerAttributes{FirstName: "Jane", LastName: "Doe", Condition: 1}},
		}

		require.NoError(t, err)
		require.Len(t, customers, totalCustomers)
		require.Equal(t, customersExpected, customers)
	})
}
