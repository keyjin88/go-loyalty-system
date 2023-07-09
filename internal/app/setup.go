package app

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/config"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/middleware/compressor"
	"net/http"
)

type API struct {
	config   *config.Config
	router   *gin.Engine
	handlers *handlers.Handler
}

func New() *API {
	return &API{
		config: config.NewConfig(),
	}
}

func (api *API) Start() error {
	if err := logger.Initialize(api.config.LogLevel); err != nil {
		return err
	}

	api.config.InitConfig()
	api.configureHandlers()
	api.configureRouter()

	logger.Log.Infof("Running server. Address: %s |DB URI: %s |Gin release mode: %v |Log level: %s |accrual system address: %s",
		api.config.ServerAddress, api.config.DataBaseURI, api.config.GinReleaseMode, api.config.LogLevel, api.config.AccrualSystemAddress)
	return http.ListenAndServe(api.config.ServerAddress, api.router)
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler()
}

func (api *API) configureRouter() {
	if api.config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	// router.Use(middleware.AuthMiddleware)
	router.Use(compressor.CompressionMiddleware())
	router.Use(gin.Logger())
	{
		router.POST("api/user/orders", func(c *gin.Context) { api.handlers.ProcessUserOrder(c) })
	}
	api.router = router
}
