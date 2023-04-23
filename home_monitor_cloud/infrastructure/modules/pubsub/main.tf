resource "google_pubsub_topic" "topic" {
  name    = var.topic_name
  project = var.project

  dynamic "schema_settings" {
    for_each = var.schema_settings
    content {
      schema   = schema_settings.value.schema
      encoding = schema_settings.value.encoding
    }
  }
}

resource "google_pubsub_topic" "topic_dead_letter" {
  name    = "${var.topic_name}-deadletter"
  project = var.project
}

resource "google_pubsub_subscription" "topic_subscription" {
  name    = "${var.topic_name}-sub"
  project = var.project
  topic   = google_pubsub_topic.topic.name

  dynamic "bigquery_config" {
    for_each = var.bigquery_config
    content {
      table            = bigquery_config.value.table
      use_topic_schema = bigquery_config.value.use_topic_schema
    }
  }

  ack_deadline_seconds = var.ack_deadline_seconds

  dynamic "retry_policy" {
    for_each = var.retry_policy
    content {
      minimum_backoff = retry_policy.value.minimum_backoff
      maximum_backoff = retry_policy.value.maximum_backoff
    }
  }

  dynamic "push_config" {
    for_each = var.push_config
    content {
      push_endpoint = push_config.value.push_endpoint
    }
  }

  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.topic_dead_letter.id
    max_delivery_attempts = var.max_delivery_attempts
  }

  depends_on = [
    google_pubsub_topic.topic,
    google_pubsub_topic.topic_dead_letter
  ]
}

resource "google_pubsub_subscription" "topic_subscription_dead_letter" {
  name    = "${var.topic_name}-sub-deadletter"
  project = var.project
  topic   = google_pubsub_topic.topic_dead_letter.name

  ack_deadline_seconds = var.ack_deadline_seconds

  dynamic "retry_policy" {
    for_each = var.retry_policy
    content {
      minimum_backoff = retry_policy.value.minimum_backoff
      maximum_backoff = retry_policy.value.maximum_backoff
    }
  }

  depends_on = [
    google_pubsub_topic.topic_dead_letter
  ]
}

resource "google_pubsub_schema" "schema" {
  count      = var.schema_name == "" ? 0 : 1
  name       = var.schema_name
  project    = var.project
  type       = var.schema_type
  definition = var.schema_definition
}
