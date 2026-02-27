APP_NAME := adam
APP_VERSION := v1.0.0
REGISTRY_URL ?=
DOCKERFILE := Dockerfile
CRON_DOCKERFILE := cron/Dockerfile

# Build image names
ifeq ($(REGISTRY_URL),)
    IMAGE_NAME := $(APP_NAME)
else
    IMAGE_NAME := $(REGISTRY_URL)/$(APP_NAME)
endif

.PHONY: all build run build-local build-amd64 push-amd64 tag clean help

# Local Go build and run
run:
	go build -o $(APP_NAME) . && ./$(APP_NAME)

all: build-amd64 push

build:
	docker buildx build --platform linux/amd64,linux/arm64 -f $(DOCKERFILE) -t $(IMAGE_NAME):$(APP_VERSION) -t $(IMAGE_NAME):latest --push .
	docker buildx build --platform linux/amd64,linux/arm64 -f $(CRON_DOCKERFILE) -t $(IMAGE_NAME)-cron:$(APP_VERSION) -t $(IMAGE_NAME)-cron:latest --push .

build-amd64:
	docker buildx build --platform linux/amd64 -f $(DOCKERFILE) -t $(IMAGE_NAME):$(APP_VERSION) -t $(IMAGE_NAME):latest --load .
	docker buildx build --platform linux/amd64 -f $(CRON_DOCKERFILE) -t $(IMAGE_NAME)-cron:$(APP_VERSION) -t $(IMAGE_NAME)-cron:latest --load .

push:
	docker push $(IMAGE_NAME):$(APP_VERSION)
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME)-cron:$(APP_VERSION)
	docker push $(IMAGE_NAME)-cron:latest

push-amd64:
	docker buildx build --platform linux/amd64 -f $(DOCKERFILE) -t $(IMAGE_NAME):$(APP_VERSION) -t $(IMAGE_NAME):latest --push .
	docker buildx build --platform linux/amd64 -f $(CRON_DOCKERFILE) -t $(IMAGE_NAME)-cron:$(APP_VERSION) -t $(IMAGE_NAME)-cron:latest --push .

build-local:
	docker build -f $(DOCKERFILE) -t $(APP_NAME):$(APP_VERSION) -t $(APP_NAME):latest .
	docker build -f $(CRON_DOCKERFILE) -t $(APP_NAME)-cron:$(APP_VERSION) -t $(APP_NAME)-cron:latest .

tag:
	@echo "Current tags:"
	@docker images | grep $(APP_NAME)

clean:
	@echo "Cleaning images..."
	-docker rmi $(IMAGE_NAME):$(APP_VERSION) 2>/dev/null || true
	-docker rmi $(IMAGE_NAME):latest 2>/dev/null || true
	-docker rmi $(IMAGE_NAME)-cron:$(APP_VERSION) 2>/dev/null || true
	-docker rmi $(IMAGE_NAME)-cron:latest 2>/dev/null || true

help:
	@echo "Available targets:"
	@echo "  make run                   - Build and run Go binary locally"
	@echo "  make all                   - Build multi-platform and push (requires REGISTRY_URL)"
	@echo "  make build                 - Build multi-platform image"
	@echo "  make build-amd64           - Build AMD64 only and load locally"
	@echo "  make push                  - Push image to registry"
	@echo "  make push-amd64            - Build AMD64 only and push to registry"
	@echo "  make build-local           - Build without registry prefix"
	@echo "  make tag                   - Show image tags"
	@echo "  make clean                 - Remove built images"
	@echo "  make help                  - Show this help"
	@echo ""
	@echo "Usage:"
	@echo "  REGISTRY_URL=registry.example.com make build-amd64"
	@echo "  REGISTRY_URL=registry.example.com make push-amd64"
	@echo "  REGISTRY_URL=registry.gitlab.com/username make build-amd64 push"