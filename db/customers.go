package db

import (
	"context"
)

type Customer struct {
	CustomerID   string
	CompanyName  string
	ContactName  string
	ContactTitle string
	Address      string
	City         string
	Region       string
	PostalCode   string
	Country      string
	Phone        string
	Fax          string
}

func GetAllCustomers() ([]Customer, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT customer_id, company_name, contact_name, contact_title,
		       address, city, COALESCE(region, ''), COALESCE(postal_code, ''),
		       country, phone, COALESCE(fax, '')
		FROM customers
		ORDER BY company_name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		err := rows.Scan(
			&c.CustomerID, &c.CompanyName, &c.ContactName, &c.ContactTitle,
			&c.Address, &c.City, &c.Region, &c.PostalCode,
			&c.Country, &c.Phone, &c.Fax,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}