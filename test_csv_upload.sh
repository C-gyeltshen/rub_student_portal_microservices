#!/bin/bash

# Test script for CSV upload to student management service

# Color codes for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
SERVICE_URL="http://localhost:8086"
CSV_FILE="new.csv"

echo -e "${YELLOW}Testing CSV Upload to Student Management Service${NC}\n"

# Check if CSV file exists
if [ ! -f "$CSV_FILE" ]; then
    echo -e "${RED}Error: CSV file '$CSV_FILE' not found!${NC}"
    exit 1
fi

echo -e "${YELLOW}Uploading CSV file: $CSV_FILE${NC}"
echo -e "${YELLOW}Target URL: $SERVICE_URL/api/students/bulk/csv${NC}\n"

# Upload the CSV file
response=$(curl -s -w "\n%{http_code}" -X POST \
  -F "file=@$CSV_FILE" \
  "$SERVICE_URL/api/students/bulk/csv")

# Split response and status code
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | sed '$d')

echo -e "${YELLOW}HTTP Status Code: $http_code${NC}"
echo -e "${YELLOW}Response:${NC}"
echo "$body" | jq '.' 2>/dev/null || echo "$body"

if [ "$http_code" = "201" ]; then
    echo -e "\n${GREEN}✓ CSV upload successful!${NC}"
else
    echo -e "\n${RED}✗ CSV upload failed!${NC}"
fi

echo -e "\n${YELLOW}You can also test with the old JSON endpoint:${NC}"
echo -e "curl -X POST $SERVICE_URL/api/students/bulk -H 'Content-Type: application/json' -d @students.json"
