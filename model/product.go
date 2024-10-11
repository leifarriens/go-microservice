package model

import (
	"time"
)

type Product struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Available bool      `json:"available"`
}

type ProductDto struct {
	Name      string  `json:"name" validate:"required"`
	Price     float64 `json:"price" validate:"min=0,required"`
	Available bool    `json:"available" validate:"required"`
}
