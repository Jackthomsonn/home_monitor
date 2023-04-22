resource "google_pubsub_schema" "schema" {
  name       = var.schema_name
  project    = var.project
  type       = var.schema_type
  definition = var.schema_definition
}
