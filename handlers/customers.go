// handlers/customers.go
package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/pedwoo/first-go-project/db"
)

func CustomersIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout/base.html",
			"templates/pages/customers.html",
		))

		customers, err := db.GetAllCustomers()
		if err != nil {
			log.Printf("Error fetching customers: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Title":     "Customers",
			"Customers": customers,
		}

		tmpl.ExecuteTemplate(w, "base.html", data)
	}
}

func CustomersSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/partials/customer_rows.html",
		))

		query := r.URL.Query().Get("q")
		customers, err := db.SearchCustomers(query)
		if err != nil {
			log.Printf("Error searching customers: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "customer_rows.html", customers)
	}
}

func CustomersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if err := db.DeleteCustomer(id); err != nil {
			log.Printf("Error deleting customer: %v\n", err)
			w.Header().Set("HX-Reswap", "none")
			http.Error(w, "Cannot delete a customer that has orders.", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}