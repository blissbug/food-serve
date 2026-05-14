package app

import (
	"food-serve.com/internal/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewApplication(database *gorm.DB, r *gin.Engine) {
	userStore := user.NewStore(database) //now userStore has db inside it
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(r)
}
