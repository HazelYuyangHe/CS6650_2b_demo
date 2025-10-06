## Instructions

### Prepare Credentials

Retrieve you temporary credentials from Learner's Lab.
Enter your configuration when prompted:
```
aws configure
```

And set your session token:
```
aws configure set aws_session_token <YOUR-TEMP-SESSION-TOKEBN>
```

### Apply Infrastructure
```
cd terraform
terraform init
terraform apply -auto-approve
```

### Send Requests
Make sure you are in terraform folder

Get public ip address
```
aws ec2 describe-network-interfaces \
--network-interface-ids $(
    aws ecs describe-tasks \
    --cluster $(terraform output -raw ecs_cluster_name) \
    --tasks $(
        aws ecs list-tasks \
        --cluster $(terraform output -raw ecs_cluster_name) \
        --service-name $(terraform output -raw ecs_service_name) \
        --query 'taskArns[0]' --output text
    ) \
    --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" \
    --output text
) \
--query 'NetworkInterfaces[0].Association.PublicIp' \
--output text
```

Send some requests:
```
curl http://<PUBLIC-IP-ADDRESS>:8080/albums
```

## Clean Up
```
terraform destroy -auto-approve
```

# CS6650 HW5 - Product API

A RESTful Product API implementation based on OpenAPI specification, deployed on AWS ECS using Terraform.

## Project Structure

CS6650_2b_demo/
├── src/                    # Go source code
│   ├── main.go            # Entry point
│   ├── handlers.go        # HTTP handlers
│   ├── models.go          # Data models
│   ├── storage.go         # In-memory storage
│   ├── go.mod             # Go dependencies
│   └── go.sum
├── terraform/             # Infrastructure as Code
│   ├── main.tf           # Main configuration
│   ├── provider.tf       # Provider settings
│   ├── variables.tf      # Variables
│   ├── outputs.tf        # Outputs
│   └── modules/          # Terraform modules
├── Dockerfile            # Container definition
└── README.md

## API Endpoints

### 1. GET /products/{productId}
Get product details by ID.

### 2. POST /products/{productId}/details
Add or update product details.

## Local Development

### Prerequisites
- Go 1.21+
- Docker (optional)

