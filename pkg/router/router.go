package router

import (
	"media-service/internal/department/handler"
	"media-service/internal/department/repository"
	"media-service/internal/department/route"
	"media-service/internal/department/service"
	"media-service/internal/gateway"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(consulClient *api.Client, departmentCollection *mongo.Collection, regionCollection *mongo.Collection) *gin.Engine {
	r := gin.Default()

	// gateway
	userGateway := gateway.NewUserGateway("go-main-service", consulClient)
	menuGateway := gateway.NewMenuGateway("go-main-service", consulClient)
	messageLangGw := gateway.NewMessageLanguageGateway("go-main-service", consulClient)
	classroomGW := gateway.NewClassroomGateway("inventory-service", consulClient)

	// region
	regionRepo := repository.NewRegionRepository(regionCollection)
	regionService := service.NewRegionService(regionRepo, userGateway)
	regionHandler := handler.NewRegionHandler(regionService)

	// department
	departmentRepo := repository.NewDepartmentRepository(departmentCollection)
	departmentService := service.NewDepartmentService(departmentRepo, userGateway, messageLangGw, menuGateway, classroomGW, regionRepo)
	departmentHandler := handler.NewDepartmentHandler(departmentService)

	// Register routes
	route.RegisterDepartmentRoutes(r, departmentHandler, regionHandler)
	//route.RegisterRegionRoutes(r, regionHandler)
	return r
}
