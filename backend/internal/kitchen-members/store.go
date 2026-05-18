package kitchenmembers

import (
	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type KitchenMembersStore struct {
	store *gorm.DB
}

func NewStore(store *gorm.DB) KitchenMembersStore {
	return KitchenMembersStore{
		store: store,
	}
}

func (kStore KitchenMembersStore) CreateKitchenMember(userId uint, kitchenId uint, role string) error {
	member := types.KitchenMember{
		UserID:    userId,
		KitchenID: kitchenId,
		Role:      role,
	}
	res := kStore.store.Create(member)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (kStore KitchenMembersStore) GetKitchenMembersByKitchenId(kitchenId uint) ([]types.KitchenMember, error) {
	var members []types.KitchenMember
	res := kStore.store.Where("kitchen_id = ?", kitchenId).Find(&members)
	if res.Error != nil {
		return nil, res.Error
	}
	return members, nil
}

func (kStore KitchenMembersStore) GetKitchensSubscribedByAUser(userId uint) ([]types.KitchenMember, error) {
	var members []types.KitchenMember
	res := kStore.store.Where("user_id = ?", userId).Find(&members)
	if res.Error != nil {
		return nil, res.Error
	}
	return members, nil
}

func (kStore KitchenMembersStore) IsUserAdminOfThisKitchen(userId uint, kitchenId uint) error {
	res := kStore.store.Where("user_id = ? AND kitchen_id = ? AND role = ?", userId, kitchenId, "admin").First(&types.KitchenMember{})
	if res.Error != nil {
		return res.Error
	}
	return nil
}
