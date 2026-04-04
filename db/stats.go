package db

import "context"

type DashboardStats struct {
	TotalCustomers int
	TotalOrders    int
	TotalProducts  int
	TotalRevenue   float64
}

type RecentOrder struct {
	OrderID     string
	CustomerName string
	OrderDate   string
	Freight     float64
	Shipped     bool
}

type TopProduct struct {
	ProductName string
	Revenue     float64
}

type OrdersByCountry struct {
	Country string
	Total   int
}

func GetDashboardStats() (DashboardStats, error) {
	var stats DashboardStats

	err := Pool.QueryRow(context.Background(), `
		SELECT
			(SELECT COUNT(*) FROM customers),
			(SELECT COUNT(*) FROM orders),
			(SELECT COUNT(*) FROM products),
			(SELECT COALESCE(SUM(od.unit_price * od.quantity * (1 - od.discount)), 0)
			 FROM order_details od)
	`).Scan(
		&stats.TotalCustomers,
		&stats.TotalOrders,
		&stats.TotalProducts,
		&stats.TotalRevenue,
	)

	return stats, err
}

func GetRecentOrders() ([]RecentOrder, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT o.order_id, c.company_name, 
		       TO_CHAR(o.order_date, 'Mon DD, YYYY'),
		       o.freight,
		       o.shipped_date IS NOT NULL
		FROM orders o
		JOIN customers c ON o.customer_id = c.customer_id
		ORDER BY o.order_date DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []RecentOrder
	for rows.Next() {
		var o RecentOrder
		err := rows.Scan(&o.OrderID, &o.CustomerName, &o.OrderDate, &o.Freight, &o.Shipped)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func GetTopProducts() ([]TopProduct, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT p.product_name,
		       SUM(od.unit_price * od.quantity * (1 - od.discount)) as revenue
		FROM order_details od
		JOIN products p ON od.product_id = p.product_id
		GROUP BY p.product_name
		ORDER BY revenue DESC
		LIMIT 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []TopProduct
	for rows.Next() {
		var p TopProduct
		err := rows.Scan(&p.ProductName, &p.Revenue)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetOrdersByCountry() ([]OrdersByCountry, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT ship_country, COUNT(*) as total
		FROM orders
		GROUP BY ship_country
		ORDER BY total DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []OrdersByCountry
	for rows.Next() {
		var c OrdersByCountry
		err := rows.Scan(&c.Country, &c.Total)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}
	return countries, nil
}