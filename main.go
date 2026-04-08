package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/pedwoo/first-go-project/db"
	"github.com/pedwoo/first-go-project/handlers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	db.Connect()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("").Funcs(funcMap).ParseFiles(
			"templates/layout/base.html",
			"templates/pages/dashboard.html",
		))

		stats, err := db.GetDashboardStats()
		if err != nil {
			log.Printf("Error fetching stats: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		recentOrders, err := db.GetRecentOrders()
		if err != nil {
			log.Printf("Error fetching recent orders: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		topProducts, err := db.GetTopProducts()
		if err != nil {
			log.Printf("Error fetching top products: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		ordersByCountry, err := db.GetOrdersByCountry()
		if err != nil {
			log.Printf("Error fetching orders by country: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		data := map[string]any{
			"Title":           "Dashboard",
			"Stats":           stats,
			"RecentOrders":    recentOrders,
			"TopProducts":     topProducts,
			"OrdersByCountry": ordersByCountry,
		}
		tmpl.ExecuteTemplate(w, "base.html", data)
	})
	r.Get("/orders/{id}", handlers.OrderDetail())

	r.Get("/customers", handlers.CustomersIndex())
	r.Get("/customers/search", handlers.CustomersSearch())
	r.Delete("/customers/{id}", handlers.CustomersDelete())

	r.Get("/orders", handlers.OrdersIndex())
	// r.Get("/orders/search", handlers.OrdersSearch())
	// r.Get("/orders/{id}", handlers.OrderDetail())

	// r.Get("/products", handlers.ProductsIndex())
	// r.Get("/products/search", handlers.ProductsSearch())

	// r.Get("/employees", handlers.EmployeesIndex())
	// r.Get("/employees/search", handlers.EmployeesSearch())

	// r.Get("/suppliers", handlers.SuppliersIndex())
	// r.Get("/suppliers/search", handlers.SuppliersSearch())

	// r.Get("/categories", handlers.CategoriesIndex())
	// r.Get("/categories/search", handlers.CategoriesSearch())

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
