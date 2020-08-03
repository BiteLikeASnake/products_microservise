package model

type IDb interface {
	Close() error

	GetProductsList() ([]Product, error)
	GetProductById(id int64) (*Product, error)
	GetProductsByFilter(filter map[string]interface{}) ([]Product, error)
	SetProduct(newProduct *Product) error
	UpdateProductById(id int64, params map[string]interface{}) (isUpdated bool, err error)
	DeleteProductById(id int64) (isDeleted bool, err error)

	GetCategoriesList() ([]Category, error)
	GetCategoryById(id int64) (*Category, error)
	GetCategoriesByFilter(filter map[string]interface{}) ([]Category, error)
	SetCategory(newCategory *Category) error
	UpdateCategoryById(id int64, params map[string]interface{}) (isUpdated bool, err error)
	DeleteCategoryById(id int64) (isDeleted bool, err error)

	GetUserProductsList() ([]UserProduct, error)
	GetUserProductsListFilteredByCategoryId(id int64) ([]UserProduct, error)
	GetUserProductById(id int64) (*UserProduct, error)
	GetProductsByUserFilter(filter FilterUProduct) ([]UserProduct, error)
}