### Run Locally
```bash
cd src
go run main.go handlers.go models.go storage.go
Server starts on port 8080 (configurable via PORT environment variable).
Deployment to AWS
Prerequisites

AWS CLI configured
Terraform installed
AWS Learner Lab credentials

Step 1: Configure AWS Credentials
Get temporary credentials from AWS Learner Lab:
bashaws configure
Set session token:
bashaws configure set aws_session_token <YOUR-SESSION-TOKEN>
Step 2: Deploy Infrastructure
bashcd terraform
terraform init
terraform apply -auto-approve
This creates:

ECR repository
ECS cluster
ECS service
Security groups
CloudWatch log groups

Step 3: Build and Push Docker Image
Note: Due to local Docker environment issues, we used AWS CloudShell for building:
bash# In AWS CloudShell
git clone https://github.com/HazelYuyangHe/CS6650_2b_demo.git
cd CS6650_2b_demo

# Login to ECR
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin 590184071414.dkr.ecr.us-west-2.amazonaws.com

# Build and push
docker build -t ecr_service:latest -f Dockerfile src/
docker tag ecr_service:latest 590184071414.dkr.ecr.us-west-2.amazonaws.com/ecr_service:latest
docker push 590184071414.dkr.ecr.us-west-2.amazonaws.com/ecr_service:latest
Step 4: Update ECS Service
bashaws ecs update-service --cluster CS6650L2-cluster --service CS6650L2 --force-new-deployment --region us-west-2
Step 5: Get Public IP Address
bashaws ec2 describe-network-interfaces \
--network-interface-ids $(
    aws ecs describe-tasks \
    --cluster CS6650L2-cluster \
    --tasks $(
        aws ecs list-tasks \
        --cluster CS6650L2-cluster \
        --service-name CS6650L2 \
        --query 'taskArns[0]' --output text
    ) \
    --query "tasks[0].attachments[0].details[?name=='networkInterfaceId'].value" \
    --output text \
    --region us-west-2
) \
--query 'NetworkInterfaces[0].Association.PublicIp' \
--output text \
--region us-west-2
Current Deployment IP: 44.243.181.3
API Testing Examples
All test examples below show the expected HTTP response codes defined in the OpenAPI specification.
Health Check
bashcurl http://44.243.181.3:8080/health
# Expected: OK
200 OK - Get Product Successfully
bash# First, add a product
curl -X POST http://44.243.181.3:8080/products/1/details \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "sku": "ABC-123",
    "manufacturer": "Test Corp",
    "category_id": 100,
    "weight": 500,
    "some_other_id": 999
  }'

# Then get it
curl http://44.243.181.3:8080/products/1
# Expected: {"product_id":1,"sku":"ABC-123",...}
204 No Content - Product Added Successfully
bashcurl -i -X POST http://44.243.181.3:8080/products/1/details \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "sku": "ABC-123",
    "manufacturer": "Test Corp",
    "category_id": 100,
    "weight": 500,
    "some_other_id": 999
  }'
# Expected: HTTP/1.1 204 No Content
400 Bad Request - Invalid Product ID
bashcurl -i http://44.243.181.3:8080/products/abc
# Expected: HTTP/1.1 400 Bad Request
# Response: {"error":"INVALID_INPUT","message":"Invalid product ID",...}
400 Bad Request - Invalid Input Data
bashcurl -i -X POST http://44.243.181.3:8080/products/1/details \
  -H "Content-Type: application/json" \
  -d '{"sku":""}'
# Expected: HTTP/1.1 400 Bad Request
# Response: {"error":"INVALID_INPUT",...}
404 Not Found - Product Does Not Exist
bashcurl -i http://44.243.181.3:8080/products/999
# Expected: HTTP/1.1 404 Not Found
# Response: {"error":"NOT_FOUND","message":"Product not found",...}
Implementation Notes
Docker Build Issue
During development, encountered EOF errors when building Docker images locally. This was resolved by:

Using AWS CloudShell for building images
Fixing Go version mismatch (go.mod required 1.25.1, but Dockerfile used 1.21)

Terraform Configuration

Modified main.tf to comment out automatic Docker image building
Manually built and pushed images to ECR
Used --force-new-deployment to update ECS service with new images

Clean Up
To avoid AWS charges, destroy resources when done:
bashcd terraform
terraform destroy -auto-approve
Architecture

Language: Go 1.21
Framework: gorilla/mux for routing
Storage: In-memory (HashMap)
Container: Docker with multi-stage build
Cloud: AWS ECS Fargate
IaC: Terraform
Region: us-west-2

Assignment Completion

✅ Part II: Product API implementation with validation
✅ Part III: Terraform infrastructure and AWS deployment
## Part IV - Load Testing Results

### Test Configuration
- Tool: Locust
- Load patterns tested: 10, 50, 100 concurrent users
- Test duration: ~3 minutes per test
- Request ratio: 75% GET, 25% POST

### Performance Results

| Users | RPS  | Median (ms) | 95%ile (ms) | Failure Rate |
|-------|------|-------------|-------------|--------------|
| 10    | 17.7 | 31-37       | ~38         | 22%          |
| 50    | 24.8 | 32-37       | ~42-76      | 49%          |
| 100   | 49.3 | 33-39       | 76-83       | 7%           |

### Key Findings
- System handles 100 concurrent users with 49.3 RPS
- Response times remain consistently fast (30-40ms median)
- Most "failures" are expected 404 responses for non-existent products
- Single ECS Fargate instance (0.5 vCPU, 1GB RAM) performs well

Author
Hazel Yuyang He
Course: CS6650 - Building Scalable Distributed Systems