package router

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/mamh1019/lambda-go-server-boilerplate/db"
	"github.com/mamh1019/lambda-go-server-boilerplate/handler"
	"github.com/mamh1019/lambda-go-server-boilerplate/kafka"
	"github.com/mamh1019/lambda-go-server-boilerplate/order"
	"github.com/mamh1019/lambda-go-server-boilerplate/user"
)

func SetupRouter() *gin.Engine {
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	d, err := db.NewDB()
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	var rdb *redis.Client
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	kafkaProducer := kafka.NewProducer()

	userRepo := user.NewRepository(d, rdb)
	userHandler := handler.NewUserHandler(userRepo, kafkaProducer)

	orderRepo := order.NewRepository(d)
	orderHandler := handler.NewOrderHandler(orderRepo)

	r.GET("/health", handler.Health)

	r.GET("/users", userHandler.List)
	r.GET("/users/:id", userHandler.Get)
	r.POST("/users", userHandler.Create)
	r.PUT("/users/:id", userHandler.Update)
	r.DELETE("/users/:id", userHandler.Delete)

	r.GET("/orders", orderHandler.List)
	r.GET("/orders/:id", orderHandler.Get)
	r.POST("/orders", orderHandler.Create)
	r.PUT("/orders/:id", orderHandler.Update)
	r.DELETE("/orders/:id", orderHandler.Delete)

	return r
}
