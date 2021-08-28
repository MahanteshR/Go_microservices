package handlers

import (
	"context"
	"fmt"
	protos "gRPC/currency/protos/currency"
	"github.com/gorilla/mux"
	"log"
	"microservices/working/data"
	"net/http"
	"strconv"
)

type Products struct {
	l  *log.Logger
	cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, cc protos.CurrencyClient) *Products {
	return &Products{l, cc}
}

func getProductId(r *http.Request) int {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		panic(err)
	}
	return id
}

//func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodGet {
//		p.getProducts(rw, r)
//		return
//	}
//
//	if r.Method == http.MethodPost {
//		p.addProducts(rw, r)
//		return
//	}
//
//	if r.Method == http.MethodPut {
//		reg := regexp.MustCompile(`/([0-9]+)`)
//		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
//
//		if len(g) != 1 {
//			http.Error(rw, "invalid URI", http.StatusBadRequest)
//			return
//		}
//
//		if len(g[0]) != 2 {
//			http.Error(rw, "invalid URI", http.StatusBadRequest)
//			return
//		}
//
//		idString := g[0][1]
//		id, err := strconv.Atoi(idString)
//		if err != nil {
//			http.Error(rw, "invalid URI", http.StatusBadRequest)
//			return
//		}
//
//		p.l.Println("got id", id)
//		p.updateProducts(id, rw, r)
//		return
//	}
//	rw.WriteHeader(http.StatusMethodNotAllowed)
//}

//func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
//	p.l.Println("Handle POST Product")
//	prod := r.Context().Value(KeyProduct{}).(data.Product)
//	data.AddProduct(&prod)
//}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	p.l.Println("Update product handler")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := data.FromJSON(prod, r.Body)
		if err != nil {
			http.Error(rw, "error reading product", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("validation error: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}

type GenericError struct {
	Message string `json:"message"`
}
