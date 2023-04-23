resource "google_kms_key_ring" "home_monitor_keyring" {
  name     = "home-monitor-keyring"
  project  = var.project
  location = "global"
}

resource "google_kms_crypto_key" "home_monitor_key" {
  name     = "home-monitor-key"
  key_ring = google_kms_key_ring.home_monitor_keyring.id

  lifecycle {
    prevent_destroy = true
  }
}
