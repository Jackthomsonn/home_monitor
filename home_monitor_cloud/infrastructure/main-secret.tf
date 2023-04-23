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