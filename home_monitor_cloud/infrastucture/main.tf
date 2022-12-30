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
      },
      {
        "name": "client_id",
        "type": "STRING",
        "description": "Client ID"
      },
      {
        "name": "timestamp",
        "type": "STRING",
        "description": "The timestamp at which the temperature was recorded on the device"
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
  definition = "{\r\n  \"type\" : \"record\",\r\n  \"name\" : \"HouseMonitor\",\r\n  \"fields\" : [\r\n    {\r\n      \"name\" : \"temperature\",\r\n      \"type\" : \"float\"\r\n    },\r\n    {\r\n      \"name\" : \"client_id\",\r\n      \"type\" : \"string\"\r\n    },\r\n    {\r\n      \"name\" : \"timestamp\",\r\n      \"type\" : \"string\"\r\n    }\r\n  ]\r\n}"
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

# Assign the EMQX service account the pubsub publisher role
resource "google_project_iam_member" "emqx_pubsub_publisher_iam" {
  project = var.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:${google_service_account.emqx_service_account.email}"
  depends_on = [
    module.project_services
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

##### Firewall
resource "google_compute_firewall" "ssh_rule" {
  project     = var.project
  name        = "allow-ssh"
  network     = google_compute_network.vpc_network.name
  description = "Allow SSH inbound traffic from specific IP"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [data.sops_file.secrets.data["ip_address"]]
}

resource "google_compute_firewall" "emqx_dashboard_rule" {
  project     = var.project
  name        = "emqx-dashboard"
  network     = google_compute_network.vpc_network.name
  description = "Allow ingress traffic to emqx dashboard from specific IP"

  allow {
    protocol = "tcp"
    ports    = ["18083"]
  }

  source_ranges = [data.sops_file.secrets.data["ip_address"]]
}

resource "google_compute_firewall" "emqx_tcp" {
  project     = var.project
  name        = "emqx-tcp"
  network     = google_compute_network.vpc_network.name
  description = "Allow ingress traffic to emqx tcp (IOT Devices)"

  allow {
    protocol = "tcp"
    ports    = ["1883"]
  }

  source_ranges = ["0.0.0.0/0"]
}

##### Compute Engine (EMQX)
resource "google_service_account" "emqx_instance_service_account" {
  account_id   = "emqx-instance-service-account"
  project      = var.project
  display_name = "EMQX Instance Service Account"
}

resource "google_compute_instance" "emqx_instance" {
  name         = "emqx-instance"
  project      = var.project
  zone         = var.zone
  description  = "The EMQX instance"
  machine_type = "e2-medium"

  lifecycle {
    ignore_changes = [
      labels,
      resource_policies,
    ]
  }

  can_ip_forward = false

  tags = ["http-server"]

  metadata_startup_script = <<-EOF
  docker pull emqx/emqx:5.0.9 &&
  docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx:5.0.9
  EOF


  metadata = {
    ssh-keys = "user:${data.sops_file.secrets.data["ssh_key"]}"
  }

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
  }

  boot_disk {
    initialize_params {
      image = "projects/cos-cloud/global/images/cos-stable-101-17162-40-42"
    }
  }

  network_interface {
    network = google_compute_network.vpc_network.name
    access_config {}
  }

  service_account {
    email  = google_service_account.emqx_instance_service_account.email
    scopes = ["cloud-platform"]
  }
}

