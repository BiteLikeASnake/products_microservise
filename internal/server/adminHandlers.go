package server

import (
	"github.com/call-me-snake/products_microservise/internal/convert"
	"github.com/call-me-snake/products_microservise/internal/model"
	response "github.com/call-me-snake/products_microservise/internal/server/model"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/call-me-snake/products_microservise/internal/validate"

	"github.com/gorilla/mux"
)

const InternalErrorMessage = "Внутренняя ошибка сервера"
const BadRequestMessage = "Некорректные входные данные"

//handlers для admin
//для products
func getProductsListAdminHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prodList, err := db.GetProductsList()
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductsListAdminHandler: %s", err.Error())
			return
		}
		if prodList == nil {
			http.Error(w, "No products in base", http.StatusNotFound)
			return
		}

		aprodList := make([]model.AdminProduct, 0, len(prodList))

		for _, val := range prodList {
			aprodList = append(aprodList, convert.ConvertToAdminProduct(val))
		}

		res, err := json.Marshal(aprodList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func getProductByIdAdminHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		id64 := int64(id)
		prod, err := db.GetProductById(id64)

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getProductByIdAdminHandler: %s", err.Error())
			return
		}
		if prod == nil {
			http.Error(w, fmt.Sprintf("No product with id %d in base", id64), http.StatusNotFound)
			return
		}

		aprod := convert.ConvertToAdminProduct(*prod)
		res, err := json.Marshal(aprod)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
		w.WriteHeader(http.StatusFound)
	}
}

func getProductsInCategoryAdminHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cat := &model.AdminCategory{}
		err := json.NewDecoder(r.Body).Decode(cat)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		filter := make(map[string]interface{})
		filter["category_id"] = cat.CategoryId

		prodList, err := db.GetProductsByFilter(filter)

		if prodList == nil {
			http.Error(w, "No products in category", http.StatusNotFound)
			return
		}

		aprodList := make([]model.AdminProduct, 0, len(prodList))

		for _, val := range prodList {
			aprodList = append(aprodList, convert.ConvertToAdminProduct(val))
		}

		res, err := json.Marshal(aprodList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func setProductHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aprod := &model.AdminProduct{}
		err := json.NewDecoder(r.Body).Decode(aprod)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}
		prod := convert.ConvertAdminToDbProduct(*aprod)
		err = validate.ValidateProduct(&prod)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %s", BadRequestMessage, err.Error()), http.StatusBadRequest)
			log.Printf("setProductHandler: %s", err.Error())
			return
		}
		err = db.SetProduct(&prod)
		if err != nil {
			if validate.CatchOnDbError(err) {
				http.Error(w, fmt.Sprintf("Категории с  id = %d не существует", prod.Category_id), http.StatusBadRequest)
				return
			}
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		aprod.ProductId = prod.Product_id

		res, err := json.Marshal(aprod)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("setProductHandler: %s", err.Error())
			return
		}

		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func updateProductAdminHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		aprod := &model.AdminProduct{}
		err := json.NewDecoder(r.Body).Decode(aprod)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		params := convert.ConvertAprodToMap(*aprod)
		if err = validate.ValidateProductMap(params); err != nil {
			http.Error(w, fmt.Sprintf("%s: %s", BadRequestMessage, err.Error()), http.StatusBadRequest)
			log.Printf("updateProductAdminHandler: %s", err.Error())
			return
		}
		id := aprod.ProductId
		isUpdated, err := db.UpdateProductById(id, params)
		if err != nil {
			if validate.CatchOnDbError(err) {
				http.Error(w, fmt.Sprintf("Категории с таким id не существует"), http.StatusBadRequest)
				return
			}
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("updateProductAdminHandler: %s", err.Error())
			return
		}
		if isUpdated == false {
			http.Error(w, fmt.Sprintf("No product with id %d in base", id), http.StatusNotFound)
			return
		}
		resp := response.ResponseMessage{Message: fmt.Sprintf("Product with id %d was updated successfully", id)}
		res, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func deleteProductByIdHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			log.Printf("deleteProductByIdHandler: %v", err)
			return
		}
		id64 := int64(id)

		isDeleted, err := db.DeleteProductById(id64)

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("deleteProductByIdHandler: %s", err.Error())
			return
		}
		resp := response.ResponseMessage{}
		if isDeleted == false {
			resp.Message = fmt.Sprintf("Product with id %d was not found", id64)
		} else {
			resp.Message = fmt.Sprintf("Product with id %d was deleted successfully", id64)
		}
		res, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

//для categories
func getCategoriesListAdminHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		catList, err := db.GetCategoriesList()
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getCategoriesListAdminHandler: %s", err.Error())
			return
		}

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getCategoriesListAdminHandler: %s", err.Error())
			return
		}
		if catList == nil {
			http.Error(w, "No categories in base", http.StatusNotFound)
			return
		}
		acatList := make([]model.AdminCategory, 0, len(catList))
		for _, val := range catList {
			acatList = append(acatList, convert.ConvertToAdminCategory(val))
		}

		res, err := json.Marshal(acatList)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func getCategoryByIdAdminHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}

		id64 := int64(id)
		cat, err := db.GetCategoryById(id64)

		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("getCategoryByIdAdminHandler: %s", err.Error())
			return
		}
		if cat == nil {
			http.Error(w, fmt.Sprintf("No category with id %d in base", id64), http.StatusNotFound)
			return
		}

		acat := convert.ConvertToAdminCategory(*cat)
		res, err := json.Marshal(acat)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func setCategoryHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acat := &model.AdminCategory{}
		err := json.NewDecoder(r.Body).Decode(acat)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}
		cat := convert.ConvertAdminToDbCategory(*acat)
		err = validate.ValidateCategory(&cat)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %s", BadRequestMessage, err.Error()), http.StatusBadRequest)
			return
		}
		err = db.SetCategory(&cat)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("setCategoryHandler: %v", err.Error())
			return
		}
		acat.CategoryId = cat.Category_id
		res, err := json.Marshal(acat)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func updateCategoryHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		acat := &model.AdminCategory{}
		err := json.NewDecoder(r.Body).Decode(acat)
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}
		params := convert.ConvertAcatToMap(*acat)
		err = validate.ValidateCategoryMap(params)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s: %s", BadRequestMessage, err.Error()), http.StatusBadRequest)
			return
		}
		id := acat.CategoryId
		isUpdated, err := db.UpdateCategoryById(id, params)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("updateCategoryAdminHandler: %v", err.Error())
			return
		}
		if isUpdated == false {
			http.Error(w, fmt.Sprintf("No category with id %d in base", id), http.StatusNotFound)
			return
		}
		resp := response.ResponseMessage{}
		resp.Message = fmt.Sprintf("Category with id %d was updated successfully", acat.CategoryId)
		res, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}

func deleteCategoryByIdHandler(db model.IDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, BadRequestMessage, http.StatusBadRequest)
			return
		}
		id64 := int64(id)

		isDeleted, err := db.DeleteCategoryById(id64)

		if err != nil {
			if validate.CatchOnDbError(err) {
				http.Error(w, fmt.Sprintf("Категория не может быть удалена"), http.StatusForbidden)
				return
			}
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			log.Printf("deleteCategoryByIdHandler: %s", err.Error())
			return
		}
		resp := response.ResponseMessage{}
		if isDeleted == false {
			resp.Message = fmt.Sprintf("Category with id %d was not found", id64)
		} else {
			resp.Message = fmt.Sprintf("Category with id %d was deleted successfully", id64)
		}
		res, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, InternalErrorMessage, http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(res)
	}
}
