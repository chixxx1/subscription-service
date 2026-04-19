package transport_http

import (
	_ "github.com/chixxx1/subscription-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	sub_service "github.com/chixxx1/subscription-service/internal/service/subscription"
	"github.com/chixxx1/subscription-service/internal/transport/http/handler"
	"github.com/chixxx1/subscription-service/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func InitRoutes(svc *sub_service.SubscriptionService, logger *zap.Logger) *gin.Engine {
	r := gin.New()

	r.Use(middleware.Logger(logger))
	r.Use(gin.Recovery())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		subHandler := handler.NewSubscriptionHandler(svc, logger)
		subscriptions := api.Group("/subscriptions")
		{
			subscriptions.POST("", subHandler.CreateSubscription)
			subscriptions.GET("/:id", subHandler.GetByID)
			subscriptions.GET("", subHandler.List)
			subscriptions.PUT("/:id", subHandler.Update)
			subscriptions.DELETE("/:id", subHandler.Delete)
			subscriptions.GET("/total-cost", subHandler.GetTotalCost)
		}
	}

	return r
}
