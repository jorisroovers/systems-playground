################################################################################
# AWS                                                                          #
################################################################################
variable "aws_access_key" {} # Variables can also be lists or maps
variable "aws_secret_key" {}
variable "aws_region" {
  default = "us-east-1"
}

provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}

resource "aws_instance" "myexample-instance" {
  ami           = "ami-0922553b7b0369273"
  instance_type = "m3.medium"
}

# Use 'terraform output' to print these
output "aws-instance-ip" {
  value = "${aws_instance.myexample-instance.public_ip}"
}

################################################################################
# GCP                                                                          #
################################################################################
variable "gcp_project" {}
variable "gcp_region" {
  default = "us-central1"
}

provider "google" {
  credentials = "${file("/tmp/gcp-account.json")}"
  project     = "${var.gcp_project}"
  region      = "${var.gcp_region}"
}

resource "google_compute_instance" "default" {
  name         = "myexample-gcp-instance"
  machine_type = "n1-standard-1"
  zone         = "us-central1-c"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"
  }

}
