package app

import (
	"food-serve.com/internal/middleware"
	"food-serve.com/internal/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewApplication(database *gorm.DB, r *gin.Engine) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	//health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "db": "connected", "redis": "connected"})
	})

	userStore := user.NewStore(database) //now userStore has db inside it
	userService := user.NewHandler(userStore)
	userService.RegisterRoutes(r)
}
