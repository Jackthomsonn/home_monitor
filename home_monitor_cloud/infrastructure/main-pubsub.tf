module "state-pubsub" {
  source      = "./modules/pubsub"
  project     = var.project
  topic_name  = "state"
  schema_name = "house-monitor-schema"
  schema_type = "AVRO"
  schema_definition = jsonencode(
    {
      fields = [
        {
          name = "temperature"
          type = "float"
        },
        {
          name = "client_id"
          type = "string"
        },
        {
          name = "timestamp"
          type = "string"
        },
      ]
      name = "HouseMonitor"
      type = "record"
    }
  )
  schema_settings = [{
    encoding = "JSON"
    schema   = "projects/${var.project}/schemas/house-monitor-schema"
  }]

  bigquery_config = [{
    table            = "${var.project}:${module.bigquery.dataset_name}.home_monitor_table"
    use_topic_schema = true
  }]

  depends_on_config = [
    module.project_services
  ]
}

module "consumption-ingestion-pubsub" {
  source     = "./modules/pubsub"
  project    = var.project
  topic_name = "consumption-ingestion"

  retry_policy = [{
    minimum_backoff = "300s"
    maximum_backoff = "300s"
  }]

  push_config = [{
    push_endpoint = "https://${var.region}-${var.project}.cloudfunctions.net/IngestConsumptionData"
  }]

  depends_on_config = [
    module.project_services
  ]
}

module "carbon-intensity-ingestion-pubsub" {
  source     = "./modules/pubsub"
  project    = var.project
  topic_name = "carbon-intensity-ingestion"

  retry_policy = [{
    minimum_backoff = "10s"
    maximum_backoff = "20s"
  }]
  max_delivery_attempts = 10

  push_config = [{
    push_endpoint = "https://${var.region}-${var.project}.cloudfunctions.net/IngestCarbonIntensityData"
  }]

  depends_on_config = [
    module.project_services
  ]
}

module "home-totals-ingestion-pubsub" {
  source     = "./modules/pubsub"
  project    = var.project
  topic_name = "home-totals-ingestion"

  retry_policy = [{
    minimum_backoff = "10s"
    maximum_backoff = "20s"
  }]
  max_delivery_attempts = 10

  push_config = [{
    push_endpoint = "https://${var.region}-${var.project}.cloudfunctions.net/IngestHomeTotals"
  }]

  depends_on_config = [
    module.project_services
  ]
}