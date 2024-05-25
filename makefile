# Makefile

# Variables
BINARY_NAME = easy-api-prom-sms-alert
BUILD_DIR = build
VERSION ?= dev

ARCH = amd64
OS = linux

ARCHIVE = $(BINARY_NAME)-$(OS)-$(ARCH).tar.gz

# Targets
all: clean build archive

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build binary
build: $(BUILD_DIR)
	GO111MODULE=on GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-$(OS)-$(ARCH) .

# Archive the build
archive: build
	tar -czvf $(ARCHIVE) -C $(BUILD_DIR) $(BINARY_NAME)-$(VERSION)-$(OS)-$(ARCH)

.PHONY: all clean build archive