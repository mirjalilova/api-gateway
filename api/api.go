package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mirjalilova/api-gateway/api/handlers"
	"github.com/mirjalilova/api-gateway/api/middlerware"
	_ "github.com/mirjalilova/api-gateway/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/exp/slog"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @BasePath  /router
// @description					Description for what is this security definition being used
func Engine(handler *handlers.Handlers) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for your specific origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	enforcer, err := casbin.NewEnforcer("gateway/casbin/model.conf", "gateway/casbin/policy.csv")
	if err != nil {
		slog.Error("Error while initializing Casbin: ", err)
	}
	// router.Use(middlerware.NewAuth(enforcer))

	router.GET("api/v1/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("swagger/doc.json")))

	med := router.Group("/medical-records").Use(middlerware.NewAuth(enforcer))
	{
		med.GET("", handler.ListMedicalRecords)
		med.GET("/:id", handler.GetMedicalRecord)
		med.POST("", handler.CreateMedicalRecord)
		med.PUT("/:id", handler.UpdateMedicalRecord)
		med.DELETE("/:id", handler.DeleteMedicalRecord)
	}

	genetic := router.Group("/genetic-data")
	{
		genetic.GET("", handler.ListGeneticData)
		genetic.GET("/:id", handler.GetGeneticData)
		genetic.POST("", handler.CreateGeneticData)
		genetic.PUT("/:id", handler.UpdateGeneticData)
		genetic.DELETE("/:id", handler.DeleteGeneticData)
	}

	life := router.Group("/lifestyle")
	{
		life.GET("", handler.ListLifestyle)
		life.GET("/:id", handler.GetLifestyle)
		life.POST("", handler.CreateLifestyle)
		life.PUT("/:id", handler.UpdateLifestyle)
		life.DELETE("/:id", handler.DeleteLifestyle)
	}

	wearable := router.Group("/wearable-data")
	{
		wearable.GET("", handler.ListWearableData)
		wearable.GET("/:id", handler.GetWearableData)
		wearable.POST("", handler.CreateWearableData)
		wearable.PUT("/:id", handler.UpdateWearableData)
		wearable.DELETE("/:id", handler.DeleteWearableData)
	}

	recommendation := router.Group("/recommendation")
	{
		recommendation.GET("", handler.ListRecommendations)
		recommendation.POST("", handler.CreateRecommendation)
		recommendation.GET("/:id", handler.GetRecommendation)
		recommendation.PUT("/:id", handler.UpdateRecommendation)
		recommendation.DELETE("/:id", handler.DeleteRecommendation)
	}

	monitoring := router.Group("/health-monitoring")
	{
		monitoring.GET("/real-time/:id", handler.HealthMonitoringRealTime)
		monitoring.GET("/daily-summary/:id", handler.HealthMonitoringDailyReport)
		monitoring.GET("/weekly-summary/:id", handler.HealthMonitoringWeeklyReport)
	}

	notification := router.Group("/notifications")
	{
		notification.POST("", handler.CreateNotification)
		notification.GET("/:id", handler.GetNotification)
		notification.GET("/", handler.ListNotifications)
	}

	return router
}
