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
