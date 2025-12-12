#!/bin/sh

API_URL="${API_URL:-http://localhost:8080}"
SHOP_ID="${SHOP_ID:-1}"
NUMBER="${NUMBER:-TEST-001}"
TOTAL="${TOTAL:-1000}"
CUSTOMER="${CUSTOMER:-Test}"

curl -X POST "${API_URL}/shops/${SHOP_ID}/orders" \
	-H "Content-Type: application/json" \
	-d "{\"number\":\"${NUMBER}\",\"total\":${TOTAL},\"customerName\":\"${CUSTOMER}\"}"

