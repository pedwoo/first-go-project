package db

import (
	"context"

	"github.com/google/uuid"
)

type Customer struct {
	CustomerID   uuid.UUID
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

func SearchCustomers(query string) ([]Customer, error) {
	rows, err := Pool.Query(context.Background(), `
		SELECT customer_id, company_name, contact_name, contact_title,
		       address, city, COALESCE(region, ''), COALESCE(postal_code, ''),
		       country, phone, COALESCE(fax, '')
		FROM customers
		WHERE LOWER(company_name) LIKE LOWER($1)
		OR LOWER(contact_name) LIKE LOWER($1)
		OR LOWER(country) LIKE LOWER($1)
		ORDER BY company_name
	`, "%"+query+"%")
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

func DeleteCustomer(id uuid.UUID) error {
	_, err := Pool.Exec(context.Background(), `
		DELETE FROM customers WHERE customer_id = $1
	`, id)
	return err
}