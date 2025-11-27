#!/bin/bash

# Test script for Finance Service Integration
# This script tests the integration between Student Management Service and Finance Service

BASE_URL="http://localhost:8084"
FINANCE_URL="http://localhost:8085"

echo "=================================="
echo "Finance Service Integration Tests"
echo "=================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Check if Student Management Service is running
echo -e "${YELLOW}Test 1: Checking Student Management Service...${NC}"
if curl -s "${BASE_URL}/api/students" > /dev/null; then
    echo -e "${GREEN}✓ Student Management Service is running${NC}"
else
    echo -e "${RED}✗ Student Management Service is not accessible${NC}"
    exit 1
fi
echo ""

# Test 2: Check if Finance Service is running
echo -e "${YELLOW}Test 2: Checking Finance Service...${NC}"
if curl -s "${FINANCE_URL}/health" > /dev/null; then
    echo -e "${GREEN}✓ Finance Service is running${NC}"
else
    echo -e "${RED}✗ Finance Service is not accessible${NC}"
    echo -e "${YELLOW}Note: Integration will use graceful degradation${NC}"
fi
echo ""

# Test 3: Test Stipend Calculation (requires Finance Service)
echo -e "${YELLOW}Test 3: Testing Stipend Calculation with Deductions...${NC}"
CALC_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/stipend/calculate" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 1,
    "stipend_type": "scholarship",
    "amount": 5000.00
  }')

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Calculation endpoint is accessible${NC}"
    echo "Response: $CALC_RESPONSE"
else
    echo -e "${RED}✗ Calculation request failed${NC}"
fi
echo ""

# Test 4: Create Stipend Allocation (with auto-calculation)
echo -e "${YELLOW}Test 4: Creating Stipend Allocation (with auto-calculation)...${NC}"
ALLOC_RESPONSE=$(curl -s -X POST "${BASE_URL}/api/stipend/allocations" \
  -H "Content-Type: application/json" \
  -d '{
    "allocation_id": "TEST-'$(date +%s)'",
    "student_id": 1,
    "amount": 5000.00,
    "semester": 1,
    "academic_year": "2025",
    "status": "pending"
  }')

if echo "$ALLOC_RESPONSE" | grep -q "allocation_id"; then
    echo -e "${GREEN}✓ Stipend allocation created successfully${NC}"
    echo "Response: $ALLOC_RESPONSE"
else
    echo -e "${RED}✗ Failed to create stipend allocation${NC}"
    echo "Response: $ALLOC_RESPONSE"
fi
echo ""

# Test 5: Get Student Finance Stipends
echo -e "${YELLOW}Test 5: Retrieving Finance Stipends for Student...${NC}"
STIPENDS_RESPONSE=$(curl -s "${BASE_URL}/api/students/1/finance-stipends")

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Finance stipends endpoint is accessible${NC}"
    echo "Response: $STIPENDS_RESPONSE"
else
    echo -e "${RED}✗ Failed to get finance stipends${NC}"
fi
echo ""

# Test 6: Check Stipend Eligibility
echo -e "${YELLOW}Test 6: Checking Student Eligibility...${NC}"
ELIGIBILITY_RESPONSE=$(curl -s "${BASE_URL}/api/stipend/eligibility/1")

if echo "$ELIGIBILITY_RESPONSE" | grep -q "is_eligible"; then
    echo -e "${GREEN}✓ Eligibility check successful${NC}"
    echo "Response: $ELIGIBILITY_RESPONSE"
else
    echo -e "${RED}✗ Eligibility check failed${NC}"
    echo "Response: $ELIGIBILITY_RESPONSE"
fi
echo ""

echo "=================================="
echo "Test Summary"
echo "=================================="
echo -e "${GREEN}Integration tests completed!${NC}"
echo ""
echo "Notes:"
echo "- If Finance Service is unavailable, the Student Management Service"
echo "  will continue to work with graceful degradation"
echo "- Check the service logs for detailed integration information"
echo "- Ensure FINANCE_GRPC_URL environment variable is set correctly"
echo ""
