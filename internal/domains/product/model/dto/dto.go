package dto

type Product struct {
	ID          int     `json:"id" gorm:"column:id;primaryKey;autoIncrement;"`
	Name        string  `json:"name" gorm:"column:name;"`
	Price       float64 `json:"price" gorm:"column:price;"`
	Description string  `json:"description" gorm:"column:description;"`
	Category    string  `json:"category" gorm:"column:category;"`
	Status      string  `json:"status" gorm:"column:status;"`
	Image       *string `json:"image,omitempty"`
}

func (Product) TableName() string {
	return "product"
}

type Products []Product
