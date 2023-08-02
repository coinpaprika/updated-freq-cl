#!/bin/bash
API_KEY="api_your_key_here"

last_price=""
echo "ticker_id,now,now_timestamp,price,price_changed"
while true; do
  price=$(curl -s https://api-pro.coinpaprika.com/v1/tickers -H "Authorization: $API_KEY" | jq '.[] | select(.id == "btc-bitcoin") | .quotes.USD.price')
  now=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  now_timestamp=$(date +%s)
  changed=false
  if [ "$price" != "$last_price" ]; then
    if [ "$last_price" ]; then
      changed=true
    fi
    last_price=$price
  fi
  echo "'btc-bitcoin',$now,$now_timestamp,$price,$changed"
  sleep 10
done