resource "google_cloud_scheduler_job" "job" {
  name             = "consumption-ingestion-job"
  project          = var.project
  region           = var.region
  description      = "Gets the latest consumption data from the meter and ingests that data into big query"
  schedule         = "00 08 * * *"
  time_zone        = "Europe/London"
  attempt_deadline = "60s"

  retry_config {
    retry_count = 5
  }

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-${var.project}.cloudfunctions.net/TriggerConsumptionData"
    body = base64encode(jsonencode({
      topic_name : "consumption-ingestion"
    }))
  }

  depends_on = [
    module.project_services
  ]
}

resource "google_cloud_scheduler_job" "ingest_carbon_intensity_job" {
  name             = "carbon-intensity-ingestion-job"
  project          = var.project
  region           = var.region
  description      = "Gets the latest carbon intensity data from the grid and ingests that data into big query"
  schedule         = "*/31 * * * *"
  time_zone        = "Europe/London"
  attempt_deadline = "60s"

  retry_config {
    retry_count = 5
  }

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-${var.project}.cloudfunctions.net/TriggerConsumptionData"
    body = base64encode(jsonencode({
      topic_name : "carbon-intensity-ingestion"
    }))
  }

  depends_on = [
    module.project_services
  ]
}

resource "google_cloud_scheduler_job" "ingest_home_totals_job" {
  name             = "home-totals-ingestion-job"
  project          = var.project
  region           = var.region
  description      = "Gets the latest home totals data from big query and ingests that data into the datastore"
  schedule         = "05 08 * * *"
  time_zone        = "Europe/London"
  attempt_deadline = "60s"

  retry_config {
    retry_count = 5
  }

  http_target {
    http_method = "POST"
    uri         = "https://${var.region}-${var.project}.cloudfunctions.net/TriggerConsumptionData"
    body = base64encode(jsonencode({
      topic_name : "home-totals-ingestion"
    }))
  }

  depends_on = [
    module.project_services
  ]
}
