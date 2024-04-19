provider "google" {
  credentials = file("/Users/christophebuffard/Documents/Dev/github.com/secrets/rabbitmq-test/rabbitmq-test-420700-ea5fca0a1408.json")
  project     = "rabbitmq-test-420700"
  region      = "us-central1"
}