resource "google_service_account" "service_account" {
  account_id   = var.service_account_id
  display_name = var.service_account_display_name
  description  = var.service_account_description
  project      = var.project
}

resource "google_service_account_key" "service_account_key" {
  service_account_id = google_service_account.service_account.name
  public_key_type    = "TYPE_X509_PEM_FILE"
  depends_on = [
    google_service_account.service_account
  ]
}
