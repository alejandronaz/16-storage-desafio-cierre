package handler

import (
	"net/http"
	"strconv"

	"app/internal"
	"app/platform/web/request"
	"app/platform/web/response"

	"github.com/go-chi/chi/v5"
)

// NewInvoicesDefault returns a new InvoicesDefault
func NewInvoicesDefault(sv internal.ServiceInvoice) *InvoicesDefault {
	return &InvoicesDefault{sv: sv}
}

// InvoicesDefault is a struct that returns the invoice handlers
type InvoicesDefault struct {
	// sv is the invoice's service
	sv internal.ServiceInvoice
}

// InvoiceJSON is a struct that represents a invoice in JSON format
type InvoiceJSON struct {
	Id         int     `json:"id"`
	Datetime   string  `json:"datetime"`
	Total      float64 `json:"total"`
	CustomerId int     `json:"customer_id"`
}

// GetAll returns all invoices
func (h *InvoicesDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// ...

		// process
		i, err := h.sv.FindAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error getting invoices")
			return
		}

		// response
		// - serialize
		ivJSON := make([]InvoiceJSON, len(i))
		for ix, v := range i {
			ivJSON[ix] = InvoiceJSON{
				Id:         v.Id,
				Datetime:   v.Datetime,
				Total:      v.Total,
				CustomerId: v.CustomerId,
			}
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "invoices found",
			"data":    ivJSON,
		})
	}
}

// RequestBodyInvoice is a struct that represents the request body for a invoice
type RequestBodyInvoice struct {
	Datetime   string  `json:"datetime"`
	Total      float64 `json:"total"`
	CustomerId int     `json:"customer_id"`
}

// Create creates a new invoice
func (h *InvoicesDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var reqBody RequestBodyInvoice
		err := request.JSON(r, &reqBody)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error parsing request body")
			return
		}

		// process
		// - deserialize
		i := internal.Invoice{
			InvoiceAttributes: internal.InvoiceAttributes{
				Datetime:   reqBody.Datetime,
				Total:      reqBody.Total,
				CustomerId: reqBody.CustomerId,
			},
		}
		// - save
		err = h.sv.Save(&i)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error saving invoice")
			return
		}

		// response
		// - serialize
		iv := InvoiceJSON{
			Id:         i.Id,
			Datetime:   i.Datetime,
			Total:      i.Total,
			CustomerId: i.CustomerId,
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "invoice created",
			"data":    iv,
		})
	}
}

// UpdateTotalPrice updates the total price of all invoices
func (h *InvoicesDefault) UpdateTotalPrice(w http.ResponseWriter, r *http.Request) {
	// call the service
	err := h.sv.UpdateTotalPrice()
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error updating total price")
		return
	}

	// response
	response.JSON(w, http.StatusOK, map[string]any{
		"message": "total price updated",
	})
}

// GetTotalPriceByCondition returns the total price of all invoices that match the condition
func (h *InvoicesDefault) GetTotalPriceByCondition(w http.ResponseWriter, r *http.Request) {
	// get path param
	condition := chi.URLParam(r, "condition")

	// parse to int
	conditionInt, err := strconv.Atoi(condition)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "error parsing condition; must be integer")
		return
	}

	// call the service
	total, err := h.sv.GetTotalPriceByCondition(conditionInt)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error getting total price")
		return
	}

	// response
	response.JSON(w, http.StatusOK, map[string]any{
		"message": "success",
		"data":    total,
	})
}
