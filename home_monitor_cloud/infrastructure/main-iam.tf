resource "google_project_service" "iam" {
  project = var.project
  service = "iam.googleapis.com"
}

# Assign the service account the pubsub publisher role (global)
resource "google_project_iam_member" "pubsub_publisher_iam" {
  project = var.project
  role    = "roles/pubsub.publisher"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"

  depends_on = [
    google_project_service.iam
  ]
}

# Assign the service account the pubsub subscriber role (global)
resource "google_project_iam_member" "pubsub_subscriber_iam" {
  project = var.project
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"

  depends_on = [
    google_project_service.iam
  ]
}

# Assign the service account the biq query viewer role (global)
resource "google_project_iam_member" "viewer" {
  project = var.project
  role    = "roles/bigquery.metadataViewer"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"

  depends_on = [
    google_project_service.iam
  ]
}

# Assign the service account the biq query editor role (global)
resource "google_project_iam_member" "editor" {
  project = var.project
  role    = "roles/bigquery.dataEditor"
  member  = "serviceAccount:service-${data.google_project.project.number}@gcp-sa-pubsub.iam.gserviceaccount.com"

  depends_on = [
    google_project_service.iam
  ]
}

###############################################################################################

##### Ingest data service account
resource "google_service_account" "ingest_data_iam_service_account" {
  account_id   = "ingest-data-iam-sa"
  project      = var.project
  display_name = "Ingest Data IAM Service Account used for the Ingest data cloud function"
}

resource "google_project_iam_member" "ingest_data_iam_service_account_member_roles" {
  project = var.project
  for_each = toset([
    "roles/bigquery.dataEditor",
    "roles/secretmanager.secretAccessor"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.ingest_data_iam_service_account.email}"
}

###############################################################################################

#### Get totals for home service account
resource "google_service_account" "get_totals_for_home_service_account" {
  account_id   = "get-totals-for-home-iam-sa"
  project      = var.project
  display_name = "Get Totals For Home Service Account used for the Get Totals For Home function"
}

resource "google_project_iam_member" "get_totals_for_home_service_account_member_roles" {
  project = var.project
  for_each = toset([
    "roles/datastore.viewer",
    "roles/secretmanager.secretAccessor"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.get_totals_for_home_service_account.email}"
}

###############################################################################################

#### Ingest home totals service account
resource "google_service_account" "ingest_home_totals_service_account" {
  account_id   = "ingest-home-totals-iam-sa"
  project      = var.project
  display_name = "Ingest Home Totals Service Account used for the Ingest Home Totals function"
}

resource "google_project_iam_member" "ingest_home_totals_service_account_member_roles" {
  project = var.project
  for_each = toset([
    "roles/bigquery.dataViewer",
    "roles/secretmanager.secretAccessor",
    "roles/bigquery.jobUser",
    "roles/datastore.user"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.ingest_home_totals_service_account.email}"
}

###############################################################################################

#### Discover Devices service account
resource "google_service_account" "discover_devices_service_account" {
  account_id   = "discover-devices-iam-sa"
  project      = var.project
  display_name = "Discover Devices Service Account used for the Discover Devices function"
}

resource "google_project_iam_member" "discover_devices_service_account_member_roles" {
  project = var.project
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/datastore.user"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.discover_devices_service_account.email}"
}

###############################################################################################

#### Get Devices service account
resource "google_service_account" "get_devices_service_account" {
  account_id   = "get-devices-iam-sa"
  project      = var.project
  display_name = "Get Devices Service Account used for the Get Devices function"
}

resource "google_project_iam_member" "get_devices_service_account_member_roles" {
  project = var.project
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/datastore.user"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.get_devices_service_account.email}"
}

###############################################################################################

#### Send Command service account
resource "google_service_account" "send_command_service_account" {
  account_id   = "send-command-iam-sa"
  project      = var.project
  display_name = "Send Command Service Account used for the Send Command function"
}

resource "google_project_iam_member" "send_command_service_account_member_roles" {
  project = var.project
  for_each = toset([
    "roles/secretmanager.secretAccessor",
    "roles/datastore.user"
  ])
  role   = each.key
  member = "serviceAccount:${google_service_account.send_command_service_account.email}"
}
