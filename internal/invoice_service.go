package internal

// ServiceInvoice is the interface that wraps the basic methods that an invoice service should implement.
type ServiceInvoice interface {
	// FindAll returns all invoices
	FindAll() (i []Invoice, err error)
	// Save saves an invoice
	Save(i *Invoice) (err error)
	// UpdateTotalPrice updates the total price of all invoices
	UpdateTotalPrice() (err error)
	// GetTotalPriceByCondition returns the total price of all invoices that match the condition
	GetTotalPriceByCondition(condition int) (total float64, err error)
}
