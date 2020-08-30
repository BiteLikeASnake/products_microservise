package db

import (
	"fmt"
	"log"
	"os"

	"github.com/call-me-snake/products_microservise/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	product_example1 = model.Product{
		Product_name:   "Test_product1",
		Product_active: true,
	}

	product_example2 = model.Product{
		Product_name:   "Test_product2",
		Product_active: false,
	}

	updateProductParams = map[string]interface{}{
		"Product_name":          "Updated",
		"Default_product_price": 100.01,
		"Sale_product_price":    200.01,
		"Product_active":        false,
	}

	current_prod_id int64 //for getting product by id
)

type ProductTestSuite struct {
	suite.Suite
	Db *Db
}

//SetupSuite implements interface SetupAllSuite. method, which will run before the tests in the suite are run.
func (mySuite *ProductTestSuite) SetupSuite() {
	testproductDbAdress := os.Getenv("TEST_DB")
	log.Println("SetupSuite began")
	var err error

	mySuite.Db, err = New(testproductDbAdress)
	//mySuite.DB.DB.LogMode(true)

	if err != nil {
		log.Fatal(fmt.Errorf("SetupSuite: %v", err))
	}

	deleted := mySuite.Db.Database.Delete(&model.Product{})
	log.Printf("Deleted rows: %d", deleted.RowsAffected)
	log.Println("SetupSuite ended")
}

//SetupTest implements interface SetupTestSuite. method, which will run before each test in the suite.
func (mySuite *ProductTestSuite) SetupTest() {
	log.Println("SetupTest began")
	prodToWorkWith := product_example1
	err := mySuite.Db.SetProduct(&prodToWorkWith)
	if err != nil {
		log.Fatal(fmt.Errorf("SetupTest: %s. Problem with inserting %v", err, product_example1))
	}
	current_prod_id = prodToWorkWith.Product_id

	log.Println("SetupTest ended")
}

//TearDownSuite implements interface TearDownAllSuite. method, which will run after all the tests in the suite have been run.
func (mySuite *ProductTestSuite) TearDownSuite() {
	log.Println("TearDownSuite begun")

	deleted := mySuite.Db.Database.Delete(&model.Product{})
	log.Printf("Deleted rows: %d", deleted.RowsAffected)

	mySuite.Db.Close()
	log.Println("TearDownSuite ended")
}

//TearDownTest implements interface TearDownTestSuite. Method, which will run after each test in the suite.
func (mySuite *ProductTestSuite) TearDownTest() {
	log.Println("TearDownTest begun")

	//mySuite.DB.DB = mySuite.DB.DB.Delete(&model.Employee{})
	deleted := mySuite.Db.Database.Delete(&model.Product{})
	log.Printf("Deleted rows: %d", deleted.RowsAffected)

	log.Println("TearDownTest ended")
}

func (mySuite *ProductTestSuite) TestGetprodegoriesList() {

	prod, err := mySuite.Db.GetProductsList()
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), prod)
}

func (mySuite *ProductTestSuite) TestGetEmptyProductsList() {
	mySuite.TearDownTest()

	prod, err := mySuite.Db.GetProductsList()
	assert.NoError(mySuite.T(), err)
	assert.Nil(mySuite.T(), prod)
}

func (mySuite *ProductTestSuite) TestGetProductById() {
	prod, err := mySuite.Db.GetProductById(current_prod_id)
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), prod)
	assert.Equal(mySuite.T(), product_example1.Product_name, prod.Product_name)
	assert.Equal(mySuite.T(), product_example1.Product_active, prod.Product_active)
}

func (mySuite *ProductTestSuite) TestGetProductByNonPresentId() {
	prod, err := mySuite.Db.GetProductById(current_prod_id + 100)
	assert.NoError(mySuite.T(), err)
	assert.Nil(mySuite.T(), prod)
}

func (mySuite *ProductTestSuite) TestSetProduct() {
	prodToWorkWith := product_example2
	err := mySuite.Db.SetProduct(&prodToWorkWith)
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), prodToWorkWith)
	assert.Equal(mySuite.T(), product_example2.Product_name, prodToWorkWith.Product_name)
	assert.Equal(mySuite.T(), product_example2.Product_active, prodToWorkWith.Product_active)
	assert.Equal(mySuite.T(), current_prod_id+1, prodToWorkWith.Product_id)
}

func (mySuite *ProductTestSuite) TestSetProductWithPresentId() {
	prodToWorkWith := product_example2
	prodToWorkWith.Product_id = current_prod_id
	err := mySuite.Db.SetProduct(&prodToWorkWith)
	assert.Error(mySuite.T(), err)
	assert.Nil(mySuite.T(), prodToWorkWith)
}

func (mySuite *ProductTestSuite) TestUpdateproductById() {
	isUpdated, err := mySuite.Db.UpdateProductById(current_prod_id, updateProductParams)
	assert.NoError(mySuite.T(), err)
	assert.True(mySuite.T(), isUpdated)
	prod, err := mySuite.Db.GetProductById(current_prod_id)
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), prod)
	assert.Equal(mySuite.T(), updateProductParams["Product_name"], prod.Product_name)
	assert.Equal(mySuite.T(), updateProductParams["Default_product_price"], prod.Supply_product_price)
	assert.Equal(mySuite.T(), updateProductParams["Sale_product_price"], prod.Sale_product_price)
	assert.Equal(mySuite.T(), updateProductParams["Product_active"], prod.Product_active)
}

func (mySuite *ProductTestSuite) TestUpdateproductByNonPresentId() {
	isUpdated, err := mySuite.Db.UpdateProductById(current_prod_id+100, updateProductParams)
	assert.NoError(mySuite.T(), err)
	assert.False(mySuite.T(), isUpdated)
}

func (mySuite *ProductTestSuite) TestDeleteProduct() {
	isDeleted, err := mySuite.Db.DeleteProductById(current_prod_id)
	assert.NoError(mySuite.T(), err)
	assert.True(mySuite.T(), isDeleted)
	prod, err := mySuite.Db.GetProductById(current_prod_id)
	assert.NoError(mySuite.T(), err)
	assert.Nil(mySuite.T(), prod)
}

func (mySuite *ProductTestSuite) TestDeleteProductByNonPresentId() {
	isDeleted, err := mySuite.Db.DeleteProductById(current_prod_id + 100)
	assert.NoError(mySuite.T(), err)
	assert.False(mySuite.T(), isDeleted)

}
