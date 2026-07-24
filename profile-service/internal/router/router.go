package router

import (
	"github.com/gorilla/mux"
)

type RouterRegistrar interface {
	RegisterRoutes(r *mux.Router)
}

func NewRouter(handlers ...RouterRegistrar) *mux.Router {
	r := mux.NewRouter()
	for _, h := range handlers {
		h.RegisterRoutes(r)
	}
	return r
}
