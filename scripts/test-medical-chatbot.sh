#!/bin/bash

# Test script for BabyHome Medical Chatbot API

API_URL="http://localhost:8080"

echo "Testing BabyHome Medical Chatbot API..."

# Test health check
echo "1. Testing health check..."
curl -s "$API_URL/health" | jq '.' || echo "Health check failed"

echo -e "\n2. Testing welcome message..."
curl -s "$API_URL/whatsapp/welcome" | jq '.text.body' || echo "Welcome endpoint failed"

echo -e "\n3. Testing webhook verification..."
curl -s "$API_URL/whatsapp/webhook?hub.mode=subscribe&hub.verify_token=test_token&hub.challenge=test_challenge" || echo "Webhook verification failed"

echo -e "\n4. Testing medical consultation flow (Option A)..."
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
  }' | jq '.text.body' || echo "Medical consultation test failed"

echo -e "\n5. Testing studies reading flow (Option B)..."
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
            "from": "15551234568",
            "id": "wamid.xxx",
            "timestamp": "1234567890",
            "text": {
              "body": "B"
            },
            "type": "text"
          }]
        },
        "field": "messages"
      }]
    }]
  }' | jq '.text.body' || echo "Studies reading test failed"

echo -e "\n6. Testing appointment booking flow (Option C)..."
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
            "from": "15551234569",
            "id": "wamid.xxx",
            "timestamp": "1234567890",
            "text": {
              "body": "C"
            },
            "type": "text"
          }]
        },
        "field": "messages"
      }]
    }]
  }' | jq '.text.body' || echo "Appointment booking test failed"

echo -e "\n7. Testing BabyHome info flow (Option D)..."
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
            "from": "15551234570",
            "id": "wamid.xxx",
            "timestamp": "1234567890",
            "text": {
              "body": "D"
            },
            "type": "text"
          }]
        },
        "field": "messages"
      }]
    }]
  }' | jq '.text.body' || echo "BabyHome info test failed"

echo -e "\nMedical Chatbot API testing completed!"
