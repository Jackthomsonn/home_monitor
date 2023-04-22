resource "google_compute_firewall" "ssh_rule" {
  project     = var.project
  name        = "allow-ssh"
  network     = var.network_name
  description = "Allow SSH inbound traffic from specific IP"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.secrets.data["ip_address"]]
}

resource "google_compute_firewall" "emqx_dashboard_rule" {
  project     = var.project
  name        = "emqx-dashboard"
  network     = var.network_name
  description = "Allow ingress traffic to emqx dashboard from specific IP"

  allow {
    protocol = "tcp"
    ports    = ["18083"]
  }

  source_ranges = [var.secrets.data["ip_address"]]
}

resource "google_compute_firewall" "emqx_tcp" {
  project     = var.project
  name        = "emqx-tcp"
  network     = var.network_name
  description = "Allow ingress traffic to emqx tcp (IOT Devices)"

  allow {
    protocol = "tcp"
    ports    = ["1883"]
  }

  source_ranges = ["0.0.0.0/0"]
}
