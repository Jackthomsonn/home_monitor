module "datastore" {
  source  = "terraform-google-modules/cloud-datastore/google"
  project = var.project
  indexes = file("index.yaml")
  depends_on = [
    module.project_services
  ]
}
