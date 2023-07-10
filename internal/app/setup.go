package app

import (
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/config"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/middleware"
	"github.com/keyjin88/go-loyalty-system/internal/app/middleware/compressor"
	"github.com/keyjin88/go-loyalty-system/internal/app/services"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type API struct {
	config          *config.Config
	router          *gin.Engine
	handlers        *handlers.Handler
	userService     *services.UserService
	orderService    *services.OrderService
	userRepository  *storage.UserRepository
	orderRepository *storage.OrderRepository
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
	api.configureRouter()
	api.configStorage()
	api.configService()
	api.configureHandlers()

	logger.Log.Infof("Running server. Address: %s |DB URI: %s |Gin release mode: %v |Log level: %s |accrual system address: %s",
		api.config.ServerAddress, api.config.DataBaseURI, api.config.GinReleaseMode, api.config.LogLevel, api.config.AccrualSystemAddress)
	return http.ListenAndServe(api.config.ServerAddress, api.router)
}

func (api *API) configureHandlers() {
	api.handlers = handlers.NewHandler(api.userService, api.orderService, api.config.SecretKey)
}

func (api *API) configureRouter() {
	if api.config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(compressor.CompressionMiddleware())
	router.Use(gin.Logger())
	authGroup := router.Group("/")
	{
		authGroup.POST("api/user/register", func(c *gin.Context) { api.handlers.RegisterUser(c) })
		authGroup.POST("api/user/login", func(c *gin.Context) { api.handlers.LoginUser(c) })
	}
	protectedGroup := router.Group("/")
	protectedGroup.Use(middleware.AuthMiddleware(api.config.SecretKey))
	{
		protectedGroup.POST("api/user/orders", func(c *gin.Context) { api.handlers.ProcessUserOrder(c) })
	}
	api.router = router
}

func (api *API) configStorage() {
	db, err := gorm.Open(postgres.Open(api.config.DataBaseURI), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database")
	}
	api.userRepository = storage.NewUserRepository(db)
	api.orderRepository = storage.NewOrderRepository(db)
}

func (api *API) configService() {
	api.userService = services.NewUserService(api.userRepository)
	api.orderService = services.NewOrderService(api.orderRepository)
}
