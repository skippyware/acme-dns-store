DOCKER_IMAGE_NAME := skippyprime/acme-dns-store
DOCKER_BUILD_PLATFORMS := linux/amd64,linux/arm64,linux/arm/v7,linux/arm/v6

build:
	docker build -f Dockerfile -t $(DOCKER_IMAGE_NAME):v1.0-alpine --platform $(DOCKER_BUILD_PLATFORMS) .

tag: build
	docker tag $(DOCKER_IMAGE_NAME):v1.0-alpine $(DOCKER_IMAGE_NAME):latest
	docker tag $(DOCKER_IMAGE_NAME):v1.0-alpine $(DOCKER_IMAGE_NAME):v1.0

push: tag
	docker push $(DOCKER_IMAGE_NAME):latest
	docker push $(DOCKER_IMAGE_NAME):v1.0-alpine
	docker push $(DOCKER_IMAGE_NAME):v1.0

.PHONY: build tag push