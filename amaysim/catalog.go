package amaysim

import (
	"encoding/json"
	"fmt"
	"io"
)

var catalog Products

// Products is a slice of Product
type Products struct {
	Catalog []Product
}

/*
InitProducts reads from a Reader and parses JSON format string to initialize Products
*/
func InitProducts(source io.Reader) (*Products, error) {
	var products []Product
	dec := json.NewDecoder(source)
	_, err := dec.Token() // '['
	if err != nil {
		return nil, err
	}
	for dec.More() {
		var product Product
		err := dec.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	catalog = Products{products}
	return &catalog, nil
}

/*
PrintCatalog prints each Product to terminal
*/
func (p Products) PrintCatalog() {
	display := "ID\tProduct Name\t\tPrice\n"
	i := 1
	for _, product := range p.Catalog {
		display += fmt.Sprintf("%d\t%s\t\t$%.02f\n", i, product.Name, product.Price)
		i++
	}
	fmt.Println(display)
}

/*
Lookup is used to search for a product given the product code
*/
func (p Products) Lookup(code string) (*Product, error) {
	for _, product := range p.Catalog {
		if code == product.Code {
			return &product, nil
		}
	}
	return nil, fmt.Errorf("Error: Product(code='%s') is missing", code)
}
