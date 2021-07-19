package handlers

import (
	"net/http"

	"github.com/SaishNaik/microservices_jn/product-api/data"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//	200: productsResponse

// ListAll handles GET requests and returns all current products
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	
	p.l.Debug("Get all records")
	rw.Header().Add("Content-Type", "application/json")
	
	currency:= r.URL.Query().Get("currency")
	
	prods,err := p.productDB.GetProducts(currency)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}


	err = data.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Serializing product","error",err)
		return
	}
}

// swagger:route GET /products/{id} products listSingleProduct
// Return a list of products from the database
// responses:
//	200: productResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductID(r)
	p.l.Debug("Get record","id", id)

	currency:= r.URL.Query().Get("currency")
	prod, err := p.productDB.GetProductByID(id,currency)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Error("Fetching product", "error",err)

		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Error("Fetching product", "error",err)

		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}


	


	err = data.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Error("Unable serializing product", "error",err)
	}
}
