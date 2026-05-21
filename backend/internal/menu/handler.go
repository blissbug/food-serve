package menu

import (
	"strconv"

	"food-serve.com/internal/middleware"
	"food-serve.com/pkg/response"
	"food-serve.com/pkg/types"
	"food-serve.com/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	MenuService *Service
}

func NewHandler(menuService *Service) *Handler {
	return &Handler{MenuService: menuService}
}

func (h Handler) CreateMenuHandler(ctx *gin.Context) {
	//get the payload
	var CreateMenuPayload types.CreateMenuPayload
	if err := ctx.ShouldBind(&CreateMenuPayload); err != nil {
		response.Error(ctx, err, 500)
		return
	}
	CreatedBy := ctx.GetUint(middleware.UserIDKey)
	value := ctx.Param("kitchenId")

	kitchenId, err := strconv.Atoi(value)

	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	err = h.MenuService.CreateMenuService(CreateMenuPayload, CreatedBy, uint(kitchenId))

	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	response.JSON(ctx, "Menu created successfully")
}

func (h Handler) PublishMenuHandler(ctx *gin.Context) {
	//we get a menu id, check if the date is of today or later and publish it
	//also if the time slots were empty initially, we fill em up
	menuId, err := utils.GetParamInUint("menuId", ctx)
	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	kitchenId, err := utils.GetParamInUint("kitchenId", ctx)

	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	var PublishMenuPayload types.PublishMenuPayload
	if err := ctx.ShouldBind(&PublishMenuPayload); err != nil {
		response.Error(ctx, err, 500)
		return
	}

	serviceErr := h.MenuService.PublishMenuService(PublishMenuPayload, menuId, kitchenId)
	if serviceErr != nil {
		response.Error(ctx, serviceErr, 500)
		return
	}
	response.JSON(ctx, "Menu published successfully")
}

func (h Handler) UpdateMenuHandler(ctx *gin.Context) {
	//if is active -> we can only update the availability of the items in menu
	//like they set something, and I fetch from the food items list in payload, close timing which cannot be before that current time ofc
	menuId, err := utils.GetParamInUint("menuId", ctx)
	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	kitchenId, err := utils.GetParamInUint("kitchenId", ctx)
	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	var UpdateMenuPayload types.UpdateMenuPayload
	if err := ctx.ShouldBind(&UpdateMenuPayload); err != nil {
		response.Error(ctx, err, 500)
		return
	}

	h.MenuService.UpdateMenuService(UpdateMenuPayload, menuId, kitchenId)
	//if its draft -> we can update the menu items, open and close time, price, date
}

func (h Handler) DeleteMenuHandler(ctx *gin.Context) {
	menuId, err := utils.GetParamInUint("menuId", ctx)
	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	kitchenId, err := utils.GetParamInUint("kitchenId", ctx)
	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	err = h.MenuService.DeleteMenuService(menuId, kitchenId)

	if err != nil {
		response.Error(ctx, err, 500)
		return
	}

	response.JSON(ctx, "Menu deleted successfully")
}
