package db

import "context"

type Order struct {
	OrderID      string
	CustomerName string
	EmployeeName string
	OrderDate    string
	RequiredDate string
	ShippedDate  string
	ShipperName  string
	Freight      float64
	ShipCity     string
	ShipCountry  string
	Shipped      bool
}

type OrderDetail struct {
	OrderID        string
	CustomerName   string
	EmployeeName   string
	OrderDate      string
	RequiredDate   string
	ShippedDate    string
	ShipperName    string
	Freight        float64
	ShipName       string
	ShipAddress    string
	ShipCity       string
	ShipRegion     string
	ShipPostalCode string
	ShipCountry    string
	Shipped        bool
}

func GetOrderByID(id string) (OrderDetail, error) {
	var o OrderDetail

	err := Pool.QueryRow(context.Background(), `
		SELECT
			o.order_id,
			c.company_name,
			e.first_name || ' ' || e.last_name,
			TO_CHAR(o.order_date, 'Mon DD, YYYY'),
			TO_CHAR(o.required_date, 'Mon DD, YYYY'),
			COALESCE(TO_CHAR(o.shipped_date, 'Mon DD, YYYY'), '—'),
			s.company_name,
			o.freight,
			o.ship_name,
			o.ship_address,
			o.ship_city,
			COALESCE(o.ship_region, '—'),
			COALESCE(o.ship_postal_code, '—'),
			o.ship_country,
			o.shipped_date IS NOT NULL
		FROM orders o
		JOIN customers c ON o.customer_id = c.customer_id
		JOIN employees e ON o.employee_id = e.employee_id
		JOIN shippers s ON o.ship_via = s.shipper_id
		WHERE o.order_id = $1
	`, id).Scan(
		&o.OrderID, &o.CustomerName, &o.EmployeeName,
		&o.OrderDate, &o.RequiredDate, &o.ShippedDate,
		&o.ShipperName, &o.Freight, &o.ShipName,
		&o.ShipAddress, &o.ShipCity, &o.ShipRegion,
		&o.ShipPostalCode, &o.ShipCountry, &o.Shipped,
	)

	return o, err
}

func GetAllOrders() ([]Order, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT
			o.order_id,
			c.company_name,
			e.first_name || ' ' || e.last_name,
			TO_CHAR(o.order_date, 'Mon DD, YYYY'),
			TO_CHAR(o.required_date, 'Mon DD, YYYY'),
			COALESCE(TO_CHAR(o.shipped_date, 'Mon DD, YYYY'), '—'),
			s.company_name,
			o.freight,
			o.ship_city,
			o.ship_country,
			o.shipped_date IS NOT NULL
		FROM orders o
		JOIN customers c ON o.customer_id = c.customer_id
		JOIN employees e ON o.employee_id = e.employee_id
		JOIN shippers s ON o.ship_via = s.shipper_id
		ORDER BY o.order_date DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.OrderID, &o.CustomerName, &o.EmployeeName,
			&o.OrderDate, &o.RequiredDate, &o.ShippedDate,
			&o.ShipperName, &o.Freight, &o.ShipCity,
			&o.ShipCountry, &o.Shipped,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}

func SearchOrders(query string) ([]Order, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT
			o.order_id,
			c.company_name,
			e.first_name || ' ' || e.last_name,
			TO_CHAR(o.order_date, 'Mon DD, YYYY'),
			TO_CHAR(o.required_date, 'Mon DD, YYYY'),
			COALESCE(TO_CHAR(o.shipped_date, 'Mon DD, YYYY'), '—'),
			s.company_name,
			o.freight,
			o.ship_city,
			o.ship_country,
			o.shipped_date IS NOT NULL
		FROM orders o
		JOIN customers c ON o.customer_id = c.customer_id
		JOIN employees e ON o.employee_id = e.employee_id
		JOIN shippers s ON o.ship_via = s.shipper_id
		WHERE LOWER(c.company_name) LIKE LOWER($1)
		OR LOWER(o.ship_country) LIKE LOWER($1)
		OR LOWER(o.ship_city) LIKE LOWER($1)
		ORDER BY o.order_date DESC
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		err := rows.Scan(
			&o.OrderID, &o.CustomerName, &o.EmployeeName,
			&o.OrderDate, &o.RequiredDate, &o.ShippedDate,
			&o.ShipperName, &o.Freight, &o.ShipCity,
			&o.ShipCountry, &o.Shipped,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
