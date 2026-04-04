package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pedwoo/first-go-project/db"
)

func OrderDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		tmpl := template.Must(template.ParseFiles(
			"templates/partials/order_detail.html",
		))

		order, err := db.GetOrderByID(id)
		if err != nil {
			log.Printf("Error fetching order: %v\n", err)
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}

		tmpl.ExecuteTemplate(w, "order_detail.html", order)
	}
}