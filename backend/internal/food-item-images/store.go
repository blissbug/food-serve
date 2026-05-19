package foodItemImages

import (
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type Store struct {
	Store *gorm.DB
}

func NewStore(store *gorm.DB) Store {
	return Store{
		Store: store,
	}
}

func (store *Store) AddImageForItem(tx *gorm.DB, foodItemImage types.FoodItemImage) error {
	resp := tx.Create(&foodItemImage)

	if resp.Error != nil {
		return resp.Error
	}
	return nil
}
