provider "google" {
  credentials = "${file("secret/spolive-dev-gcpmanager.json")}"
  project     = "spolive-dev"
  region      = "asia-northeast1"
}

resource "google_compute_address" "ftp-static-ip" {
  name   = "ftp-static-ip"
  region = "asia-northeast1"
}

// ref image project: https://cloud.google.com/compute/docs/images#os-compute-support
data "google_compute_image" "my_image" {
  family  = "centos-7"
  project = "gce-uefi-images"
}

resource "google_compute_instance" "instance_template" {
  name = "ansible-generated-by-terraform"

  machine_type = "f1-micro"
  zone         = "asia-northeast1-a"

  // boot disk
  boot_disk {
    initialize_params {
      image = "${data.google_compute_image.my_image.self_link}"
    }
  }

  network_interface {
    network    = "${google_compute_network.terraform-vpc.self_link}"
    subnetwork = "${google_compute_subnetwork.subnet1.name}"
    access_config {
      nat_ip = "${google_compute_address.ftp-static-ip.address}"
    }
  }

  // TODO: サービスアカウント周りの設定
  // service_account {
  //   scopes = ["userinfo-email", "compute-ro", "storage-ro"]
  // }

  metadata_startup_script = <<EOT
    echo hi > /test.txt
    timedatectl set-timezone Asia/Tokyo
    sudo yum install -y epel-release
    sudo yum install -y ansible
  EOT

  metadata = {
    "block-project-ssh-keys" = "true"
    "sshKeys"                = "${var.ansible_ssh_keys}"
  }
}


resource "google_compute_instance" "instance_target" {
  name = "target-generated-by-terraform"

  machine_type = "f1-micro"
  zone         = "asia-northeast1-a"

  // boot disk
  boot_disk {
    initialize_params {
      image = "${data.google_compute_image.my_image.self_link}"
    }
  }

  network_interface {
    network    = "${google_compute_network.terraform-vpc.self_link}"
    subnetwork = "${google_compute_subnetwork.subnet1.name}"
  }

  metadata_startup_script = <<EOT
    timedatectl set-timezone Asia/Tokyo
  EOT

  metadata = {
    "block-project-ssh-keys" = "true"
    "sshKeys"                = "${var.ansible_ssh_keys}"
  }
}