package food_items

import "github.com/gin-gonic/gin"

type Handler struct {
	store foodItemsStore
}

func NewHandler(store foodItemsStore) *Handler {
	return &Handler{
		store: store,
	}
}

func (h Handler) CreateFoodItem(ctx *gin.Context) {

}
