package db

import "context"

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