package app

import (
	"food-serve.com/internal/cache"
	"food-serve.com/internal/middleware"
	"food-serve.com/internal/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewApplication(database *gorm.DB, r *gin.Engine, rdb cache.Cache) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recovery())
	//health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "db": "connected", "redis": "connected"})
	})

	userStore := user.NewStore(database) //now userStore has db inside it
	userService := user.NewHandler(userStore, rdb)
	userService.RegisterRoutes(r)

	authenticated := r.Group("/")
	authenticated.Use(middleware.Authenticate())
	authenticated.GET("/me", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello", "userID": c.GetInt(middleware.UserIDKey), "email": c.GetString(middleware.EmailKey)})
	})

}
