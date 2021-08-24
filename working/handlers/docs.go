// Package classification of Product API
//
// Documentation for Product API that resembles Coffee Objects
//
// Schemes: http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//swagger:meta
package handlers

import "microservices/working/data"

// List of all products
// swagger:response productsResponse
type productsResponseWrapper struct {
	// All current products
	// in: body
	Body []data.Product
}

// Data structure representing a single product
// swagger:response productResponse
type productResponseWrapper struct {
	// Newly created product
	// in: body
	Body data.Product
}

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}
