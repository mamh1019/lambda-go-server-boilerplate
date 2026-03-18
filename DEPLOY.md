# Lambda ZIP deployment

## 1. Create the Lambda function (one time)

1. In the AWS console, go to **Lambda → Create function → Author from scratch**.
2. **Function name**:  
   - `lambda-go-server-boilerplate-prod` (production), or  
   - `lambda-go-server-boilerplate-dev` (development).
3. **Runtime**: **Amazon Linux 2023 (provided.al2023)** or **Amazon Linux 2 (provided.al2)**.
4. **Architecture**: `x86_64`.
5. After creation, go to **Configuration → General configuration → Edit**:
   - **Handler**: `bootstrap`
6. Optionally attach **Function URL**, **API Gateway**, or **ALB** to expose HTTP endpoints.

## 2. Local deploy with Makefile

```bash
# 1) Build ZIP only
make deploy

# 2) Build ZIP + upload to Lambda (requires AWS CLI)
# You must provide FUNCTION_NAME and AWS_REGION, e.g.:
FUNCTION_NAME=lambda-go-server-boilerplate-dev AWS_REGION=ap-northeast-1 make deploy-aws
```

Example environment (shell export or CI variables):

```bash
export FUNCTION_NAME=lambda-go-server-boilerplate-dev
export AWS_REGION=ap-northeast-1
```

## 3. GitHub Actions

- On `main` branch push: deploys to `lambda-go-server-boilerplate-prod`
- On `dev` branch push: deploys to `lambda-go-server-boilerplate-dev`
- Required repository **Secrets**:
  - `AWS_ACCESS_KEY_ID`
  - `AWS_SECRET_ACCESS_KEY`
  - `AWS_REGION`

