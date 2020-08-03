package model

type Category struct {
	Category_id     int64 `gorm:"primary_key"`
	Category_name   string
	Category_active bool
}

//TODO	функции проверки, определиться, нужны ли указатели на поля
