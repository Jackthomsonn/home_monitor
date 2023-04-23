variable "topic_name" {
  type        = string
  description = "The name of the topic to create"
}

variable "project" {
  type        = string
  description = "The project to create the topic in"
}

variable "schema_settings" {
  type = list(object({
    schema   = string
    encoding = string
  }))
  description = "The schema settings for the topic"
  default     = []
}

variable "bigquery_config" {
  type = list(object({
    table            = string
    use_topic_schema = bool
  }))
  description = "The bigquery config for the topic"
  default     = []
}

variable "ack_deadline_seconds" {
  type        = number
  description = "The ack deadline for the topic"
  default     = 20
}

variable "retry_policy" {
  type = list(object({
    minimum_backoff = string
    maximum_backoff = string
  }))
  description = "The retry policy for the topic"
  default     = []
}

variable "max_delivery_attempts" {
  type        = number
  description = "The max delivery attempts for the topic"
  default     = 5
}

variable "depends_on_config" {
  type        = list(any)
  description = "The depends on for the topic"
  default     = []
}

variable "push_config" {
  type = list(object({
    push_endpoint = string
  }))
  description = "The push config for the topic"
  default     = []
}

variable "schema_name" {
  type        = string
  description = "The name of the schema to create"
  default     = ""
}

variable "schema_type" {
  type        = string
  description = "The type of the schema to create"
  default     = ""
}

variable "schema_definition" {
  type        = string
  description = "The definition of the schema to create"
  default     = ""
}

