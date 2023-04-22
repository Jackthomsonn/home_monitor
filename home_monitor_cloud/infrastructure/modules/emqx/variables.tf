variable "project" {
  type        = string
  description = "The project ID to deploy to"
}

variable "region" {
  type        = string
  description = "The region to deploy to"
}

variable "zone" {
  type        = string
  description = "The zone to deploy to"
}

variable "secrets" {
  type = object({
    data = map(string)
  })
}

variable "network_name" {
  type        = string
  description = "The name of the network to attach this firewall to."
}
