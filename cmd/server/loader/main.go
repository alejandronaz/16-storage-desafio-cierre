package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func main() {
	// app
	// - config
	cfg := mysql.Config{
		User:   "root",
		Passwd: "bootcamp",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "fantasy_products",
	}

	// - db: init
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return
	}
	// - db: ping
	err = db.Ping()
	if err != nil {
		return
	}

	// - loader: init
	loader := NewLoaderDefault(db)
	// - loader: load from json
	err = loader.LoadFromJSON()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Loaded from JSON")
}

func NewLoaderDefault(db *sql.DB) *Loader {
	return &Loader{
		db: db,
	}
}

type Loader struct {
	db *sql.DB
	tx *sql.Tx
}

func (l *Loader) LoadFromJSON() (err error) {

	// begin the transaction
	tx, err := l.db.Begin()
	if err != nil {
		return
	}
	l.tx = tx

	var mysqlErr *mysql.MySQLError

	// ---- CUSTOMERS ----
	err = l.LoadCustomers("docs/db/json/customers.json")
	if err != nil {

		// if the error is not a duplicate key error, return it
		if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
			tx.Rollback()
			return err
		}

	}

	// ---- PRODUCTS ----
	err = l.LoadProducts("docs/db/json/products.json")
	if err != nil {

		// if the error is not a duplicate key error, return it
		if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
			tx.Rollback()
			return err
		}

	}

	// ---- INVOICES ----
	err = l.LoadInvoices("docs/db/json/invoices.json")
	if err != nil {

		// if the error is not a duplicate key error, return it
		if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
			tx.Rollback()
			return err
		}

	}

	// ---- SALES ----
	err = l.LoadSales("docs/db/json/sales.json")
	if err != nil {

		// if the error is not a duplicate key error, return it
		if errors.As(err, &mysqlErr) && mysqlErr.Number != 1062 {
			tx.Rollback()
			return err
		}

	}

	// commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	return nil

}

// CustomerJSON is the struct that represents the customer in the JSON file.
type CustomerJSON struct {
	ID        int    `json:"id"`
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Condition int    `json:"condition"`
}

// LoadCustomers loads customers from a JSON file.
func (l *Loader) LoadCustomers(filename string) (err error) {
	// open the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.New("could not open the file")
	}

	// unmarshal the data
	var customers []CustomerJSON
	if err = json.Unmarshal(data, &customers); err != nil {
		return errors.New("could not unmarshal the data")
	}

	// prepare the statement
	stmt, err := l.tx.Prepare("INSERT INTO customers (`id`, `first_name`, `last_name`, `condition`) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	// iterate over the customers and insert them into the database
	for _, c := range customers {
		// execute the statement
		_, err = stmt.Exec(c.ID, c.FirstName, c.LastName, c.Condition)
		if err != nil {
			return err
		}
	}

	return nil

}

// ProductJSON is the struct that represents the product in the JSON file.
type ProductJSON struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

// LoadProducts loads products from a JSON file.
func (l *Loader) LoadProducts(filename string) (err error) {
	// open the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.New("could not open the file")
	}

	// unmarshal the data
	var products []ProductJSON
	if err = json.Unmarshal(data, &products); err != nil {
		return errors.New("could not unmarshal the data")
	}

	// prepare the statement
	stmt, err := l.tx.Prepare("INSERT INTO products (`id`, `description`, `price`) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	// iterate over the products and insert them into the database
	for _, p := range products {
		// execute the statement
		_, err = stmt.Exec(p.ID, p.Description, p.Price)
		if err != nil {
			return err
		}
	}

	return nil

}

// InvoiceJSON is the struct that represents the invoice in the JSON file.
type InvoiceJSON struct {
	ID         int     `json:"id"`
	Datetime   string  `json:"datetime"`
	Total      float64 `json:"total"`
	CustomerID int     `json:"customer_id"`
}

// LoadInvoices loads invoices from a JSON file.
func (l *Loader) LoadInvoices(filename string) (err error) {
	// open the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.New("could not open the file")
	}

	// unmarshal the data
	var invoices []InvoiceJSON
	if err = json.Unmarshal(data, &invoices); err != nil {
		return errors.New("could not unmarshal the data")
	}

	// prepare the statement
	stmt, err := l.tx.Prepare("INSERT INTO invoices (`id`, `datetime`, `total`, `customer_id`) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	// iterate over the invoices and insert them into the database
	for _, i := range invoices {
		// execute the statement
		_, err = stmt.Exec(i.ID, i.Datetime, i.Total, i.CustomerID)
		if err != nil {
			return err
		}
	}

	return nil

}

// SaleJSON is the struct that represents the sale in the JSON file.
type SaleJSON struct {
	ID        int `json:"id"`
	Quantity  int `json:"quantity"`
	ProductID int `json:"product_id"`
	InvoiceID int `json:"invoice_id"`
}

// LoadSales loads sales from a JSON file.
func (l *Loader) LoadSales(filename string) (err error) {
	// open the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return errors.New("could not open the file")
	}

	// unmarshal the data
	var sales []SaleJSON
	if err = json.Unmarshal(data, &sales); err != nil {
		return errors.New("could not unmarshal the data")
	}

	// prepare the statement
	stmt, err := l.tx.Prepare("INSERT INTO sales (`id`, `quantity`, `product_id`, `invoice_id`) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}

	// iterate over the sales and insert them into the database
	for _, s := range sales {
		// execute the statement
		_, err = stmt.Exec(s.ID, s.Quantity, s.ProductID, s.InvoiceID)
		if err != nil {
			return err
		}
	}

	return nil

}
