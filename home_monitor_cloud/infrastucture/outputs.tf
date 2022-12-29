output "emqx_service_account_key" {
  value     = google_service_account_key.emqx_service_account_key.private_key
  sensitive = true
}
