#!/bin/bash

# Test script for WhatsApp Chatbot API

API_URL="http://localhost:8080"

echo "Testing WhatsApp Chatbot API..."

# Test health check
echo "1. Testing health check..."
curl -s "$API_URL/health" | jq '.' || echo "Health check failed"

echo -e "\n2. Testing stats endpoint..."
curl -s "$API_URL/stats" | jq '.' || echo "Stats endpoint failed"

echo -e "\n3. Testing welcome message..."
curl -s "$API_URL/whatsapp/welcome" | jq '.' || echo "Welcome endpoint failed"

echo -e "\n4. Testing webhook verification..."
curl -s "$API_URL/whatsapp/webhook?hub.mode=subscribe&hub.verify_token=test_token&hub.challenge=test_challenge" || echo "Webhook verification failed"

echo -e "\n5. Testing webhook with sample message..."
curl -X POST "$API_URL/whatsapp/webhook" \
  -H "Content-Type: application/json" \
  -d '{
    "object": "whatsapp_business_account",
    "entry": [{
      "id": "ENTRY_ID",
      "changes": [{
        "value": {
          "messaging_product": "whatsapp",
          "metadata": {
            "display_phone_number": "15551234567",
            "phone_number_id": "PHONE_NUMBER_ID"
          },
          "messages": [{
            "from": "15551234567",
            "id": "wamid.xxx",
            "timestamp": "1234567890",
            "text": {
              "body": "A"
            },
            "type": "text"
          }]
        },
        "field": "messages"
      }]
    }]
  }' | jq '.' || echo "Webhook test failed"

echo -e "\nAPI testing completed!"
