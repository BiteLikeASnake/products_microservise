package db

import (
	"fmt"

	"github.com/BiteLikeASnake/products_microservise/internal/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//GetProductsList ...
func (db Db) GetProductsList() ([]model.Product, error) {
	var productsList []model.Product = make([]model.Product, 0)

	query := db.Database.Find(&productsList)

	if query.Error != nil {
		return nil, fmt.Errorf("db.GetProductsList: %v", query.Error)
	}

	if query.RowsAffected == 0 {
		return nil, nil
	}
	return productsList, nil
}

//GetUserProductsList ...
func (db Db) GetUserProductsList() ([]model.UserProduct, error) {
	uproductsList := make([]model.UserProduct, 0)

	query := db.Database.Table("products").
		Select("products.product_id, products.product_name, products.product_name, products.supply_product_price, products.sale_product_price, categories.category_id, categories.category_name, categories.category_active, products.product_quantity, products.product_active").
		Joins("JOIN categories ON categories.category_id = products.category_id ").
		Where("categories.category_active = ? AND products.product_active = ?", true, true).
		Find(&uproductsList)

	if query.Error != nil {
		return nil, fmt.Errorf("db.GetUserProductsList: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}
	return uproductsList, nil
}

func (db Db) GetUserProductsListFilteredByCategoryId(id int64) ([]model.UserProduct, error) {
	uproductsList := make([]model.UserProduct, 0)

	query := db.Database.Table("products").
		Select("products.product_id, products.product_name, products.product_name, products.supply_product_price, products.sale_product_price, categories.category_id, categories.category_name, categories.category_active, products.product_quantity, products.product_active").
		Joins("JOIN categories ON categories.category_id = products.category_id ").
		Where("categories.category_active = ? AND products.product_active = ? AND categories.category_id = ?", true, true, id).
		Find(&uproductsList)

	if query.Error != nil {
		return nil, fmt.Errorf("db.GetUserProductsList: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}
	return uproductsList, nil
}

func (db Db) GetUserProductById(id int64) (*model.UserProduct, error) {
	uprod := &model.UserProduct{}
	query := db.Database.Table("products").
		Select("products.product_id, products.product_name, products.product_name, products.supply_product_price, products.sale_product_price, categories.category_id, categories.category_name, categories.category_active, products.product_quantity, products.product_active").
		Joins("JOIN categories ON categories.category_id = products.category_id ").
		Where("categories.category_active = ? AND products.product_active = ? AND products.product_id = ?", true, true, id).
		First(uprod)

	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("GetUserProductById: %v", query.Error)
	}
	return uprod, nil
}

//GetProductById ...
func (db Db) GetProductById(id int64) (*model.Product, error) {
	resultProduct := &model.Product{}

	query := db.Database.First(resultProduct, id)

	if query.Error != nil {
		if query.Error.Error() == "record not found" {
			return nil, nil
		}

		return nil, fmt.Errorf("db.GetProductById: %v", query.Error)
	}

	return resultProduct, nil
}

//GetProductByFilter ...
//В filter map необходимые ключи должны быть с маленькой буквы (а в Update почему-то можно и с большой)
//TODO проверить, можно ли обойтись без rowsaffected
func (db Db) GetProductsByFilter(filter map[string]interface{}) ([]model.Product, error) {
	filteredProds := make([]model.Product, 0)
	query := db.Database.Where(filter).Find(&filteredProds)
	if query.Error != nil {
		return nil, fmt.Errorf("GetProductByFilter: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}
	return filteredProds, nil
}

//SetProduct ...
func (db Db) SetProduct(newProduct *model.Product) error {
	if isNew := db.Database.NewRecord(newProduct); !isNew {
		return fmt.Errorf("db.SetProduct: параметры на входе имеют ненулевой pk %+v", newProduct)
	}

	set := db.Database.Create(newProduct)

	if set.Error != nil {
		return fmt.Errorf("db.SetProduct: %v", set.Error)
	}
	return nil
}

//UpdateProduct ...
func (db Db) UpdateProductById(id int64, params map[string]interface{}) (isUpdated bool, err error) {
	update := db.Database.Model(model.Product{}).Where(id).Updates(params)

	if update.Error != nil {
		isUpdated = false
		err = fmt.Errorf("db.UpdateProduct: %v", update.Error)
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

//DeleteProduct ...
func (db Db) DeleteProductById(id int64) (isDeleted bool, err error) {
	delete := db.Database.Where("Product_id =  ?", id).Delete(model.Product{})

	if delete.Error != nil {
		isDeleted = false
		err = fmt.Errorf("db.DeleteProduct: %v", delete.Error)
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

func (db Db) GetProductsByUserFilter(filter model.FilterUProduct) ([]model.UserProduct, error) {
	var productsList = make([]model.UserProduct, 0)
	query := db.Database.Table("products").
		Select("products.product_id, products.product_name, products.product_name, products.supply_product_price, products.sale_product_price, categories.category_id, categories.category_name, categories.category_active, products.product_quantity, products.product_active").
		Joins("JOIN categories ON categories.category_id = products.category_id ").
		Where("categories.category_active = ? AND products.product_active = ?", true, true)
	if filter.ProductName != nil {
		query = query.Where("products.product_name LIKE ?", fmt.Sprintf("%%%s%%", *filter.ProductName))
	}
	if filter.MaxProductPrice != nil {
		query = query.Where("products.sale_product_price <= ?", filter.MaxProductPrice)
	}
	if filter.MinProductPrice != nil {
		query = query.Where("products.sale_product_price >= ?", filter.MinProductPrice)
	}
	if filter.CategoryId != nil {
		query = query.Where("categories.category_id = ?", filter.CategoryId)
	}
	query = query.Find(&productsList)
	if query.Error != nil {
		return nil, fmt.Errorf("db.GetProductsByUserFilter: %v", query.Error)
	}
	if query.RowsAffected == 0 {
		return nil, nil
	}
	return productsList, nil
}
