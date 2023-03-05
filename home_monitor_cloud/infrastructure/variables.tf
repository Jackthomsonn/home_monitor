variable "project" {
  type    = string
  default = "home-monitor-373013"
}

variable "region" {
  type        = string
  default     = "europe-west1"
  description = "Region to deploy to"
}

variable "zone" {
  type        = string
  default     = "europe-west1-c"
  description = "Zone to deploy to"
}
