![diagram](https://user-images.githubusercontent.com/11717131/210117086-a79049fe-e4d5-47b3-87d0-021cf94efb21.png)


![Home Monitor Diagram](https://user-images.githubusercontent.com/11717131/222989343-6fd65048-555e-48a3-893b-fb8b762a5a21.jpeg)

## Technologies / services used

- Elixir for the firmware running on a Raspberry PI 4
- DHT22 sensor for temperature/humidity readings
- TP Link plug for energy monitoring
- Go lang for serverless functions (Hosted in GCP)
- EMQX for the MQTT broker client (Hosted on its own compute within the VPC)
