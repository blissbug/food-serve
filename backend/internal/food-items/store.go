package food_items

import "gorm.io/gorm"

type foodItemsStore struct {
	foodItemsStore *gorm.DB
}

func NewStore(store *gorm.DB) foodItemsStore {
	return foodItemsStore{
		foodItemsStore: store,
	}
}

func (fStore foodItemsStore) CreateFoodItem() {

}
