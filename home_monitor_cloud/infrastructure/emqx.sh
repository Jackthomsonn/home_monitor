response=$(curl -u '<username>:<password>' \
     -X 'POST' '<url>/api/v5/api_key' \
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

curl -X POST '<url>/api/v5/bridges' \
     -u ${api_key}:${api_secret} \
     -H "Content-Type: application/json" \
     -d '{
 "body": "{\"payload\":\"{current_ma: ${payload.current_ma}, power_mw: ${payload.power_mw}, total_wh: ${payload.total_wh}, voltage_mv: ${payload.voltage_mv}, ip: ${payload.ip}, alias: ${payload.alias} }\",\"topic\":\"${topic}\",\"timestamp\":\"${timestamp}\",\"client_id\":\"${client_id}\",\"type\":\"energy\"}",
 "connect_timeout": "10s",
  "enable": true,
  "enable_pipelining": 100,
  "max_retries": 5,
  "method": "post",
  "name": "energy_data_publisher",
  "pool_size": 4,
  "pool_type": "random",
  "request_timeout": "10s",
  "ssl": {
    "enable": false
  },
  "type": "webhook",
  "url": "<url>"
}'


curl -X POST '<url>/api/v5/rules' \
     -u ${api_key}:${api_secret} \
     -H "Content-Type: application/json" \
     -d '{
  "name": "publish_energy_data",
  "sql": "SELECT payload.current_ma,
       payload.power_mw,
       payload.total_wh,
       payload.voltage_mv,
       payload.ip,
       payload.alias,
       topic,
       clientid as client_id,
       payload.timestamp as timestamp
       FROM \"reports/+/energy\"",
  "actions": [
    "webhook:energy_data_publisher"
  ],
  "enable": true,
  "metadata": {}
}'
