variable "network_name" {
  type        = string
  description = "The name of the network to attach this firewall to."
}

variable "project" {
  type        = string
  description = "The ID of the project in which the resource belongs."
}

variable "secrets" {
  type = object({
    data = map(string)
  })
}
