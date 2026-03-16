package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/mamh1019/go-boilerplate/router"
)

func main() {
	_ = godotenv.Load(".env")

	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
