#!/bin/bash

# BitPin API Token Refresh Script
# Usage: ./refresh_token.sh

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# API endpoint
API_URL="https://api.bitpin.ir/api/v1/usr/refresh_token/"

# Refresh token
REFRESH_TOKEN="eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJ0b2tlbl90eXBlIjoicmVmcmVzaCIsImV4cCI6MTc1MTE3NTY5OCwianRpIjoiN2E1MWJhNGM4NDI3NGM5Nzg1ZmMwNTVmOTEzYmU2MzYiLCJ1c2VyX2lkIjo0OTM0MDk0LCJpcCI6WyI1LjEyNy4yMTAuMjE2IiwiMTA0LjI4LjE5Ny4xNCIsIjUuMjMzLjE5Ni41MSIsIjEwNC4yOC4yMjkuMTMiLCIyLjE4NS40OS43NyIsIjUuMTI3LjIxMi4xNzciXSwiYXBpX2NyZWRlbnRpYWxfaWQiOjMxNDN9.cC8YiI-bTIu_VAMKx8fnuWITlkd5l5V32XPlFn9jk3AU26UBqTWej_DYkG6L75-2C4PDKFxKuAPc5MN-o8Ge0Hj9hUJB8bb4m-KTu4u9ggzRJu1rhEu2-GnEPgV45L0PDiACoJzrtfgBUtO9HrJKbN3gulLBaon_NpnjMDGIqkhmyTmfaa8Ru7TbVLokiMZ7_hvAy6e-BPYVfiMfx1F0VqV4w60MsVSLpq26xoRMHmpYWyDF-m3YskcXsqP12qxvXQ_AOc4pDQ5BL4Nz6hv8Vx80NJnODWePsMzBseetIMqT4bvnL4a-uyqobMXXa0PSKnlEFu5PsEW-1v6ToZeqsnEijyZ6kEgnkJWSbxkO8xaOdRy3khEeQYMnHweFIL_p7gLhfza80u0xsGU6HiDHAM6FgVy-2LhLufx9KjIJNNj7EPsDLQqrGlPbUt6J7hxLgx4oo2tdHul4fqJnhWtje3317UCEHRBvbIo-XmnoDMGcyF6nkk7TmhNxNxUx1Bk5rhEV3GRHD8jvVidF2rssnpuBZQZU42xmpFX59OnYJLbdrL8XgUQioP0TmVeANML_h6l7-VyNslHaCknFqFVcwEsZSF1dhJiepTq5nv-1nOmtsLXLbvpUXQjNmGDmg88JVMKjhcIpkv9YimZcSLjAgkYX-P1TkVjXXbITwbK1Vlk"

echo -e "${YELLOW}Sending token refresh request to BitPin API...${NC}"
echo "URL: $API_URL"
echo ""

# Send the curl request and capture response and HTTP status
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$API_URL" \
    -H "Content-Type: application/json" \
    -d "{\"refresh\":\"$REFRESH_TOKEN\"}")

# Split response body and HTTP status code
HTTP_BODY=$(echo "$RESPONSE" | head -n -1)
HTTP_STATUS=$(echo "$RESPONSE" | tail -n 1)

echo -e "${YELLOW}Response Status Code:${NC} $HTTP_STATUS"
echo ""

# Check if the request was successful
if [ "$HTTP_STATUS" -eq 200 ]; then
    echo -e "${GREEN}✓ Success! Token refreshed successfully.${NC}"
    echo ""
    echo -e "${YELLOW}Response Body:${NC}"
    echo "$HTTP_BODY" | python3 -m json.tool 2>/dev/null || echo "$HTTP_BODY"

    # Extract new tokens if present
    if command -v jq &> /dev/null; then
        ACCESS_TOKEN=$(echo "$HTTP_BODY" | jq -r '.access // empty')
        NEW_REFRESH_TOKEN=$(echo "$HTTP_BODY" | jq -r '.refresh // empty')

        if [ ! -z "$ACCESS_TOKEN" ] && [ "$ACCESS_TOKEN" != "null" ]; then
            echo ""
            echo -e "${GREEN}New Access Token:${NC}"
            echo "$ACCESS_TOKEN"
        fi

        if [ ! -z "$NEW_REFRESH_TOKEN" ] && [ "$NEW_REFRESH_TOKEN" != "null" ]; then
            echo ""
            echo -e "${GREEN}New Refresh Token:${NC}"
            echo "$NEW_REFRESH_TOKEN"
        fi
    fi
else
    echo -e "${RED}✗ Request failed with status code: $HTTP_STATUS${NC}"
    echo ""
    echo -e "${YELLOW}Error Response:${NC}"
    echo "$HTTP_BODY" | python3 -m json.tool 2>/dev/null || echo "$HTTP_BODY"
fi

echo ""
echo -e "${YELLOW}Raw curl command used:${NC}"
echo "curl -X POST '$API_URL' -H 'Content-Type: application/json' -d '{\"refresh\":\"$REFRESH_TOKEN\"}'"
