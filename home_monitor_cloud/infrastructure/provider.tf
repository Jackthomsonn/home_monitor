terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.47.0"
    }
    sops = {
      source  = "carlpett/sops"
      version = "~> 0.5"
    }
  }

  backend "gcs" {
    bucket = "home-monitor-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google-beta" {
  project = var.project
  region  = var.region
  zone    = var.zone
}
