package router

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/mamh1019/go-boilerplate/db"
	"github.com/mamh1019/go-boilerplate/handler"
	"github.com/mamh1019/go-boilerplate/kafka"
	"github.com/mamh1019/go-boilerplate/order"
	"github.com/mamh1019/go-boilerplate/user"
)

// SetupRouter는 모든 라우팅을 설정하고 gin.Engine 을 반환합니다.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// DB 초기화
	d, err := db.NewDB()
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	// Redis 초기화 (없으면 nil 로 두고, 캐시는 비활성화)
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	var rdb *redis.Client
	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Kafka 프로듀서 초기화
	kafkaProducer := kafka.NewProducer()

	// 리포지토리 & 핸들러
	userRepo := user.NewRepository(d, rdb)
	userHandler := handler.NewUserHandler(userRepo, kafkaProducer)

	orderRepo := order.NewRepository(d)
	orderHandler := handler.NewOrderHandler(orderRepo)

	// 헬스체크
	r.GET("/health", handler.Health)

	// 사용자 CRUD
	r.GET("/users", userHandler.List)
	r.GET("/users/:id", userHandler.Get)
	r.POST("/users", userHandler.Create)
	r.PUT("/users/:id", userHandler.Update)
	r.DELETE("/users/:id", userHandler.Delete)

	// 주문 CRUD
	r.GET("/orders", orderHandler.List)
	r.GET("/orders/:id", orderHandler.Get)
	r.POST("/orders", orderHandler.Create)
	r.PUT("/orders/:id", orderHandler.Update)
	r.DELETE("/orders/:id", orderHandler.Delete)

	return r
}
