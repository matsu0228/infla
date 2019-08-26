# Firewall Rule設定
resource "google_compute_firewall" "icmp" {
  name    = "terraform-allow-icmp"
  network = "${google_compute_network.terraform-vpc.self_link}"

  allow {
    protocol = "icmp"
  }
}

resource "google_compute_firewall" "ssh" {
  name    = "terraform-allow-ssh"
  network = "${google_compute_network.terraform-vpc.self_link}"

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }
}

