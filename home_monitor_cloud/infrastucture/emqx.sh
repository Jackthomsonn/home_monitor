response=$(curl -u 'admin:<password>' \
     -X 'POST' '<host>/api/v5/api_key' \
     -H 'accept: application/json' \
     -H 'Content-Type: application/json' \
     -d '{
            "name": "EMQX-SETUP-KEY",
            "expired_at": "2099-01-01T00:00:00.000Z",
            "desc": "The key used to setup the EMQX broker",
            "enable": true
        }')

api_key=$(echo $response | jq -r '.api_key')
api_secret=$(echo $response | jq -r '.api_secret')

curl -X POST '<host>/api/v5/bridges' \
     -u ${api_key}:${api_secret} \
     -H "Content-Type: application/json" \
     -d '{
 "body": "{\"temperature\":${temperature},\"client_id\":\"${clientid}\",\"timestamp\":\"${timestamp}\"}",
 "connect_timeout": "10s",
  "enable": true,
  "enable_pipelining": 100,
  "max_retries": 5,
  "method": "post",
  "name": "data_publisher",
  "pool_size": 4,
  "pool_type": "random",
  "request_timeout": "10s",
  "ssl": {
    "enable": false
  },
  "type": "webhook",
  "url": "<url>"
}'


curl -X POST '<host>/api/v5/rules' \
     -u ${api_key}:${api_secret} \
     -H "Content-Type: application/json" \
     -d '{
  "name": "state_pub",
  "sql": "SELECT payload.temperature as temperature,
       payload.timestamp as timestamp,
       clientid as clientid

FROM \"reports/#\"",
  "actions": [
    "webhook:data_publisher"
  ],
  "enable": true,
  "metadata": {}
}'
