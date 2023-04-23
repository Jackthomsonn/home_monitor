module "bigquery" {
  source       = "./modules/bigquery"
  project      = var.project
  dataset_name = "home_monitor_dataset"
  tables = [{
    name                = "home_monitor_table"
    deletion_protection = true
    schema              = <<EOF
    [
      {
        "name": "temperature",
        "type": "FLOAT",
        "description": "Temperature in degrees celsius"
      },
      {
        "name": "client_id",
        "type": "STRING",
        "description": "Client ID"
      },
      {
        "name": "timestamp",
        "type": "TIMESTAMP",
        "description": "The timestamp at which the temperature was recorded on the device"
      }
    ]
    EOF
    }, {
    name                = "home_monitor_consumption_table"
    deletion_protection = true
    schema              = <<EOF
    [
      {
        "name": "timestamp",
        "type": "TIMESTAMP",
        "description": "The timestamp at which the consumption data was recorded"
      },
      {
        "name": "value",
        "type": "FLOAT",
        "description": "The value of the consumption in kWh"
      }
    ]
    EOF
    }, {
    name                = "home_monitor_carbon_intensity"
    deletion_protection = true
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
