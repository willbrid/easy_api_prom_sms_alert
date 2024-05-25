# Makefile

# Variables
BINARY_NAME = easy-api-prom-sms-alert
BUILD_DIR = build
VERSION ?= dev

ARCH = amd64
OS = linux

ARCHIVE = $(BINARY_NAME)-$(OS)-$(ARCH).tar.gz

# Go build command template
BUILD_CMD = GOOS=$(OS) GOARCH=$(ARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-$(OS)-$(ARCH)

# Targets
all: clean build

# Clean build directory
clean:
	rm -rf $(BUILD_DIR)

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Build binary
build: $(BUILD_DIR)
    $(BUILD_CMD)

# Archive the build
archive: build
    tar -czvf $(ARCHIVE) -C $(BUILD_DIR) $(BINARY_NAME)-$(VERSION)-$(OS)-$(ARCH)

.PHONY: all clean build archive