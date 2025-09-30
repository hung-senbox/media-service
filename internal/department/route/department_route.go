package route

import (
	"media-service/internal/department/handler"
	"media-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDepartmentRoutes(r *gin.Engine, h *handler.DepartmentHandler, hRegion *handler.RegionHandler) {
	// Admin routes
	adminGroup := r.Group("/api/v1/admin")
	adminGroup.Use(middleware.Secured())
	{
		departmentsAdmin := adminGroup.Group("/departments")
		{
			departmentsAdmin.GET("", h.GetDepartments4Web)
			departmentsAdmin.GET("/:department_id", h.GetDepartmentDetail4Web)
			departmentsAdmin.POST("", h.UploadDepartment)
			departmentsAdmin.PUT("/:department_id", h.UpdateDepartment)
			departmentsAdmin.POST("/home-menus", h.UploadDepartmentMenu)
			departmentsAdmin.POST("/organization-menus", h.UploadDepartmentMenuOrganization)

			// assign - remove
			departmentsAdmin.POST("/assign/leader", h.AssignLeader)
			departmentsAdmin.POST("/assign/staff", h.AssignStaff)
			departmentsAdmin.POST("/remove/staff", h.RemoveStaff)
			departmentsAdmin.POST("/remove/leader", h.RemoveLeader)

			departmentsAdmin.POST("/regions", hRegion.CreateRegion)
			departmentsAdmin.PUT("/regions/:region_id", hRegion.UpdateRegionName)
		}
	}

	// User routes
	userGroup := r.Group("/api/v1/user")
	userGroup.Use(middleware.Secured())
	{
		departmentsUser := userGroup.Group("/departments")
		{
			departmentsUser.GET("", h.GetDepartments4App)
		}
	}

	// gateway routes
	gatewayGroup := r.Group("/api/v1/gateway")
	gatewayGroup.Use(middleware.Secured())
	{
		departmentsGateway := gatewayGroup.Group("/departments")
		{
			departmentsGateway.GET("", h.GetDepartments4Gateway)
			departmentsGateway.GET("/organization/:organization_id", h.GetDepartmentsByOrganization4Gateway)
		}
	}

}
