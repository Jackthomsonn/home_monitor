resource "google_service_account" "emqx_instance_service_account" {
  account_id   = "emqx-instance-service-account"
  project      = var.project
  display_name = "EMQX Instance Service Account"
}

resource "google_compute_address" "emqx_static_ip" {
  name    = "ipv4-address"
  project = var.project
  region  = var.region
}

resource "google_compute_instance" "emqx_instance" {
  name                      = "emqx-instance"
  project                   = var.project
  zone                      = var.zone
  description               = "The EMQX instance"
  machine_type              = "e2-small"
  allow_stopping_for_update = true

  lifecycle {
    ignore_changes = [
      labels,
      resource_policies,
    ]
  }

  can_ip_forward = false

  tags = ["http-server"]

  metadata_startup_script = <<-EOF
  docker pull emqx/emqx:5.0.9 &&
  docker run -d --name emqx -p 1883:1883 -p 8083:8083 -p 8084:8084 -p 8883:8883 -p 18083:18083 emqx/emqx:5.0.9
  EOF


  metadata = {
    ssh-keys = "user:${var.secrets.data["ssh_key"]}"
  }

  scheduling {
    automatic_restart   = true
    on_host_maintenance = "MIGRATE"
  }

  boot_disk {
    initialize_params {
      image = "projects/cos-cloud/global/images/cos-stable-101-17162-40-42"
    }
  }

  network_interface {
    network = var.network_name
    access_config {
      nat_ip = google_compute_address.emqx_static_ip.address
    }
  }

  service_account {
    email  = google_service_account.emqx_instance_service_account.email
    scopes = ["cloud-platform"]
  }
}
