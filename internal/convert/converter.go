package convert

import "github.com/call-me-snake/products_microservise/internal/model"

//ConvertToAdminProduct
//converts product.Product to admin.AdminProduct
func ConvertToAdminProduct(prod model.Product) model.AdminProduct {
	aprod := model.AdminProduct{}
	aprod.ProductId = prod.Product_id
	aprod.ProductName = prod.Product_name
	aprod.SupplyProductPrice = prod.Supply_product_price
	aprod.SaleProductPrice = prod.Sale_product_price
	if prod.Category_id == nil {
		aprod.CategoryId = 0
	} else {
		aprod.CategoryId = *prod.Category_id
	}
	aprod.ProductQuantity = prod.Product_quantity
	aprod.ProductActive = prod.Product_active
	return aprod
}

//ConvertAdminToDbProduct
//convert admin.AdminProduct to product.Product
func ConvertAdminToDbProduct(aprod model.AdminProduct) model.Product {
	prod := model.Product{}
	prod.Product_id = aprod.ProductId
	prod.Product_name = aprod.ProductName
	prod.Supply_product_price = aprod.SupplyProductPrice
	prod.Sale_product_price = aprod.SaleProductPrice
	if aprod.CategoryId == 0 {
		prod.Category_id = nil
	} else {
		prod.Category_id = &aprod.CategoryId
	}
	prod.Product_quantity = aprod.ProductQuantity
	prod.Product_active = aprod.ProductActive
	return prod

}

//ConvertAprodToMap
//convert admin.AdminProduct to map[string]interface{}
func ConvertAprodToMap(aprod model.AdminProduct) map[string]interface{} {
	params := make(map[string]interface{})
	params["product_name"] = aprod.ProductName
	params["supply_product_price"] = aprod.SupplyProductPrice
	params["sale_product_price"] = aprod.SaleProductPrice
	if aprod.CategoryId != 0 {
		params["category_id"] = aprod.CategoryId
	}
	params["product_quantity"] = aprod.ProductQuantity
	params["product_active"] = aprod.ProductActive
	return params
}

//ConvertToAdminCategory
//convert category.Category to admin.AdminCategory
func ConvertToAdminCategory(cat model.Category) model.AdminCategory {
	acat := model.AdminCategory{}
	acat.CategoryId = cat.Category_id
	acat.CategoryName = cat.Category_name
	acat.CategoryActive = cat.Category_active
	return acat
}

//ConvertAdminToDbCategory
//convert admin.AdminCategory to category.Category
func ConvertAdminToDbCategory(acat model.AdminCategory) model.Category {
	cat := model.Category{}
	cat.Category_id = acat.CategoryId
	cat.Category_name = acat.CategoryName
	cat.Category_active = acat.CategoryActive
	return cat
}

//ConvertAcatToMap
//convert admin.AdminCategory to map[string]interface{}
func ConvertAcatToMap(acat model.AdminCategory) map[string]interface{} {
	params := make(map[string]interface{})
	params["category_name"] = acat.CategoryName
	params["category_active"] = acat.CategoryActive
	return params
}

//ConvertToJsonProduct
func ConvertToJsonProduct(prod model.UserProduct) model.JsonProduct {
	uprod := model.JsonProduct{}
	uprod.ProductId = prod.Product_id
	uprod.ProductName = prod.Product_name
	uprod.ProductPrice = prod.Sale_product_price
	uprod.ProductQuantity = prod.Product_quantity
	uprod.CategoryId = prod.Category_id
	uprod.CategoryName = prod.Category_name
	return uprod
}

//ConvertToUserCategory
func ConvertToUserCategory(cat model.Category) model.UserCategory {
	ucat := model.UserCategory{}
	ucat.CategoryId = cat.Category_id
	ucat.CategoryName = cat.Category_name
	return ucat
}
