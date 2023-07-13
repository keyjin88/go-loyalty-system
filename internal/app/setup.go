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
	"time"
)

type API struct {
	ctx                context.Context
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
		ctx:    context.Background(),
		config: config.NewConfig(),
	}
}

func (api *API) Start() error {
	ctx, cancel := context.WithCancel(api.ctx)
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

	// Ожидаем получения сигнала ОС для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Получен сигнал остановки")
	// Отменяем контекст для graceful shutdown
	cancel()
	logger.Log.Infof("Running server. Address: %s |DB URI: %s |Gin release mode: %v |Log level: %s |accrual system address: %s",
		api.config.ServerAddress, api.config.DataBaseURI, api.config.GinReleaseMode, api.config.LogLevel, api.config.AccrualSystemAddress)
	srv := &http.Server{Addr: api.config.ServerAddress, Handler: api.router}

	// Останавливаем сервер с использованием контекста
	go func() {
		<-ctx.Done()
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelShutdown()
		conn, err := db.DB()
		if err != nil {
			logger.Log.Infof("Error while connecting to database: %v", err)
			return
		}
		err = conn.Close()
		if err != nil {
			logger.Log.Infof("Error while closing connection: %v", err)
			return
		}
		err = srv.Shutdown(ctxShutdown)
		if err != nil {
			logger.Log.Infof("Error while shutting down")
			return
		}
	}()
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
