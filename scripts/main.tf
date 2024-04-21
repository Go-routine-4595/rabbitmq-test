terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 3.5"
    }
  }

  required_version = ">= 0.12"
}

// provider "google" {
//  credentials = file("<PATH_TO_YOUR_SERVICE_ACCOUNT_KEY>.json")
//  project     = "rabbitmq-test-420700"
//  region      = "us-central1"
//}

resource "google_compute_instance" "vm" {
  name         = "rabbitmq-instance"
  machine_type = "e2-micro"
  zone         = "us-central1-a"

  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-2004-lts"
    }
  }

  network_interface {
    network = "default"
    access_config {
      // Ephemeral IP
    }
  }

  metadata_startup_script = <<-EOS
    #!/bin/bash
    apt-get update
    apt-get install curl
    apt-get install -y docker.io
    systemctl start docker
    systemctl enable docker
    docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:management

    # Wait for RabbitMQ to start
    sleep 30

    # Install rabbitmqadmin from the RabbitMQ management plugin
    while ! curl -s http://localhost:15672/cli/rabbitmqadmin > /usr/local/bin/rabbitmqadmin; do
      echo "Waiting for RabbitMQ to start..."
      sleep 10
    done

    chmod +x /usr/local/bin/rabbitmqadmin

    sleep 10
    /usr/local/bin/rabbitmqadmin declare exchange name=ref type=fanout durable=true
    /usr/local/bin/rabbitmqadmin declare exchange name=hw type=fanout durable=true
    /usr/local/bin/rabbitmqadmin declare exchange name=fmi type=fanout durable=true
    /usr/local/bin/rabbitmqadmin declare queue name=alarms durable=false
    /usr/local/bin/rabbitmqadmin declare queue name=fmi_alarms durable=true
    /usr/local/bin/rabbitmqadmin declare binding source=hw destination=alarms
    /usr/local/bin/rabbitmqadmin declare binding source=fmi destination=fmi_alarms
    /usr/local/bin/rabbitmqadmin declare binding source=hw destination_type=exchange destination=fmi

    EOS
}

resource "google_compute_firewall" "rabbitmq" {
  name    = "allow-rabbitmq-2"
  network = "default"

  allow {
    protocol = "tcp"
    ports    = ["5672", "15672"]
  }

  source_ranges = ["0.0.0.0/0"]
  target_tags   = ["rabbitmq-instance"]
}
