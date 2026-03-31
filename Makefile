VERSION ?= dev
TARGET ?= linux/arm64

.PHONY: build-all
build-all:
	$(MAKE) build TARGET=linux/amd64
	$(MAKE) build TARGET=linux/arm64

.PHONY: build
build:
	docker buildx build \
		--platform $(TARGET) \
		--load \
		--build-arg VERSION=$(VERSION) \
		-t wait_for_response:$(VERSION) .

.PHONY: run-dev
run-dev:
	docker run -it --rm \
		wait_for_response:$(VERSION) \
		-url=https://ripe.net -code=200 -timeout=1000 -interval=500