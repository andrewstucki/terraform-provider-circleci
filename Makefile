VERSION=v0.0.3

OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

TARGET_BINARY=terraform-provider-circleci_$(VERSION)

TERRAFORM_PLUGIN_DIR=$(HOME)/.terraform.d/plugins/$(OS)_$(ARCH)/

.PHONY: $(TARGET_BINARY) docker

build: $(TARGET_BINARY)

$(TARGET_BINARY):
	CGO_ENABLED=0 go build -ldflags="-s -w" -a -o $(TARGET_BINARY)

docker:
	docker build . -t andrewstucki/go-terraform
	docker push andrewstucki/go-terraform

install: $(TARGET_BINARY)
	mkdir -p $(TERRAFORM_PLUGIN_DIR)
	cp ./$(TARGET_BINARY) $(TERRAFORM_PLUGIN_DIR)/
