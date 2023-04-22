resource "google_bigquery_dataset" "big_query_dataset" {
  dataset_id = var.dataset_name
  project    = var.project
}

resource "google_bigquery_table" "big_query_table" {
  for_each = { for key, value in var.tables : key => value }

  deletion_protection = each.value.deletion_protection
  project             = var.project
  table_id            = each.value.name
  dataset_id          = google_bigquery_dataset.big_query_dataset.dataset_id

  schema = each.value.schema

  depends_on = [
    google_bigquery_dataset.big_query_dataset
  ]
}
