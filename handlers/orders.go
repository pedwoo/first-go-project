package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pedwoo/first-go-project/db"
)

func OrdersIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/layout/base.html",
			"templates/pages/orders.html",
		))

		orders, err := db.GetAllOrders()
		if err != nil {
			log.Printf("Error fetching orders: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Title":  "Orders",
			"Orders": orders,
		}

		tmpl.ExecuteTemplate(w, "base.html", data)
	}
}

func OrdersSearch() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(
			"templates/partials/order_rows.html",
		))

		query := r.URL.Query().Get("q")
		orders, err := db.SearchOrders(query)
		if err != nil {
			log.Printf("Error searching orders: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "order_rows.html", orders)
	}
}

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
