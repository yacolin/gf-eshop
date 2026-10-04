ROOT_DIR    = $(shell pwd)
NAMESPACE   = "default"
DEPLOY_NAME = "template-single"
DOCKER_NAME = "template-single"

include ./hack/hack.mk

# ============================================================
# 开发常用命令
# ============================================================

.PHONY: run run.prod build stop
run: ## 开发模式启动（gf CLI 热加载）
	@gf run main.go

run.prod: ## 本机编译并前台运行
	@go build -o main . && ./main

build: ## 本机平台编译（gf build，会打包 resource 进二进制）
	@gf build -ew

stop: ## 停掉 :8000 上的开发服务
	@echo "Stopping app on :8000..."
	-@lsof -ti :8000 | xargs kill -9 2>/dev/null

# ============================================================
# 发布：本地交叉编译 → 打包 → 上传（参考 campus_express）
#
# 为什么在本地交叉编译：线上多是小内存机器，现场 go build 的链接阶段内存峰值 1G+，
# 容易被 OOM killer 干掉；一律本地编译成静态二进制（CGO_ENABLED=0，不依赖目标机
# libc/glibc 版本）再上传。平台可覆盖：make release GOARCH=arm64
# ============================================================

GOOS   ?= linux
GOARCH ?= amd64

RELEASE_ROOT  := bin/release
RELEASE_STAGE := $(RELEASE_ROOT)/gf-eshop
RELEASE_BIN   := bin/gf-eshop-$(GOOS)-$(GOARCH)
ARTIFACT      := bin/gf-eshop-$(GOOS)-$(GOARCH).tar.gz

# 导出本库结构，供服务器首次建库（只导结构不导数据；仓库里没有建库脚本）
DB_HOST     ?= 127.0.0.1
DB_PORT     ?= 3306
DB_USER     ?= root
DB_PASSWORD ?= 123456
DB_NAME     ?= eshop_db
SCHEMA_SQL  := bin/schema.sql

# 上传目标：make upload DEPLOY_HOST=root@1.2.3.4 [DEPLOY_DIR=/tmp]
DEPLOY_HOST ?=
DEPLOY_DIR  ?= /tmp

.PHONY: help release upload deploy db-schema fmt vet test

help: ## 显示所有命令
	@grep -hE '^[a-zA-Z0-9_.-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "};{printf "  %-12s %s\n", $$1, $$2}'

release: ## 交叉编译并打包部署包（默认 linux/amd64，产物 bin/*.tar.gz）
	@rm -rf $(RELEASE_STAGE)
	@mkdir -p $(RELEASE_STAGE)/manifest/config
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath -o $(RELEASE_BIN) .
	@cp $(RELEASE_BIN) $(RELEASE_STAGE)/gf-eshop
	@cp manifest/config/config.yaml $(RELEASE_STAGE)/manifest/config/config.yaml
	@cp deploy/start.sh deploy/es-watch.sh deploy/gf-eshop.service deploy/INSTALL.md deploy/env.example $(RELEASE_STAGE)/
	@mkdir -p $(RELEASE_STAGE)/db
	@if [ -f $(SCHEMA_SQL) ]; then cp $(SCHEMA_SQL) $(RELEASE_STAGE)/db/schema.sql; \
		echo "  已内置 db/schema.sql（首次部署建库用）"; \
	else echo "  提示：未找到 $(SCHEMA_SQL)。首次部署建库请先执行 make db-schema"; fi
	@chmod 755 $(RELEASE_STAGE)/gf-eshop $(RELEASE_STAGE)/start.sh $(RELEASE_STAGE)/es-watch.sh
	@chmod 644 $(RELEASE_STAGE)/gf-eshop.service $(RELEASE_STAGE)/INSTALL.md
	@[ ! -f $(RELEASE_STAGE)/db/schema.sql ] || chmod 644 $(RELEASE_STAGE)/db/schema.sql
	@chmod 600 $(RELEASE_STAGE)/manifest/config/config.yaml
	@tar -czf $(ARTIFACT) -C $(RELEASE_ROOT) gf-eshop
	@echo "== 部署包已生成 =="
	@file $(RELEASE_BIN) | sed 's/^/  /'
	@printf "  %s\n" "$$(du -h $(ARTIFACT) | cut -f1)  $(ARTIFACT)"
	@echo "  sha256: $$(shasum -a 256 $(ARTIFACT) | cut -d' ' -f1)"
	@go version -m $(RELEASE_BIN) 2>/dev/null | grep -E 'vcs\.(revision|modified)' | sed 's/^/  /' || true
	@echo "  内含: gf-eshop 二进制 + manifest/config/config.yaml(600)"
	@echo "        + start.sh + es-watch.sh + gf-eshop.service + INSTALL.md + env.example + db/schema.sql(如有)"
	@echo "  配置: 服务器真实配置放 /etc/gf-eshop.env（参考 env.example），它会覆盖 config.yaml 同名项"

upload: ## 上传部署包（需 DEPLOY_HOST，如 make upload DEPLOY_HOST=root@1.2.3.4）
	@[ -n "$(DEPLOY_HOST)" ] || { echo "缺少 DEPLOY_HOST。用法：make upload DEPLOY_HOST=root@1.2.3.4 [DEPLOY_DIR=/tmp]"; exit 1; }
	@[ -f $(ARTIFACT) ] || { echo "$(ARTIFACT) 不存在，请先执行 make release"; exit 1; }
	scp $(ARTIFACT) $(DEPLOY_HOST):$(DEPLOY_DIR)/
	@echo "已上传到 $(DEPLOY_HOST):$(DEPLOY_DIR)/$(notdir $(ARTIFACT))，接着在服务器上执行："
	@echo "  sudo tar -xzf $(DEPLOY_DIR)/$(notdir $(ARTIFACT)) -C /opt && sudo systemctl restart gf-eshop"
	@echo "  （首次部署请先按 INSTALL.md 建库、改配置、装 systemd 单元）"

deploy: release upload ## 打包并上传（release + upload）

db-schema: ## 导出本地库结构到 bin/schema.sql（只导结构，随 release 一起打包）
	@mkdir -p bin
	@{ \
		echo "-- gf-eshop 表结构（make db-schema 生成，仅结构、不含数据）"; \
		echo "-- 生成时间: $$(date '+%Y-%m-%d %H:%M:%S')  来源: $$(hostname):$(DB_NAME)"; \
		echo "--"; \
		echo "-- ⚠️ 本文件含 DROP TABLE IF EXISTS：在【已有数据】的库上执行会先删表再重建，等于清空数据。"; \
		echo "--    仅用于首次建库，或你明确要重建结构的场景。"; \
		echo "--    已去掉 DEFINER / SQL_LOG_BIN / GTID_PURGED 等需要 SUPER 权限的语句，可直接导入云 RDS。"; \
		mysqldump -h$(DB_HOST) -P$(DB_PORT) -u$(DB_USER) -p$(DB_PASSWORD) \
			--no-data --skip-comments --single-transaction --set-gtid-purged=OFF \
			$(DB_NAME); \
	} > $(SCHEMA_SQL)
	@echo "已导出 $(SCHEMA_SQL)：$$(grep -c 'CREATE TABLE' $(SCHEMA_SQL)) 张表、$$(grep -c '^DROP TABLE' $(SCHEMA_SQL)) 条 DROP"

fmt: ## 格式化代码
	gofmt -l -w .

vet: ## 静态检查
	go vet ./...

test: ## 运行测试
	go test ./...
