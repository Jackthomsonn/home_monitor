terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.47.0"
    }
  }

  backend "gcs" {
    bucket = "home-monitor-terraform-state"
    prefix = "terraform/state"
  }
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
    "artifactregistry.googleapis.com"
  ]

  disable_services_on_destroy = true
  disable_dependent_services  = true
}

data "google_project" "project" {
  project_id = var.project
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

resource "google_compute_network_peering" "emqxx_peering" {
  name         = "emqxx-peering"
  network      = google_compute_network.vpc_network.self_link
  peer_network = "projects/emq-x-cloud-324802/global/networks/e4bf24df"
  depends_on = [
    google_compute_network.vpc_network
  ]
}

resource "google_compute_network_firewall_policy" "home_monitor_network_firewall_policy" {
  name        = "home-monitor-firewall-policy"
  project     = var.project
  description = "This is a simple firewall policy for the home-monitor project"
}

resource "google_compute_network_firewall_policy_rule" "primary" {
  firewall_policy = google_compute_network_firewall_policy.home_monitor_network_firewall_policy.name
  action          = "allow"
  direction       = "INGRESS"
  priority        = 1000
  rule_name       = "emqxx rule"
  description     = "Allow traffic from emqxx"
  project         = var.project
  match {
    src_ip_ranges = ["10.25.27.0/24"]
    layer4_configs {
      ip_protocol = "all"
    }
  }
  disabled = false
  depends_on = [
    google_compute_network_firewall_policy.home_monitor_network_firewall_policy
  ]
}

resource "google_compute_network_firewall_policy_association" "primary" {
  name              = "home-monitor-firewall-policy to VPC network"
  attachment_target = google_compute_network.vpc_network.id
  firewall_policy   = google_compute_network_firewall_policy.home_monitor_network_firewall_policy.name
  project           = var.project
  depends_on = [
    google_compute_network_firewall_policy.home_monitor_network_firewall_policy
  ]
}

##### Big Query
resource "google_bigquery_dataset" "home_monitor" {
  dataset_id = "home_monitor_dataset"
  project    = var.project
}

resource "google_bigquery_table" "home_monitor" {
  deletion_protection = false
  project             = var.project
  table_id            = "home_monitor_table"
  dataset_id          = google_bigquery_dataset.home_monitor.dataset_id

  schema = <<EOF
    [
      {
        "name": "temperature",
        "type": "FLOAT",
        "description": "Temperature in degrees celsius"
      }
    ]
    EOF

  depends_on = [
    google_bigquery_dataset.home_monitor
  ]
}

##### Pub sub
resource "google_pubsub_schema" "house_monitor_schema" {
  name       = "house-monitor-schema"
  project    = var.project
  type       = "AVRO"
  definition = "{\r\n  \"type\" : \"record\",\r\n  \"name\" : \"HouseMonitor\",\r\n  \"fields\" : [\r\n    {\r\n      \"name\" : \"temperature\",\r\n      \"type\" : \"float\"\r\n    }\r\n  ]\r\n}"
}

resource "google_pubsub_topic" "topic" {
  name    = "state"
  project = var.project
  schema_settings {
    encoding = "JSON"
    schema   = "projects/${var.project}/schemas/house-monitor-schema"
  }

  depends_on = [
    google_pubsub_schema.house_monitor_schema
  ]
}

resource "google_pubsub_topic" "topic_dead_letter" {
  name    = "state-deadletter"
  project = var.project
}

resource "google_pubsub_subscription" "topic_sub" {
  name    = "state-sub"
  project = var.project
  topic   = google_pubsub_topic.topic.name

  bigquery_config {
    table            = "${google_bigquery_table.home_monitor.project}:${google_bigquery_table.home_monitor.dataset_id}.${google_bigquery_table.home_monitor.table_id}"
    use_topic_schema = true
  }

  ack_deadline_seconds = 20

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.topic_dead_letter.id
    max_delivery_attempts = 5
  }

  depends_on = [
    google_pubsub_topic.topic,
    google_bigquery_table.home_monitor
  ]
}

resource "google_pubsub_subscription" "topic_dead_letter_sub" {
  name    = "state-sub-deadletter"
  project = var.project
  topic   = google_pubsub_topic.topic_dead_letter.name

  ack_deadline_seconds = 600

  depends_on = [
    google_pubsub_topic.topic_dead_letter
  ]
}

##### Service Accounts
module "emqxx_service_account" {
  source                       = "./modules/service_account"
  project                      = var.project
  service_account_id           = "emqxx-service-account"
  service_account_display_name = "EMQXX Service Account"
  depends_on = [
    module.project_services
  ]
}

##### IAM bindings

# Assign the service account the pubsub publisher role (global)
resource "google_project_iam_member" "pubsub_publisher_iam" {
  project = var.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  depends_on = [
    module.project_services
  ]
}

# Assign the service account the pubsub subscriber role (global)
resource "google_project_iam_member" "pubsub_subscriber_iam" {
  project = var.project
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
  depends_on = [
    module.project_services
  ]
}

# Assign the service account the biq query viewer role (global)
resource "google_project_iam_member" "viewer" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.metadataViewer"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

# Assign the service account the biq query editor role (global)
resource "google_project_iam_member" "editor" {
  project = data.google_project.project.project_id
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"
}

# Assign the EMQXX service account the pubsub publisher role
resource "google_project_iam_member" "emqxx_pubsub_publisher_iam" {
  project = var.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${module.emqxx_service_account.email}"
  depends_on = [
    module.project_services
  ]
}
