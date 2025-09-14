#!/bin/bash

# WhatsApp Chatbot Deployment Script
# This script deploys the application to AWS

set -e

# Configuration
AWS_REGION=${AWS_REGION:-us-east-1}
ENVIRONMENT=${ENVIRONMENT:-dev}
ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REPOSITORY="chatbot-wsp"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}Starting deployment process...${NC}"

# Check if AWS CLI is installed
if ! command -v aws &> /dev/null; then
    echo -e "${RED}AWS CLI is not installed. Please install it first.${NC}"
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Docker is not installed. Please install it first.${NC}"
    exit 1
fi

# Login to ECR
echo -e "${YELLOW}Logging in to ECR...${NC}"
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com

# Create ECR repository if it doesn't exist
echo -e "${YELLOW}Creating ECR repository if it doesn't exist...${NC}"
aws ecr describe-repositories --repository-names $ECR_REPOSITORY --region $AWS_REGION 2>/dev/null || \
aws ecr create-repository --repository-name $ECR_REPOSITORY --region $AWS_REGION

# Build Docker image
echo -e "${YELLOW}Building Docker image...${NC}"
docker build -t $ECR_REPOSITORY:latest .

# Tag image for ECR
docker tag $ECR_REPOSITORY:latest $ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPOSITORY:latest

# Push image to ECR
echo -e "${YELLOW}Pushing image to ECR...${NC}"
docker push $ACCOUNT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/$ECR_REPOSITORY:latest

# Deploy using CloudFormation
echo -e "${YELLOW}Deploying CloudFormation stack...${NC}"
aws cloudformation deploy \
    --template-file aws/cloudformation/template.yaml \
    --stack-name $ENVIRONMENT-chatbot-wsp \
    --parameter-overrides \
        Environment=$ENVIRONMENT \
        WhatsAppVerifyToken=$WHATSAPP_VERIFY_TOKEN \
        WhatsAppAccessToken=$WHATSAPP_ACCESS_TOKEN \
    --capabilities CAPABILITY_IAM \
    --region $AWS_REGION

# Get API URL
API_URL=$(aws cloudformation describe-stacks \
    --stack-name $ENVIRONMENT-chatbot-wsp \
    --query 'Stacks[0].Outputs[?OutputKey==`ApiUrl`].OutputValue' \
    --output text \
    --region $AWS_REGION)

echo -e "${GREEN}Deployment completed successfully!${NC}"
echo -e "${GREEN}API URL: $API_URL${NC}"
echo -e "${GREEN}Webhook URL: ${API_URL}whatsapp/webhook${NC}"

# Test the deployment
echo -e "${YELLOW}Testing deployment...${NC}"
curl -f $API_URL/health || echo -e "${RED}Health check failed${NC}"
