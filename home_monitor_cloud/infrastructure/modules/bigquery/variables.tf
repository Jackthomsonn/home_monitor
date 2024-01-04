variable "tables" {
  description = "List of tables to create"
  type = list(object({
    name                = string
    schema              = string
    deletion_protection = bool
    clustering          = list(string)
  }))
}

variable "dataset_name" {
  description = "Name of the dataset to create"
  type        = string
}

variable "project" {
  description = "Project to create the dataset in"
  type        = string
}
