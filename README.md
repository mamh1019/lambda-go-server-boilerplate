## Go Boilerplate

Minimal starter template for a Go web service using Gin, MySQL, Redis, and Kafka.

### Quick start (local)

```bash
git clone git@github.com:mamh1019/lambda-go-server-boilerplate.git
cd lambda-go-server-boilerplate
go mod tidy
go run ./...
```

### Deploy to AWS Lambda

Lambda entrypoint is in `lambda/main.go`. You can deploy in two ways:

**Manual build & upload**

```bash
GOOS=linux GOARCH=amd64 go build -o bootstrap ./lambda
zip function.zip bootstrap
```

Upload `function.zip` to AWS Lambda with:
- Runtime: `provided.al2`
- Handler: `bootstrap`
- Connect via API Gateway HTTP API or ALB to expose the HTTP endpoints.

**Using Makefile**

```bash
make deploy           # build build/function.zip
FUNCTION_NAME=lambda-go-server-boilerplate-dev make deploy-aws
```

**GitHub Actions (CI/CD)**

On push:
- `main` branch → deploys to `lambda-go-server-boilerplate-prod`
- `dev` branch → deploys to `lambda-go-server-boilerplate-dev`

Configured in `.github/workflows/deploy-lambda.yml` (requires AWS credentials in repository secrets).

### Environment

- `MYSQL_DSN`
- `REDIS_ADDR`
- `KAFKA_BROKERS`

# go-boilerplate
Opinionated Go boilerplate for REST APIs with Gin, MySQL, Redis caching, and Kafka integration.
