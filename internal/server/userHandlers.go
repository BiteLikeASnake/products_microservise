package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/call-me-snake/products_microservise/internal/convert"
	"github.com/call-me-snake/products_microservise/internal/model"
	"github.com/gorilla/mux"
)

func getProductsListUserHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prodList, err := db.GetUserProductsList()
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductsListUserHandler: %s", err.Error())
			return
		}

		if prodList == nil {
			http.Error(w, "No products present", http.StatusNotFound)
			return
		}
		jprodList := make([]model.JsonProduct, 0, len(prodList))
		for _, val := range prodList {
			jprodList = append(jprodList, convert.ConvertToJsonProduct(val))
		}

		res, err := json.Marshal(jprodList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}
func getCategoriesListUserHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		catList, err := db.GetCategoriesByFilter(map[string]interface{}{"category_active": true})
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		if catList == nil {
			http.Error(w, "No categories present", http.StatusNotFound)
			return
		}

		ucatList := make([]model.UserCategory, 0, len(catList))

		for _, val := range catList {
			ucatList = append(ucatList, convert.ConvertToUserCategory(val))
		}

		res, err := json.Marshal(ucatList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func getProductsInCategoryUserHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ucat := &model.UserCategory{}
		err := json.NewDecoder(r.Body).Decode(ucat)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		prodList, err := db.GetUserProductsListFilteredByCategoryId(ucat.CategoryId)
		if prodList == nil {
			http.Error(w, fmt.Sprintf("No products in category %s", ucat.CategoryName), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductsInCategoryUserHandler: %s", err.Error())
			return
		}

		jprodList := make([]model.JsonProduct, 0, len(prodList))

		for _, val := range prodList {
			jprodList = append(jprodList, convert.ConvertToJsonProduct(val))
		}
		res, err := json.Marshal(jprodList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductsInCategoryUserHandler: %s", err.Error())
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func getProductsByFilterUserHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uprod := &model.FilterUProduct{}
		err := json.NewDecoder(r.Body).Decode(uprod)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		uprodList, err := db.GetProductsByUserFilter(*uprod)

		if uprodList == nil {
			http.Error(w, fmt.Sprintf("No products filtered present"), http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductsByFilterUserHandler: %s", err.Error())
			return
		}

		jprodList := make([]model.JsonProduct, 0, len(uprodList))

		for _, val := range uprodList {
			jprodList = append(jprodList, convert.ConvertToJsonProduct(val))
		}
		res, err := json.Marshal(jprodList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductsInCategoryUserHandler: %s", err.Error())
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)

	}
}

func getProductByIdUserHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		id64 := int64(id)
		prod, err := db.GetUserProductById(id64)

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("GetProductByIdUserHandler: %s", err.Error())
			return
		}
		if prod == nil {
			http.Error(w, fmt.Sprintf("No product with id %d present", id64), http.StatusNotFound)
			return
		}

		jprod := convert.ConvertToJsonProduct(*prod)

		res, err := json.Marshal(jprod)

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("GetProductByIdUserHandler: %s", err.Error())
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func updateProductUserHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uprod := &model.UserUpdateProduct{}
		err := json.NewDecoder(r.Body).Decode(uprod)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}
		if uprod.AmountPurchased <= 0 {
			http.Error(w, fmt.Sprintf("Неверное количество товара"), http.StatusBadRequest)
			return
		}
		currentprod, err := db.GetProductById(uprod.ProductId)

		if currentprod.Product_active == false {
			http.Error(w, fmt.Sprintf("Товар не найден"), http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("updateProductByIdUserHandler: %s", err.Error())
			return
		}

		if uprod.AmountPurchased > currentprod.Product_quantity {
			http.Error(w, fmt.Sprintf("Покупка невозможна, недостаточно товара на складе"), http.StatusBadRequest)
			return
		}
		productsRemains := currentprod.Product_quantity - uprod.AmountPurchased

		isUpdated, err := db.UpdateProductById(uprod.ProductId, map[string]interface{}{"product_quantity": productsRemains})

		if isUpdated == false || err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			if err != nil {
				log.Printf("updateProductByIdUserHandler: %s", err.Error())
			}
			return
		}

		fmt.Fprintf(w, "Покупка совершена успешно.")
	}
}
