# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_DIR=../bin
CONSUMER_BINARY_NAME=consumer
PRODUCER_BINARY_NAME=producer

# Build both projects
all: consumer producer

# Build the consumer
consumer:
	cd cmd/consumer && $(GOBUILD) -o $(BINARY_DIR)/$(CONSUMER_BINARY_NAME) -v

# Build the producer
producer:
	cd cmd/producer && $(GOBUILD) -o $(BINARY_DIR)/$(PRODUCER_BINARY_NAME) -v

# Clean up binaries
clean:
	$(GOCLEAN)
	rm -f $(BINARY_DIR)/$(CONSUMER_BINARY_NAME)
	rm -f $(BINARY_DIR)/$(PRODUCER_BINARY_NAME)

.PHONY: all consumer producer clean
