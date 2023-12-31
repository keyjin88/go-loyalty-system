package app

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/keyjin88/go-loyalty-system/internal/app/config"
	"github.com/keyjin88/go-loyalty-system/internal/app/daemons"
	"github.com/keyjin88/go-loyalty-system/internal/app/handlers"
	"github.com/keyjin88/go-loyalty-system/internal/app/logger"
	"github.com/keyjin88/go-loyalty-system/internal/app/middleware"
	"github.com/keyjin88/go-loyalty-system/internal/app/middleware/compressor"
	"github.com/keyjin88/go-loyalty-system/internal/app/model/entities"
	"github.com/keyjin88/go-loyalty-system/internal/app/services"
	"github.com/keyjin88/go-loyalty-system/internal/app/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type API struct {
	config             *config.Config
	router             *gin.Engine
	handlers           *handlers.Handler
	userService        *services.UserService
	orderService       *services.OrderService
	withdrawService    *services.WithdrawService
	userRepository     *storage.UserRepository
	orderRepository    *storage.OrderRepository
	withdrawRepository *storage.WithdrawRepository
}

func New() *API {
	return &API{
		config: config.NewConfig(),
	}
}
func (api *API) Start() error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := logger.Initialize(api.config.LogLevel); err != nil {
		return err
	}

	api.config.InitConfig()
	api.configureRouter()
	db := api.ConfigDBConnection()
	api.configStorage(db)
	// Канал для обработки заказов через сервер Accrual
	// Если уже есть пулл горутин, то насколько важна буферизация канала? Или я чего-то не понял?
	orderProcessingChannel := make(chan entities.Order, api.config.ProcessingChannelBufferSize)
	mutex := &sync.Mutex{}
	api.configService(orderProcessingChannel, mutex)
	api.configHandlers()
	api.configWorkers(db, orderProcessingChannel, mutex)

	// Создаем HTTP-сервер
	srv := &http.Server{
		Addr:    api.config.ServerAddress,
		Handler: api.router,
	}

	// Запускаем HTTP-сервер в отдельной горутине
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Infof("Error while start server")
		}
	}()
	logger.Log.Infof("Server started")
	// Ожидаем получения сигнала остановки
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Infof("Stop signal received")
	// Отменяем контекст для graceful shutdown
	cancel()
	// Устанавливаем таймаут для graceful shutdown
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	// Останавливаем HTTP-сервер
	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Log.Infof("Error shutting down")
	}
	log.Println("Сервер остановлен")
	return nil
}

func (api *API) ConfigDBConnection() *gorm.DB {
	// Создание пула соединений
	db, err := gorm.Open(postgres.Open(api.config.DataBaseURI), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Установка максимального количества подключений в пуле
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	return db
}

func (api *API) configHandlers() {
	api.handlers = handlers.NewHandler(api.userService, api.orderService, api.withdrawService, api.config.SecretKey)
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
		protectedGroup.GET("api/user/orders", func(c *gin.Context) { api.handlers.GetAllOrders(c) })
		protectedGroup.GET("api/user/balance", func(c *gin.Context) { api.handlers.GetBalance(c) })
		protectedGroup.GET("api/user/withdrawals", func(c *gin.Context) { api.handlers.GetAllWithdrawals(c) })
		protectedGroup.POST("api/user/balance/withdraw", func(c *gin.Context) { api.handlers.SaveWithdraw(c) })
	}
	api.router = router
}

func (api *API) configStorage(db *gorm.DB) {
	api.userRepository = storage.NewUserRepository(db)
	api.orderRepository = storage.NewOrderRepository(db)
	api.withdrawRepository = storage.NewWithdrawRepository(db)
}

func (api *API) configService(channel chan entities.Order, mutex *sync.Mutex) {
	api.userService = services.NewUserService(api.userRepository)
	api.withdrawService = services.NewWithdrawService(api.withdrawRepository, api.userRepository, mutex)
	api.orderService = services.NewOrderService(
		api.orderRepository,
		channel,
	)
}

func (api *API) configWorkers(db *gorm.DB, channel chan entities.Order, mutex *sync.Mutex) {
	go daemons.WorkerProcessingOrders(channel, api.config.AccrualSystemAddress, db, api.config.WorkerPoolSize, mutex)
}
