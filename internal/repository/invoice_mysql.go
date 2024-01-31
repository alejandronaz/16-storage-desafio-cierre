package repository

import (
	"database/sql"
	"fmt"
	"strconv"

	"app/internal"
)

// NewInvoicesMySQL creates new mysql repository for invoice entity.
func NewInvoicesMySQL(db *sql.DB) *InvoicesMySQL {
	return &InvoicesMySQL{db}
}

// InvoicesMySQL is the MySQL repository implementation for invoice entity.
type InvoicesMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all invoices from the database.
func (r *InvoicesMySQL) FindAll() (i []internal.Invoice, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `datetime`, `total`, `customer_id` FROM invoices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var iv internal.Invoice
		// scan the row into the invoice
		err := rows.Scan(&iv.Id, &iv.Datetime, &iv.Total, &iv.CustomerId)
		if err != nil {
			return nil, err
		}
		// append the invoice to the slice
		i = append(i, iv)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the invoice into the database.
func (r *InvoicesMySQL) Save(i *internal.Invoice) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO invoices (`datetime`, `total`, `customer_id`) VALUES (?, ?, ?)",
		(*i).Datetime, (*i).Total, (*i).CustomerId,
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
	(*i).Id = int(id)

	return
}

// UpdateTotalPrice updates the total price of all invoices.
func (r *InvoicesMySQL) UpdateTotalPrice() (err error) {

	// query to update the total price of all invoices
	queryTotal := `
		UPDATE invoices i 
		SET total = (
			SELECT SUM(s.quantity * p.price) 
			FROM sales s JOIN products p ON (s.product_id = p.id) 
			WHERE s.invoice_id = i.id
		)
	`
	_, err = r.db.Exec(queryTotal)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}

// GetTotalPriceByCondition returns the total price of all invoices that match the condition
func (i *InvoicesMySQL) GetTotalPriceByCondition(condition int) (total float64, err error) {

	// query to retrieve the total price by customer condition
	query := `
		SELECT SUM(i.total)
		FROM customers c 
			JOIN invoices i ON (c.id = i.customer_id)
		WHERE c.condition = ?
	`

	// execute the query
	row := i.db.QueryRow(query, condition)
	if row.Err() != nil {
		return 0, row.Err()
	}

	// scan the row into the total
	err = row.Scan(&total)
	if err != nil {
		return 0, err
	}

	// round total to 2 decimals
	twoDec := fmt.Sprintf("%.2f", total)
	total, _ = strconv.ParseFloat(twoDec, 64)

	return total, nil

}
