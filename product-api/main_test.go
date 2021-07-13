package main

import (
	"fmt"
	"testing"

	"github.com/SaishNaik/microservices_jn/product-api/sdk/client"
	"github.com/SaishNaik/microservices_jn/product-api/sdk/client/products"
)

func TestOurClient(t *testing.T){
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil,cfg)
	params := products.NewListProductsParams()
	prods,err := c.Products.ListProducts(params)
	if err != nil{
		t.Fatal(err)
	}
	fmt.Printf("%#v",prods.GetPayload()[0])
	t.Fail()
} 

