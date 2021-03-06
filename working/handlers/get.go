package handlers

import (
	"context"
	protos "gRPC/currency/protos/currency"
	"microservices/working/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Return a list of products from the database
// responses:
//		200: productsResponse

// ListAll returns all the products from the data store
func (p *Products) ListAll(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := data.ToJSON(lp, rw)
	if err != nil {
		http.Error(rw, "unable to convert it to json", http.StatusInternalServerError)
	}
}

// swagger:route GET /products/{id} products listSingle
// Return a list of products from the database
// responses:
//		202: productResponse
// 		404: errorResponse

// ListSingle handles GET Request
func (p *Products) ListSingle(rw http.ResponseWriter, r *http.Request) {
	id := getProductId(r)
	p.l.Println("[DEBUG] get record id", id)
	prod, err := data.GetProductById(id)
	p.l.Println("prod", prod)

	switch err {
	case nil:

	case data.ErrProductNotFound:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusNotFound)
		data.ToJSON(prod, rw)
		return

	default:
		p.l.Println("[ERROR] fetching product", err)
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	// get exchange rate
	rr := &protos.RateRequest{
		Base:        protos.Currencies(protos.Currencies_value["EUR"]),
		Destination: protos.Currencies(protos.Currencies_value["USD"]),
	}

	resp, err := p.cc.GetRate(context.Background(), rr)
	if err != nil {
		p.l.Println("[Error] error getting new rate", err)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	p.l.Printf("Resp %#v", resp)

	prod.Price = prod.Price * resp.Rate

	err = data.ToJSON(prod, rw)
	if err != nil {
		p.l.Println("[ERROR] serializing product", err)
	}
}
