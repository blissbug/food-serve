package user

import (
	"fmt"

	"food-serve.com/pkg/types"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (userStore UserStore) CreateUser(user types.User) (uint, error) {
	res := userStore.db.Create(&user)

	if res.Error != nil {
		return 0, res.Error
	}

	return user.ID, nil
}

func (userStore UserStore) FindUserByEmail(email string) (user types.User, err error) {
	res := userStore.db.Find(&user, "email = ?", email)

	if res.RowsAffected == 0 {
		return user, fmt.Errorf("No such user exists, please register!")
	}

	if res.Error != nil {
		return user, res.Error
	}

	return user, nil
}

func (userStore UserStore) VerifyUser(user types.User) (bool, error) {
	res := userStore.db.Where("email = ?", user.Email).First(&user).Update("is_verified", true)
	if res.Error != nil {
		return false, res.Error
	}
	return true, nil
}
