# Makefile for Gin Web Application

.PHONY: help run build test clean deps migrate docker-up docker-down

# デフォルトターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  make run          - アプリケーションを起動"
	@echo "  make build        - バイナリをビルド"
	@echo "  make test         - テストを実行"
	@echo "  make clean        - ビルド成果物を削除"
	@echo "  make deps         - 依存パッケージをインストール"
	@echo "  make lint         - コードの静的解析"
	@echo "  make fmt          - コードをフォーマット"
	@echo "  make docker-up    - Dockerコンテナを起動"
	@echo "  make docker-down  - Dockerコンテナを停止"

# アプリケーションを起動
run:
	go run cmd/api/main.go

# バイナリをビルド
build:
	go build -o bin/api cmd/api/main.go

# テストを実行
test:
	go test -v ./...

# カバレッジ付きテスト
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# ビルド成果物を削除
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# 依存パッケージをインストール
deps:
	go mod download
	go mod tidy

# コードの静的解析
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run

# コードをフォーマット
fmt:
	go fmt ./...
	goimports -w .

# Dockerコンテナを起動（PostgreSQL）
docker-up:
	docker-compose up -d

# Dockerコンテナを停止
docker-down:
	docker-compose down

# データベースをリセット
db-reset:
	docker-compose down -v
	docker-compose up -d
	sleep 3
	make run

# 本番用ビルド
build-prod:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go
