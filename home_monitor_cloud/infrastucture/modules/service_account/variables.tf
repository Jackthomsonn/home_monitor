variable "service_account_id" {
  type        = string
  description = "The ID of the service account"
}

variable "service_account_display_name" {
  type        = string
  description = "The display name of the service account"
}

variable "service_account_description" {
  type        = string
  description = "The description of the service account"
  default     = "Service Account Description not set (Created via TF)"
}

variable "project" {
  type        = string
  description = "The project to create the service account in"
}

