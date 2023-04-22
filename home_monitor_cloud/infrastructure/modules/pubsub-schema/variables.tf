variable "project" {
  type        = string
  description = "The project ID to deploy to"
}

variable "schema_name" {
  type        = string
  description = "The name of the schema to create"
}

variable "schema_type" {
  type        = string
  description = "The type of the schema to create"
}

variable "schema_definition" {
  type        = string
  description = "The definition of the schema to create"
}
