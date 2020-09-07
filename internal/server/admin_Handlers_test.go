package server

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/call-me-snake/products_microservise/internal/model"
	mock_model "github.com/call-me-snake/products_microservise/internal/model/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	testCatId1             int64   = 1
	testCatId2             int64   = 2
	testProdName1                  = "test1"
	testProdName2                  = "test2"
	testCatName1                   = "testc1"
	testProdId1            int64   = 1
	testProdId2            int64   = 2
	testSupPrice1          float64 = 1000
	testSupPrice2          float64 = 500
	testSalePrice1         float64 = 1500
	testSalePrice2         float64 = 700
	testProdQuantity1      int64   = 10
	testProdQuantity2      int64   = 5
	testNoProductsQuantity int64   = 0
	testProdList                   = []model.Product{
		{
			Product_id:           testProdId1,
			Product_name:         testProdName1,
			Supply_product_price: testSupPrice1,
			Sale_product_price:   testSalePrice1,
			Category_id:          &testCatId1,
			Product_quantity:     testProdQuantity1,
			Product_active:       true,
		},
		{
			Product_id:           testProdId2,
			Product_name:         testProdName2,
			Supply_product_price: testSupPrice2,
			Sale_product_price:   testSalePrice2,
			Category_id:          &testCatId2,
			Product_quantity:     testProdQuantity2,
			Product_active:       false,
		},
	}

	testAProdList = []model.AdminProduct{
		{
			ProductId:          testProdId1,
			ProductName:        testProdName1,
			SupplyProductPrice: testSupPrice1,
			SaleProductPrice:   testSalePrice1,
			CategoryId:         testCatId1,
			ProductQuantity:    testProdQuantity1,
			ProductActive:      true,
		},
		{
			ProductId:          testProdId2,
			ProductName:        testProdName2,
			SupplyProductPrice: testSupPrice2,
			SaleProductPrice:   testSalePrice2,
			CategoryId:         testCatId2,
			ProductQuantity:    testProdQuantity2,
			ProductActive:      false,
		},
	}
	testACat = model.AdminCategory{
		CategoryId:     testCatId1,
		CategoryName:   testCatName1,
		CategoryActive: true,
	}
	setAProd = model.AdminProduct{
		ProductId:          0,
		ProductName:        testProdName1,
		SupplyProductPrice: testSupPrice1,
		SaleProductPrice:   testSalePrice1,
		CategoryId:         testCatId1,
		ProductQuantity:    testProdQuantity1,
		ProductActive:      true,
	}
	setProd = model.Product{
		Product_id:           0,
		Product_name:         testProdName1,
		Supply_product_price: testSupPrice1,
		Sale_product_price:   testSalePrice1,
		Category_id:          &testCatId1,
		Product_quantity:     testProdQuantity1,
		Product_active:       true,
	}
)

func TestGetProductsListAdminHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		GetProductsList().
		Return(testProdList, nil)

	req, err := http.NewRequest("GET", "/products/", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductsListAdminHandler(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	res, err := json.Marshal(testAProdList)
	assert.Equal(t, res, rr.Body.Bytes())
}

func TestGetEmptyProductsListAdminHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		GetProductsList().
		Return(nil, nil)
	req, err := http.NewRequest("GET", "/products/", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductsListAdminHandler(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetProductByIdAdminHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		GetProductById(testProdId1).
		Return(&testProdList[0], nil)

	router := mux.NewRouter()
	router.HandleFunc("/products/{id:[0-9]+}", getProductByIdAdminHandler(mockdb)).Methods("GET")
	req, err := http.NewRequest("GET", "/products/1", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	res, err := json.Marshal(testAProdList[0])
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, res, rr.Body.Bytes())
}

func TestGetNonPresentProductByIdAdminHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		GetProductById(int64(1000)).
		Return(nil, nil)

	router := mux.NewRouter()
	router.HandleFunc("/products/{id:[0-9]+}", getProductByIdAdminHandler(mockdb)).Methods("GET")
	req, err := http.NewRequest("GET", "/products/1000", nil)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetProductsInCategoryAdminHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		GetProductsByFilter(map[string]interface{}{"category_id": testCatId1}).
		Return(testProdList[:1], nil)

	inputJson, err := json.Marshal(testACat)
	if err != nil {
		log.Fatal(err)
	}
	rbody := bytes.NewReader(inputJson)
	req, err := http.NewRequest("POST", "/categories/products", rbody)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductsInCategoryAdminHandler(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	res, err := json.Marshal(testAProdList[:1])
	assert.Equal(t, res, rr.Body.Bytes())
}

func TestGetProductsInEmptyCategoryAdminHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		GetProductsByFilter(map[string]interface{}{"category_id": testCatId1}).
		Return(nil, nil)

	inputJson, err := json.Marshal(testACat)
	if err != nil {
		log.Fatal(err)
	}
	rbody := bytes.NewReader(inputJson)
	req, err := http.NewRequest("POST", "/categories/products", rbody)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getProductsInCategoryAdminHandler(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestSetProductHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockdb := mock_model.NewMockIDb(ctrl)
	mockdb.
		EXPECT().
		SetProduct(&setProd).
		Return(nil)

	inputJson, _ := json.Marshal(&setAProd)
	rbody := bytes.NewReader(inputJson)
	req, err := http.NewRequest("POST", "/products", rbody)
	if err != nil {
		log.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(setProductHandler(mockdb))
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
