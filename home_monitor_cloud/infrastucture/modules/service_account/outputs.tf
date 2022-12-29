output "email" {
  value = google_service_account.service_account.email
}

output "private_key" {
  value = google_service_account_key.service_account_key.private_key
}
