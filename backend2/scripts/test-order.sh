#!/bin/sh

API_URL="${API_URL:-http://localhost:3000}"
SHOP_ID="${SHOP_ID:-1}"
NUMBER="${NUMBER:-TEST-001}"
TOTAL="${TOTAL:-1000}"
CUSTOMER="${CUSTOMER:-Test}"

echo "Creating test order:"
echo "  API_URL: $API_URL"
echo "  SHOP_ID: $SHOP_ID"
echo "  NUMBER: $NUMBER"
echo "  TOTAL: $TOTAL"
echo "  CUSTOMER: $CUSTOMER"
echo ""

curl -X POST "${API_URL}/shops/${SHOP_ID}/orders" \
	-H "Content-Type: application/json" \
	-d "{\"number\":\"${NUMBER}\",\"total\":${TOTAL},\"customerName\":\"${CUSTOMER}\"}"

echo ""
