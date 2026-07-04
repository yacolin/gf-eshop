ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack.mk

# ============================================================
# 开发常用命令
# ============================================================

# Build and run the server (hot-reload with gf CLI).
.PHONY: run
run:
	@gf run main.go

# Production build and run.
.PHONY: run.prod
run.prod:
	@go build -o main . && ./main

# Build binary.
.PHONY: build
build:
	@gf build -ew


stop:
	@echo "Stopping app on :8000..."
	-@lsof -ti :8000 | xargs kill -9 2>/dev/null
