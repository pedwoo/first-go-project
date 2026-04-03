// handlers/customers.go
package handlers

import (
	"html/template"
	"log"
	"net/http"

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