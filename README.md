![diagram](https://user-images.githubusercontent.com/11717131/210117086-a79049fe-e4d5-47b3-87d0-021cf94efb21.png)

## Technologies / services used

- Elixir for the firmware running on a Raspberry PI 4
- DHT22 sensor for temperature/humidity readings
- Go lang for serverless functions (Hosted in GCP)
- EMQX for the MQTT broker client (Hosted on its own compute within the VPC)
- Terraform IAC
