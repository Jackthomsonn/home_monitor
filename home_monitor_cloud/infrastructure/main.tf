data "sops_file" "secrets" {
  source_file = "../secrets/secrets.yaml"
}

data "google_project" "project" {
  project_id = var.project
}

module "project_services" {
  source     = "terraform-google-modules/project-factory/google//modules/project_services"
  version    = "3.3.0"
  project_id = var.project

  activate_apis = [
    "run.googleapis.com",
    "cloudkms.googleapis.com",
    "servicenetworking.googleapis.com",
    "sqladmin.googleapis.com",
    "compute.googleapis.com",
    "cloudfunctions.googleapis.com",
    "cloudbuild.googleapis.com",
    "artifactregistry.googleapis.com",
    "secretmanager.googleapis.com",
    "cloudscheduler.googleapis.com",
    "datastore.googleapis.com"
  ]

  disable_services_on_destroy = true
  disable_dependent_services  = true
}
