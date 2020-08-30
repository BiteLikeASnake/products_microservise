package db

import (
	"fmt"

	"github.com/call-me-snake/products_microservise/internal/model"
)

//GetAllCategories ...
func (db Db) GetCategoriesList() ([]model.Category, error) {
	categoriesList := make([]model.Category, 0)

	query := db.Database.Find(&categoriesList)

	if query.Error != nil {
		return nil, fmt.Errorf("db.GetCategoriesList: %v", query.Error)
	}

	if query.RowsAffected == 0 {
		return nil, nil
	}
	return categoriesList, nil
}

//GetCategoryById ...
func (db Db) GetCategoryById(id int64) (*model.Category, error) {
	resultCategory := &model.Category{}

	query := db.Database.First(resultCategory, id)

	if query.Error != nil {
		if query.Error.Error() == "record not found" {
			return nil, nil
		}

		return nil, fmt.Errorf("db.GetCategoryById: %v", query.Error)
	}

	return resultCategory, nil
}

//GetCategoryByFilter ...
func (db Db) GetCategoriesByFilter(filter map[string]interface{}) ([]model.Category, error) {
	filteredCats := make([]model.Category, 0)
	query := db.Database.Where(filter).Find(&filteredCats)
	if query.Error != nil {
		return nil, fmt.Errorf("GetCategoryByFilter: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}
	return filteredCats, nil
}

//SetCategory ...
func (db Db) SetCategory(newCategory *model.Category) error {
	if isNew := db.Database.NewRecord(newCategory); !isNew {
		return fmt.Errorf("db.SetCategory: параметры на входе имеют ненулевой pk %+v", newCategory)
	}

	set := db.Database.Create(newCategory)

	if set.Error != nil {
		return fmt.Errorf("db.SetCategory: %v", set.Error)
	}
	return nil
}

//UpdateCategory ...
func (db Db) UpdateCategoryById(id int64, params map[string]interface{}) (isUpdated bool, err error) {
	update := db.Database.Model(model.Category{}).Where(id).Updates(params)

	if update.Error != nil {
		isUpdated = false
		err = fmt.Errorf("db.UpdateCategory: %v", update.Error)
		return
	}

	if update.RowsAffected == 0 {
		isUpdated = false
		err = nil
		return
	}
	isUpdated = true
	err = nil
	return
}

//DeleteCategory ...
func (db Db) DeleteCategoryById(id int64) (isDeleted bool, err error) {
	delete := db.Database.Where("category_id =  ?", id).Delete(model.Category{})

	if delete.Error != nil {
		isDeleted = false
		err = fmt.Errorf("db.DeleteCategory: %v", delete.Error)
		return
	}

	if delete.RowsAffected == 0 {
		isDeleted = false
		err = nil
		return
	}
	isDeleted = true
	err = nil
	return
}
