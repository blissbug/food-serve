package food_items

import "github.com/gin-gonic/gin"

func (h Handler) RegisterRoutes(r *gin.RouterGroup) {
	foodItemsGroup := r.Group("/food/items")
	foodItemsGroup.GET("/")                                       //get all the food items, paginated
	foodItemsGroup.GET("/:id")                                    //get food item by id
	foodItemsGroup.POST("/create", h.CreateFoodItem)              //create food item - admin only
	foodItemsGroup.PATCH("/update/:foodItemId", h.UpdateFoodItem) //updates a specific food item - admin only
	foodItemsGroup.DELETE("/delete/:id")                          //deletes a specific food item - admin only
}
