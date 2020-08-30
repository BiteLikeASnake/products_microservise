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
	category_example1 = model.Category{
		Category_name:   "Test_category1",
		Category_active: true,
	}

	category_example2 = model.Category{
		Category_name:   "Test_category2",
		Category_active: false,
	}

	updateCategoryParams = map[string]interface{}{
		"Category_name":   "Updated",
		"Category_active": false,
	}

	current_cat_id int64 //for getting category by id
)

type CategoryTestSuite struct {
	suite.Suite
	Db *Db
}

//SetupSuite implements interface SetupAllSuite. method, which will run before the tests in the suite are run.
func (mySuite *CategoryTestSuite) SetupSuite() {
	testCategoryDbAdress := os.Getenv("TEST_DB")
	log.Println("SetupSuite began")
	var err error

	mySuite.Db, err = New(testCategoryDbAdress)
	//mySuite.DB.DB.LogMode(true)

	if err != nil {
		log.Fatal(fmt.Errorf("SetupSuite: %v", err))
	}

	deleted := mySuite.Db.Database.Delete(&model.Category{})
	log.Printf("Deleted rows: %d", deleted.RowsAffected)
	log.Println("SetupSuite ended")
}

//SetupTest implements interface SetupTestSuite. method, which will run before each test in the suite.
func (mySuite *CategoryTestSuite) SetupTest() {
	log.Println("SetupTest began")
	catToWorkWith := category_example1
	err := mySuite.Db.SetCategory(&catToWorkWith)
	if err != nil {
		log.Fatal(fmt.Errorf("SetupTest: %s. Problem with inserting %v", err, category_example1))
	}
	current_cat_id = catToWorkWith.Category_id

	log.Println("SetupTest ended")
}

//TearDownSuite implements interface TearDownAllSuite. method, which will run after all the tests in the suite have been run.
func (mySuite *CategoryTestSuite) TearDownSuite() {
	log.Println("TearDownSuite begun")

	deleted := mySuite.Db.Database.Delete(&model.Category{})
	log.Printf("Deleted rows: %d", deleted.RowsAffected)

	mySuite.Db.Close()
	log.Println("TearDownSuite ended")
}

//TearDownTest implements interface TearDownTestSuite. Method, which will run after each test in the suite.
func (mySuite *CategoryTestSuite) TearDownTest() {
	log.Println("TearDownTest begun")

	//mySuite.DB.DB = mySuite.DB.DB.Delete(&model.Employee{})
	deleted := mySuite.Db.Database.Delete(&model.Category{})
	log.Printf("Deleted rows: %d", deleted.RowsAffected)

	log.Println("TearDownTest ended")
}

func (mySuite *CategoryTestSuite) TestGetCategoriesList() {

	cat, err := mySuite.Db.GetCategoriesList()
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), cat)
}

func (mySuite *CategoryTestSuite) TestGetEmptyCategoriesList() {
	mySuite.TearDownTest()

	cat, err := mySuite.Db.GetCategoriesList()
	assert.NoError(mySuite.T(), err)
	assert.Nil(mySuite.T(), cat)
}

func (mySuite *CategoryTestSuite) TestGetCategoryById() {
	cat, err := mySuite.Db.GetCategoryById(current_cat_id)
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), cat)
	assert.Equal(mySuite.T(), category_example1.Category_name, cat.Category_name)
	assert.Equal(mySuite.T(), category_example1.Category_active, cat.Category_active)
}

func (mySuite *CategoryTestSuite) TestGetCategoryByNonPresentId() {
	cat, err := mySuite.Db.GetCategoryById(current_cat_id + 100)
	assert.NoError(mySuite.T(), err)
	assert.Nil(mySuite.T(), cat)
}

func (mySuite *CategoryTestSuite) TestSetCategory() {
	catToWorkWith := category_example2
	err := mySuite.Db.SetCategory(&catToWorkWith)
	assert.NoError(mySuite.T(), err)
	assert.Equal(mySuite.T(), category_example2.Category_name, catToWorkWith.Category_name)
	assert.Equal(mySuite.T(), category_example2.Category_active, catToWorkWith.Category_active)
	assert.Equal(mySuite.T(), current_cat_id+1, catToWorkWith.Category_id)
}

func (mySuite *CategoryTestSuite) TestSetCategoryWithPresentId() {
	catToWorkWith := category_example2
	catToWorkWith.Category_id = current_cat_id
	err := mySuite.Db.SetCategory(&catToWorkWith)
	assert.Error(mySuite.T(), err)
}

func (mySuite *CategoryTestSuite) TestUpdateCategoryById() {
	isUpdated, err := mySuite.Db.UpdateCategoryById(current_cat_id, updateCategoryParams)
	assert.NoError(mySuite.T(), err)
	assert.True(mySuite.T(), isUpdated)
	cat, err := mySuite.Db.GetCategoryById(current_cat_id)
	assert.NoError(mySuite.T(), err)
	assert.NotNil(mySuite.T(), cat)
	assert.Equal(mySuite.T(), updateCategoryParams["Category_name"], cat.Category_name)
	assert.Equal(mySuite.T(), updateCategoryParams["Category_active"], cat.Category_active)
}

func (mySuite *CategoryTestSuite) TestUpdateCategoryByNonPresentId() {
	isUpdated, err := mySuite.Db.UpdateCategoryById(current_cat_id+100, updateCategoryParams)
	assert.NoError(mySuite.T(), err)
	assert.False(mySuite.T(), isUpdated)
}

func (mySuite *CategoryTestSuite) TestDeleteCategory() {
	isDeleted, err := mySuite.Db.DeleteCategoryById(current_cat_id)
	assert.NoError(mySuite.T(), err)
	assert.True(mySuite.T(), isDeleted)
	cat, err := mySuite.Db.GetCategoryById(current_cat_id)
	assert.NoError(mySuite.T(), err)
	assert.Nil(mySuite.T(), cat)
}

func (mySuite *CategoryTestSuite) TestDeleteCategoryByNonPresentId() {
	isDeleted, err := mySuite.Db.DeleteCategoryById(current_cat_id + 100)
	assert.NoError(mySuite.T(), err)
	assert.False(mySuite.T(), isDeleted)

}

//TestRun runs all tests
