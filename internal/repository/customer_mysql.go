package repository

import (
	"database/sql"

	"app/internal"
)

// NewCustomersMySQL creates new mysql repository for customer entity.
func NewCustomersMySQL(db *sql.DB) *CustomersMySQL {
	repo := &CustomersMySQL{db}
	return repo
}

// CustomersMySQL is the MySQL repository implementation for customer entity.
type CustomersMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all customers from the database.
func (r *CustomersMySQL) FindAll() (c []internal.Customer, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `first_name`, `last_name`, `condition` FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var cs internal.Customer
		// scan the row into the customer
		err := rows.Scan(&cs.Id, &cs.FirstName, &cs.LastName, &cs.Condition)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		c = append(c, cs)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the customer into the database.
func (r *CustomersMySQL) Save(c *internal.Customer) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)",
		(*c).FirstName, (*c).LastName, (*c).Condition,
	)
	if err != nil {
		return err
	}

	// get the last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set the id
	(*c).Id = int(id)

	return
}

// GetTopCustomers returns the top 5 customers by quantity of purchases.
func (r *CustomersMySQL) GetTopCustomers() (c []internal.Customer, err error) {
	// query
	query := `
		SELECT c.id, c.first_name, c.last_name, c.condition
		FROM customers c JOIN invoices i ON (c.id = i.customer_id)
		GROUP BY c.id, c.first_name, c.last_name, c.condition
		ORDER BY SUM(i.total) DESC
		LIMIT 5
	`
	// execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	// iterate over the rows
	for rows.Next() {
		// scan the row into the customer
		var cust internal.Customer
		err := rows.Scan(&cust.Id, &cust.FirstName, &cust.LastName, &cust.Condition)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		c = append(c, cust)
	}

	return
}
