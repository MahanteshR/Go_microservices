package handlers

import "net/http"

type Hello struct{

}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	
}
