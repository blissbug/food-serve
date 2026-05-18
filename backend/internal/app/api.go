package app

import (
	"food-serve.com/internal/cache"
	food_items "food-serve.com/internal/food-items"
	"food-serve.com/internal/kitchen"
	kitchenmembers "food-serve.com/internal/kitchen-members"
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

	kitchenMembersStore := kitchenmembers.NewStore(database)

	userStore := user.NewStore(database) //now userStore has db inside it
	userService := user.NewHandler(userStore, rdb, &kitchenMembersStore)
	userService.RegisterRoutes(r)

	authenticated := r.Group("/")
	authenticated.Use(middleware.Authenticate())
	authenticated.GET("/me", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hello", "userID": c.GetUint(middleware.UserIDKey), "email": c.GetString(middleware.EmailKey)})
	})

	foodItemsStore := food_items.NewStore(database)
	foodItemService := food_items.NewHandler(foodItemsStore)
	foodItemService.RegisterRoutes(authenticated)

	kitchenStore := kitchen.NewStore(database)
	kitchenService := kitchen.NewHandler(kitchenStore, kitchenMembersStore, userStore)
	kitchenService.RegisterRoutes(authenticated)

}
