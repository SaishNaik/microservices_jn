package data

import (
	"context"
	"fmt"

	protos "github.com/SaishNaik/microservices_jn/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("Product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"` // Unique identifier for the product

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"required,gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"sku"`
}

// Products defines a slice of Product
type Products []*Product

type ProductsDB struct{
	currency protos.CurrencyClient
	log hclog.Logger
}

func NewProductsDB(c protos.CurrencyClient,log hclog.Logger) *ProductsDB{
	return &ProductsDB{c,log}
}

func (p *ProductsDB) getRate(currency string) (float64,error) {
	rr := protos.RateRequest{
		Base: protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value[currency]),
	}

	resp,err := p.currency.GetRate(context.Background(),&rr)
	if err != nil{
		return 0,err
	}
	return resp.Rate,err
}

// GetProducts returns all products from the database
func (p *ProductsDB) GetProducts(currency string) (Products,error) {
	if(currency == ""){
		return productList,nil
	}

	rate,err := p.getRate(currency)
	if err != nil{
		p.log.Error("Unable to get rate","currency",currency,"error",err)
		return nil,err
	}

	pr := Products{}

	for _,prod:= range productList{
		temp := *prod
		temp.Price *= rate
		pr = append(pr, &temp)
	}

	return pr,nil
	
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func (p *ProductsDB) GetProductByID(id int,currency string) (*Product, error) {
	i := findIndexByProductID(id)
	if id == -1 {
		return nil, ErrProductNotFound
	}

	if(currency == ""){
		return productList[i], nil
	}

	rate,err := p.getRate(currency)
	if err != nil{
		p.log.Error("Unable to get rate","currency",currency,"error",err)
	}

	temp := *productList[i]
	temp.Price *= rate
	
	return &temp,nil
	
}


 

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func (pdb *ProductsDB) UpdateProduct(p Product) error {
	i := findIndexByProductID(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	// update the product in the DB
	productList[i] = &p

	return nil
}

// AddProduct adds a new product to the database
func (pdb *ProductsDB) AddProduct(p Product) {
	// get the next id in sequence
	maxID := productList[len(productList)-1].ID
	p.ID = maxID + 1
	productList = append(productList, &p)
}

// DeleteProduct deletes a product from the database
func (pdb *ProductsDB) DeleteProduct(id int) error {
	i := findIndexByProductID(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1])

	return nil
}

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
func findIndexByProductID(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
	},
	&Product{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
	},
}
