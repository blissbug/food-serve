package food_items

import "github.com/gin-gonic/gin"

func (h Handler) RegisterRoutes(r *gin.RouterGroup) {
	foodItemsGroup := r.Group("/kitchen/:kitchenId/food/items")
	foodItemsGroup.GET("/", h.GetFoodItems) //get all the food items, paginated
	foodItemsGroup.GET("/:id")              //get food item by id
	//create food item - admin only
	foodItemsGroup.POST("/create", h.CreateFoodItem)
	foodItemsGroup.PATCH("/update/:foodItemId", h.UpdateFoodItem)  //updates a specific food item - admin only
	foodItemsGroup.DELETE("/delete/:foodItemId", h.DeleteFoodItem) //deletes a specific food item - admin only
}
