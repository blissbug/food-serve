package food_items

import "github.com/gin-gonic/gin"

func (h Handler) RegisterRoutes(r *gin.RouterGroup) {
	foodItemsGroup := r.Group("/food/items")
	foodItemsGroup.GET("/")
	foodItemsGroup.GET("/:id")
	foodItemsGroup.POST("/create")
	foodItemsGroup.PUT("/update/:id")
	foodItemsGroup.DELETE("/delete/:id")
}
