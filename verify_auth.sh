#!/bin/bash

BASE_URL="http://localhost:8000"
EMAIL="test_user_$(date +%s)@example.com"
PASSWORD="password123"

echo "1. Registering user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")
echo "Response: $REGISTER_RESPONSE"

if [[ $REGISTER_RESPONSE == *"id"* ]]; then
  echo "✅ Registration successful"
else
  echo "❌ Registration failed"
  exit 1
fi

echo -e "\n2. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")
echo "Response: $LOGIN_RESPONSE"

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [[ -n "$ACCESS_TOKEN" ]]; then
  echo "✅ Login successful. Token: ${ACCESS_TOKEN:0:20}..."
else
  echo "❌ Login failed"
  exit 1
fi

echo -e "\n3. Logging out..."
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "Response: $LOGOUT_RESPONSE"

if [[ $LOGOUT_RESPONSE == *"Successfully logged out"* ]]; then
  echo "✅ Logout successful"
else
  echo "❌ Logout failed"
  exit 1
fi

echo -e "\n4. Verifying token blacklist (should fail)..."
PROTECTED_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/logout" \
  -H "Authorization: Bearer $ACCESS_TOKEN")
echo "Response: $PROTECTED_RESPONSE"

if [[ $PROTECTED_RESPONSE == *"Token is blacklisted"* ]]; then
  echo "✅ Token blacklist verified"
else
  echo "❌ Token blacklist verification failed (expected 'Token is blacklisted')"
fi
