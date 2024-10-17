package v1

import (
	"net/http"

	// Need for swagger rendering
	_ "vio_coding_challenge/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const ipParam = "ip_address"

type GeoCtrlI interface {
	Get(c *gin.Context)
}

// AddRoutes -.
// Swagger spec:
func AddRoutes(handler *gin.Engine, geoCtrl GeoCtrlI) {
	// Middleware
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	handler.GET("/swagger/*any", swaggerHandler)

	// Liveness probe
	handler.GET("/liveness", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Handlers
	v1 := handler.Group("/api/v1")
	{
		geolocations := v1.Group("/geolocations")
		{
			geolocations.GET("/:"+ipParam, geoCtrl.Get)
		}
	}
}
