output "emqxx_service_account_key" {
  value     = module.emqxx_service_account.private_key
  sensitive = true
}
