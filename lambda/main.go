package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/joho/godotenv"
	"github.com/mamh1019/lambda-go-server-boilerplate/router"
)

var ginLambda *ginadapter.GinLambdaV2

func loadEnv() {
	stage := os.Getenv("STAGE")
	if stage == "" {
		stage = "lambda"
	}
	_ = godotenv.Load(".env." + stage)
	_ = godotenv.Load(".env")
}

func init() {
	loadEnv()

	r := router.SetupRouter()
	ginLambda = ginadapter.NewV2(r)
}

func main() {
	lambda.Start(ginLambda.ProxyWithContext)
}
