package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"microservices/working/data"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
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

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "unable to convert it to json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST Product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert ID", http.StatusBadRequest)
		return
	}

	p.l.Println("Update product handler")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	err = data.UpdateProduct(id, &prod)
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

		err := prod.FromJSON(r.Body)
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
