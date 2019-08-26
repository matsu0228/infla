resource "google_compute_network" "terraform-vpc" {
  name = "terraform-vpc"
}

resource "google_compute_subnetwork" "subnet1" {
  name          = "subnet1"
  ip_cidr_range = "192.168.10.0/24"
  network       = "${google_compute_network.terraform-vpc.name}"
  description   = "terraform-vpc.subnet1"
  region        = "${var.region}"
}