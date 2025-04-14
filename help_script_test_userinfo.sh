#!/bin/bash
# Script to test and interact with the OAuth2 mock server
# This script:
# 1. Gets an authorization code from the /authorize endpoint
# 2. Exchanges the code for access tokens
# 3. Gets user info using the access token
# 4. Optionally updates the user info via config endpoint

set -e # Exit on error

# Default values - can be changed with command line arguments
HOST="localhost"
PORT="8080"
CLIENT_ID="test-client"
CLIENT_SECRET="test-secret"
REDIRECT_URI="http://localhost/callback"
UPDATE_USER=false

# ANSI colors for better output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Command line arguments
show_help() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -h, --host HOST          OAuth2 server host (default: localhost)"
    echo "  -p, --port PORT          OAuth2 server port (default: 8080)"
    echo "  -c, --client CLIENT_ID   Client ID (default: test-client)"
    echo "  -s, --secret SECRET      Client Secret (default: test-secret)"
    echo "  -r, --redirect URI       Redirect URI (default: http://localhost/callback)"
    echo "  -u, --update             Update user info after retrieval"
    echo "  --help                   Display this help message"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--host)
            HOST="$2"
            shift 2
            ;;
        -p|--port)
            PORT="$2"
            shift 2
            ;;
        -c|--client)
            CLIENT_ID="$2"
            shift 2
            ;;
        -s|--secret)
            CLIENT_SECRET="$2"
            shift 2
            ;;
        -r|--redirect)
            REDIRECT_URI="$2"
            shift 2
            ;;
        -u|--update)
            UPDATE_USER=true
            shift
            ;;
        --help)
            show_help
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

BASE_URL="http://${HOST}:${PORT}"
echo -e "${BLUE}OAuth2 Server URL: ${BASE_URL}${NC}"

echo -e "\n${YELLOW}Step 1: Getting authorization code...${NC}"
# Make a request to the authorize endpoint, with -i to include headers
AUTHORIZE_RESPONSE=$(curl -s -i "${BASE_URL}/authorize?client_id=${CLIENT_ID}&redirect_uri=${REDIRECT_URI}&response_type=code&scope=openid&state=test-state")
echo "Authorize response received..."

# Extract the authorization code from the Location header
LOCATION_HEADER=$(echo "$AUTHORIZE_RESPONSE" | grep -i "Location:")
if [ -z "$LOCATION_HEADER" ]; then
    echo -e "${RED}Failed to get Location header from authorize response${NC}"
    echo "$AUTHORIZE_RESPONSE"
    exit 1
fi

# Extract the code parameter using grep and sed
AUTH_CODE=$(echo "$LOCATION_HEADER" | grep -o "code=[^&]*" | sed "s/code=//")
if [ -z "$AUTH_CODE" ]; then
    echo -e "${RED}Failed to extract authorization code from response${NC}"
    echo "$LOCATION_HEADER"
    exit 1
fi

echo -e "${GREEN}Authorization code obtained: $AUTH_CODE${NC}"

echo -e "\n${YELLOW}Step 2: Exchanging code for tokens...${NC}"
# Exchange the authorization code for tokens
TOKEN_RESPONSE=$(curl -s -X POST "${BASE_URL}/token" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=authorization_code&code=${AUTH_CODE}&client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&redirect_uri=${REDIRECT_URI}")

# Extract the access token
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"access_token":"[^"]*' | sed 's/"access_token":"//')
if [ -z "$ACCESS_TOKEN" ]; then
    echo -e "${RED}Failed to extract access token from response${NC}"
    echo "$TOKEN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}Access token obtained: ${ACCESS_TOKEN}${NC}"

echo -e "\n${YELLOW}Step 3: Getting user info...${NC}"
# Get user info with the access token
USER_INFO_RESPONSE=$(curl -s -H "Authorization: Bearer ${ACCESS_TOKEN}" "${BASE_URL}/userinfo")

# Check if we got a valid JSON response
if ! echo "$USER_INFO_RESPONSE" | jq . &>/dev/null; then
    echo -e "${RED}Invalid JSON response from userinfo endpoint${NC}"
    echo "$USER_INFO_RESPONSE"
    exit 1
fi

echo -e "${GREEN}User info retrieved:${NC}"
echo "$USER_INFO_RESPONSE" | jq .

# If update flag is set, update the user info
if [ "$UPDATE_USER" = true ]; then
    echo -e "\n${YELLOW}Step 4: Updating user info...${NC}"
    
    # Create a JSON payload for the new user info
    # IMPORTANT CHANGE: Instead of sending "user_info", we need to send it in the token config
    # Based on examining the code in memory.go, the server looks for user_info within the tokens config
    NEW_USER_INFO=$(cat <<EOF
{
  "tokens": {
    "user_info": {
      "sub": "987654321",
      "id": "987654321",
      "name": "Updated Test User",
      "given_name": "Updated",
      "family_name": "User",
      "email": "updated@example.com",
      "email_verified": true,
      "picture": "https://example.com/updated.jpg",
      "locale": "en-US",
      "hd": "example.com"
    }
  }
}
EOF
)
    
    echo -e "${BLUE}New user info to be sent:${NC}"
    echo "$NEW_USER_INFO" | jq .
    
    # Send the update request
    UPDATE_RESPONSE=$(curl -s -X POST "${BASE_URL}/config" \
        -H "Content-Type: application/json" \
        -d "$NEW_USER_INFO")
    
    # Check if we got a valid JSON response
    if ! echo "$UPDATE_RESPONSE" | jq . &>/dev/null; then
        echo -e "${RED}Invalid JSON response from config endpoint${NC}"
        echo "$UPDATE_RESPONSE"
        exit 1
    fi
    
    echo -e "${GREEN}User info updated:${NC}"
    echo "$UPDATE_RESPONSE" | jq .
    
    echo -e "\n${YELLOW}Step 5: Getting updated user info...${NC}"
    # Get updated user info with the access token
    UPDATED_USER_INFO=$(curl -s -H "Authorization: Bearer ${ACCESS_TOKEN}" "${BASE_URL}/userinfo")
    
    echo -e "${GREEN}Updated user info retrieved:${NC}"
    echo "$UPDATED_USER_INFO" | jq .
fi

echo -e "\n${GREEN}OAuth2 flow completed successfully!${NC}"