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

func (userStore UserStore) CreateUser(user types.User) error {
	res := userStore.db.Find(&user, "email = ?", user.Email)

	if res.RowsAffected > 0 {
		return fmt.Errorf("A user with this email already exists")
	}

	res = userStore.db.Create(&user)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (userStore UserStore) FindUser(email string) (user types.User, err error) {
	res := userStore.db.Find(&user, "email = ?", email)

	if res.RowsAffected == 0 {
		return user, fmt.Errorf("No such user exists, please register!")
	}

	if res.Error != nil {
		return user, res.Error
	}

	return user, nil
}
