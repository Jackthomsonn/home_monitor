[![Deploy Home Monitor Workflow](https://github.com/Jackthomsonn/home_monitor/actions/workflows/deploy-it.yaml/badge.svg)](https://github.com/Jackthomsonn/home_monitor/actions/workflows/deploy-it.yaml)

![system](https://github.com/Jackthomsonn/home_monitor/assets/11717131/a88ae45b-b0df-464b-a3e9-72dc52b3259c)

![Home Monitor](https://github.com/Jackthomsonn/home_monitor/assets/11717131/c9e83460-b6bf-4cf7-bddf-4950d6be533c)

## Technologies / services used

- Elixir for the firmware running on a Raspberry PI 4
- DHT22 sensor for temperature/humidity readings
- TP Link plug for energy monitoring
- Go lang for serverless functions (Hosted in GCP)
- EMQX for the MQTT broker client (Hosted on its own compute within the VPC)
