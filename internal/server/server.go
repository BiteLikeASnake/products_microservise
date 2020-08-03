package server

import (
	"fmt"
	"net/http"

	"github.com/BiteLikeASnake/products_microservise/internal/model"
	"github.com/gorilla/mux"
)

type Connector struct {
	router *mux.Router
	addr   string
}

func New(addr string) *Connector {
	c := &Connector{}
	c.router = mux.NewRouter()
	c.addr = addr
	return c
}

func (c *Connector) ExecuteHandlers(db model.IDb) {
	s_sup := c.router.PathPrefix("/support").Subrouter()
	s_sup.Use(authAdmin)

	s_sup.HandleFunc("/products/", getProductsListAdminHandler(db)).Methods("GET")
	s_sup.HandleFunc("/products/{id:[0-9]+}", getProductByIdAdminHandler(db)).Methods("GET")
	s_sup.HandleFunc("/products/", setProductHandler(db)).Methods("POST")
	s_sup.HandleFunc("/products/", updateProductAdminHandler(db)).Methods("PUT")
	s_sup.HandleFunc("/products/", deleteProductByIdHandler(db)).Methods("DELETE")

	s_sup.HandleFunc("/categories/", getCategoriesListAdminHandler(db)).Methods("GET")
	s_sup.HandleFunc("/categories/{id:[0-9]+}", getCategoryByIdAdminHandler(db)).Methods("GET")
	s_sup.HandleFunc("/categories/products", getProductsInCategoryAdminHandler(db)).Methods("POST")
	s_sup.HandleFunc("/categories/", setCategoryHandler(db)).Methods("POST")
	s_sup.HandleFunc("/categories/", updateCategoryHandler(db)).Methods("PUT")
	s_sup.HandleFunc("/categories/", deleteCategoryByIdHandler(db)).Methods("DELETE")

	s_cust := c.router.PathPrefix("/customer").Subrouter()
	s_cust.Use(authUser)
	s_cust.HandleFunc("/products/", getProductsListUserHandler(db)).Methods("GET")
	s_cust.HandleFunc("/products/{id:[0-9]+}", getProductByIdUserHandler(db)).Methods("GET")
	s_cust.HandleFunc("/products/", updateProductUserHandler(db)).Methods("PUT")

	s_cust.HandleFunc("/products/filter/", getProductsByFilterUserHandler(db)).Methods("POST")

	s_cust.HandleFunc("/categories/", getCategoriesListUserHandler(db)).Methods("GET")
	s_cust.HandleFunc("/categories/products", getProductsInCategoryUserHandler(db)).Methods("POST")

	c.router.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "alive") })
}

func (c *Connector) Start() {
	http.ListenAndServe(c.addr, c.router)
}
