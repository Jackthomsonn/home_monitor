output "emqx_service_account_key" {
  value     = google_service_account_key.emqx_service_account_key.private_key
  sensitive = true
}

output "cicd_service_account_key" {
  value     = google_service_account_key.ci_cd_service_account_key.private_key
  sensitive = true
}
