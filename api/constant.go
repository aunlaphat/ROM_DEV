package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *Application) ConstantRoute(apiRouter *gin.RouterGroup) {
	constant := apiRouter.Group("/constant")

	constant.GET("/roles", app.GetRoles)
	constant.GET("/warehouses", app.GetWarehouses)
}

// @Summary Get roles
// @Description Get all available roles
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]response.RoleResponse}
// @Failure 500 {object} Response
// @Router /constant/roles [get]
func (app *Application) GetRoles(c *gin.Context) {
	app.Logger.Info("📋 Fetching roles")

	roles, err := app.Service.Constant.GetRoles(c.Request.Context())
	if err != nil {
		app.Logger.Error("❌ Failed to fetch roles", zap.Error(err))
		handleResponse(c, false, "❌ Failed to fetch roles", nil, http.StatusInternalServerError)
		return
	}

	app.Logger.Info("✅ Roles retrieved successfully", zap.Int("count", len(roles)))
	handleResponse(c, true, "⭐ Roles retrieved successfully ⭐", roles, http.StatusOK)
}

// @Summary Get warehouses
// @Description Get all available warehouses
// @Tags Constants
// @Accept json
// @Produce json
// @Success 200 {object} Response{data=[]response.WarehouseResponse}
// @Failure 500 {object} Response
// @Router /constant/warehouses [get]
func (app *Application) GetWarehouses(c *gin.Context) {
	app.Logger.Info("📋 Fetching warehouses")

	warehouses, err := app.Service.Constant.GetWarehouses(c.Request.Context())
	if err != nil {
		app.Logger.Error("❌ Failed to fetch warehouses", zap.Error(err))
		handleResponse(c, false, "❌ Failed to fetch warehouses", nil, http.StatusInternalServerError)
		return
	}

	app.Logger.Info("✅ Warehouses retrieved successfully", zap.Int("count", len(warehouses)))
	handleResponse(c, true, "⭐ Warehouses retrieved successfully ⭐", warehouses, http.StatusOK)
}
