package kitchen

import (
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type KitchenStore struct {
	store *gorm.DB
}

func NewStore(store *gorm.DB) KitchenStore {
	return KitchenStore{
		store: store,
	}
}

func (kStore KitchenStore) CreateKitchen(name string, slug string) (uint, error) {
	kitchen := types.Kitchen{
		Name: name,
		Slug: slug,
	}

	res := kStore.store.Create(&kitchen)

	if res.Error != nil {
		return 0, res.Error
	}

	return kitchen.ID, nil
}
