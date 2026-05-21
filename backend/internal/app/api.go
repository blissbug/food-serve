package app

import (
	"food-serve.com/internal/cache"
	foodItemImages "food-serve.com/internal/food-item-images"
	food_items "food-serve.com/internal/food-items"
	"food-serve.com/internal/kitchen"
	kitchenmembers "food-serve.com/internal/kitchen-members"
	"food-serve.com/internal/menu"
	menuItems "food-serve.com/internal/menu-items"
	"food-serve.com/internal/middleware"
	"food-serve.com/internal/user"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewApplication(database *gorm.DB, r *gin.Engine, rdb cache.Cache, cld *cloudinary.Cloudinary) {
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

	foodItemImagesStore := foodItemImages.NewStore(database)
	foodItemsStore := food_items.NewStore(database)
	foodItemService := food_items.NewService(database, foodItemsStore, foodItemImagesStore)
	foodItemHandler := food_items.NewHandler(foodItemService, cld)
	foodItemHandler.RegisterRoutes(authenticated)

	kitchenStore := kitchen.NewStore(database)
	kitchenService := kitchen.NewHandler(kitchenStore, kitchenMembersStore, userStore)
	kitchenService.RegisterRoutes(authenticated)

	menuStore := menu.NewStore(database)
	menuItemsStore := menuItems.NewStore(database)
	menuService := menu.NewService(menuStore, menuItemsStore, rdb)
	menuHandler := menu.NewHandler(menuService)
	menuHandler.RegisterRoutes(authenticated)

}
