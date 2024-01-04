module "bigquery" {
  source       = "./modules/bigquery"
  project      = var.project
  dataset_name = "home_monitor_dataset"
  tables = [{
    name                = "home_monitor_table"
    deletion_protection = true
    clustering          = ["client_id"]
    schema              = <<EOF
    [
      {
        "name": "payload",
        "type": "STRING",
        "description": "Raw rayload data"
      },
      {
        "name": "topic",
        "type": "STRING",
        "description": "The topic the message originated from"
      },
      {
        "name": "type",
        "type": "STRING",
        "description": "The type of data this message represents"
      },
      {
        "name": "client_id",
        "type": "STRING",
        "description": "Client ID"
      },
      {
        "name": "timestamp",
        "type": "TIMESTAMP",
        "description": "The timestamp at which the data was recorded on the device"
      }
    ]
    EOF
    }, {
    name                = "home_monitor_carbon_intensity"
    deletion_protection = true
    clustering          = []
    schema              = <<EOF
    [
      {
        "name": "timestamp",
        "type": "TIMESTAMP",
        "description": "The timestamp at which the consumption data was recorded"
      },
      {
        "name": "actual",
        "type": "FLOAT",
        "description": "The actual carbon intensity in gCO2/kWh"
      },
      {
        "name": "forecast",
        "type": "FLOAT",
        "description": "The forecasted carbon intensity in gCO2/kWh"
      }
    ]
    EOF
  }]

  depends_on = [
    module.project_services
  ]
}
