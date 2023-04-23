terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.47.0"
    }
    sops = {
      source  = "carlpett/sops"
      version = "~> 0.5"
    }
  }

  backend "gcs" {
    bucket = "home-monitor-terraform-state"
    prefix = "terraform/state"
  }
}

data "sops_file" "secrets" {
  source_file = "../secrets/secrets.yaml"
}

data "google_project" "project" {
  project_id = var.project
}

provider "google-beta" {
  project = var.project
  region  = var.region
  zone    = var.zone
}

module "project_services" {
  source     = "terraform-google-modules/project-factory/google//modules/project_services"
  version    = "3.3.0"
  project_id = var.project

  activate_apis = [
    "run.googleapis.com",
    "cloudkms.googleapis.com",
    "servicenetworking.googleapis.com",
    "sqladmin.googleapis.com",
    "compute.googleapis.com",
    "cloudfunctions.googleapis.com",
    "cloudbuild.googleapis.com",
    "artifactregistry.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudscheduler.googleapis.com",
    "datastore.googleapis.com"
  ]

  disable_services_on_destroy = true
  disable_dependent_services  = true
}

##### Pub sub
module "state-pubsub-schema" {
  source      = "./modules/pubsub-schema"
  project     = var.project
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
}

module "state-pubsub" {
  source     = "./modules/pubsub"
  project    = var.project
  topic_name = "state"
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

##### Networking
resource "google_compute_network" "vpc_network" {
  project                 = var.project
  name                    = "home-monitor-vpc"
  auto_create_subnetworks = true
  mtu                     = 1460
  depends_on = [
    module.project_services
  ]
}

##### Big Query
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

##### IAM
module "iam" {
  source = "./modules/iam"

  project_number = data.google_project.project.number
  project        = var.project
}

##### Firewall
module "firewall" {
  source = "./modules/firewall"

  project      = var.project
  network_name = google_compute_network.vpc_network.name
  secrets      = data.sops_file.secrets
}

##### Compute Engine (EMQX)
module "emqx" {
  source = "./modules/emqx"

  project      = var.project
  zone         = var.zone
  region       = var.region
  network_name = google_compute_network.vpc_network.name
  secrets      = data.sops_file.secrets

  depends_on = [
    module.project_services,
    module.iam,
    module.firewall
  ]
}

#### Datastore
module "datastore" {
  source  = "terraform-google-modules/cloud-datastore/google"
  project = var.project
  indexes = file("index.yaml")
  depends_on = [
    module.project_services
  ]
}

// TODO - Refactor all into modules and then adopt terragrunt

##### Service Accounts
resource "google_service_account" "emqx_service_account" {
  account_id   = "emqx-service-account"
  display_name = "EMQX Service Account"
  description  = "Service account for EMQX"
  project      = var.project
}

resource "google_service_account_key" "emqx_service_account_key" {
  service_account_id = google_service_account.emqx_service_account.name
  public_key_type    = "TYPE_X509_PEM_FILE"
  depends_on = [
    google_service_account.emqx_service_account
  ]
}

resource "google_service_account" "ci_cd_service_account" {
  account_id   = "cicd-service-account"
  display_name = "CI-CD Service Account"
  description  = "CI-CD service account"
  project      = var.project
}

resource "google_project_iam_member" "ci_cd_service_account_roles" {
  for_each = toset([
    "roles/iam.serviceAccountUser",
    "roles/cloudfunctions.developer",
    "roles/storage.objectAdmin",
    "roles/cloudkms.cryptoKeyDecrypter",
    "roles/compute.networkUser",
    "roles/secretmanager.secretAccessor",
    "roles/owner",
    "roles/cloudscheduler.admin",
    "roles/pubsub.admin"
  ])
  role    = each.key
  member  = "serviceAccount:${google_service_account.ci_cd_service_account.email}"
  project = var.project
}

resource "google_service_account_key" "ci_cd_service_account_key" {
  service_account_id = google_service_account.ci_cd_service_account.name
  public_key_type    = "TYPE_X509_PEM_FILE"
  depends_on = [
    google_service_account.ci_cd_service_account
  ]
}

##### KMS
resource "google_kms_key_ring" "home_monitor_keyring" {
  name     = "home-monitor-keyring"
  project  = var.project
  location = "global"
}

resource "google_kms_crypto_key" "home_monitor_key" {
  name     = "home-monitor-key"
  key_ring = google_kms_key_ring.home_monitor_keyring.id

  lifecycle {
    prevent_destroy = true
  }
}

##### Secrets
resource "google_secret_manager_secret" "consumption_secret" {
  secret_id = "consumption_secret"
  project   = var.project

  replication {
    automatic = true
  }

  depends_on = [
    module.project_services
  ]
}

resource "google_secret_manager_secret_version" "consumption_secret_version" {
  secret      = google_secret_manager_secret.consumption_secret.id
  secret_data = data.sops_file.secrets.data["consumption_secret"]

  depends_on = [
    google_secret_manager_secret.consumption_secret
  ]
}

resource "google_secret_manager_secret" "redis_connection_string" {
  secret_id = "redis_connection_string"
  project   = var.project

  replication {
    automatic = true
  }

  depends_on = [
    module.project_services
  ]
}
resource "google_secret_manager_secret_version" "redis_connection_string_version" {
  secret      = google_secret_manager_secret.redis_connection_string.id
  secret_data = data.sops_file.secrets.data["redis_connection_string"]

  depends_on = [
    google_secret_manager_secret.redis_connection_string
  ]
}

##### Cloud scheduler jobs
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
