package validate

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/call-me-snake/products_microservise/internal/model"
)

func ValidateCategoryMap(params map[string]interface{}) error {
	if _, ok := params["Category_id"]; ok {
		return fmt.Errorf("Ненулевой id")
	}
	if cat_name, ok := params["Category_name"]; ok {
		if utf8.RuneCountInString(cat_name.(string)) > 25 {
			return fmt.Errorf("Category_name length > 25")
		}
	}
	return nil
}

func ValidateProductMap(params map[string]interface{}) error {
	if _, ok := params["product_id"]; ok {
		return fmt.Errorf("ValidateProductMap: Ненулевой id")
	}
	if prod_name, ok := params["product_name"]; ok {
		if utf8.RuneCountInString(prod_name.(string)) > 25 {
			return fmt.Errorf("Product_name length > 25")
		}
	}
	if def_price, ok := params["supply_product_price"]; ok {
		if def_price.(float64) < 0 {
			return fmt.Errorf("Supply_product_price<0")
		}
	}
	if price, ok := params["sale_product_price"]; ok {
		if price.(float64) < 0 {
			return fmt.Errorf("Sale_product_price<0")
		}
	}
	if cat_id, ok := params["category_id "]; ok {
		if cat_id.(int64) <= 0 {
			return fmt.Errorf("Category_id<=0")
		}
	}
	if quantity, ok := params["product_quantity"]; ok {
		if quantity.(int64) < 0 {
			return fmt.Errorf("Product_quantity<0")
		}
	}
	return nil
}

func ValidateCategory(cat *model.Category) error {
	if cat.Category_id != 0 {
		return fmt.Errorf("Ненулевой Category_id")
	}
	if cat.Category_name == "" {
		return fmt.Errorf("Category_name не указан")
	}
	if utf8.RuneCountInString(cat.Category_name) > 25 {
		return fmt.Errorf("Category_name length > 25")
	}
	return nil
}

func ValidateProduct(prod *model.Product) error {
	if prod.Product_id != 0 {
		return fmt.Errorf("Ненулевой Product_id")
	}
	if prod.Product_name == "" {
		return fmt.Errorf("Product_name не указаан")
	}
	if utf8.RuneCountInString(prod.Product_name) > 25 {
		return fmt.Errorf("Product_name length > 25")
	}
	if prod.Supply_product_price <= 0 {
		return fmt.Errorf("Base_product_price<=0")
	}
	if prod.Sale_product_price <= 0 {
		return fmt.Errorf("Sale_product_price<=0")
	}
	if prod.Product_quantity < 0 {
		return fmt.Errorf("Product_quantity<0")
	}
	return nil
}

//CatchOnDbError проверяет ошибку базы на category_id_fk constraint
func CatchOnDbError(err error) bool {
	if strings.Contains(err.Error(), "category_id_fk") {
		return true
	}
	return false
}

//category_id_fk
