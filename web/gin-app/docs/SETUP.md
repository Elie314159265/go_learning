# セットアップガイド

## 前提条件

以下がインストールされている必要があります:

- Go 1.21以上
- PostgreSQL 12以上（またはDocker）
- Git

## 手順

### 1. プロジェクトのクローン

```bash
git clone <repository-url>
cd web/gin-app
```

### 2. 依存パッケージのインストール

```bash
make deps
# または
go mod download
```

### 3. データベースのセットアップ

#### オプション A: Dockerを使用（推奨）

```bash
# PostgreSQLコンテナを起動
make docker-up

# または
docker-compose up -d
```

これにより以下が起動します:
- PostgreSQL (ポート 5432)
- pgAdmin (ポート 5050) - http://localhost:5050

pgAdminログイン情報:
- Email: admin@example.com
- Password: admin

#### オプション B: ローカルのPostgreSQLを使用

PostgreSQLをインストールし、データベースを作成:

```sql
CREATE DATABASE gin_app;
```

### 4. 環境変数の設定

```bash
# .env.exampleをコピー
cp .env.example .env

# .envファイルを編集
# 必要に応じてデータベース接続情報などを変更
```

重要な環境変数:

```env
# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gin_app

# JWT設定（本番環境では必ず変更すること！）
JWT_SECRET=your-very-secret-key-change-in-production
```

### 5. アプリケーションの起動

```bash
make run
# または
go run cmd/api/main.go
```

アプリケーションは http://localhost:8080 で起動します。

### 6. 動作確認

#### ヘルスチェック

```bash
curl http://localhost:8080/health
```

期待されるレスポンス:

```json
{
  "status": "healthy",
  "database": "connected",
  "version": "1.0.0",
  "environment": "development"
}
```

#### ユーザー登録

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "Password123",
    "first_name": "Test",
    "last_name": "User"
  }'
```

#### ログイン

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Password123"
  }'
```

レスポンスからトークンを取得し、以降のリクエストで使用します。

## トラブルシューティング

### データベース接続エラー

**症状:**
```
データベース接続に失敗しました: dial tcp [::1]:5432: connect: connection refused
```

**解決方法:**
1. PostgreSQLが起動しているか確認
   ```bash
   # Dockerの場合
   docker ps | grep postgres

   # ローカルの場合
   pg_isadmin
   ```

2. 接続情報が正しいか確認
   - `.env`ファイルの`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`を確認

### ポートが既に使用されている

**症状:**
```
listen tcp :8080: bind: address already in use
```

**解決方法:**
1. 使用中のポートを確認
   ```bash
   lsof -i :8080
   ```

2. `.env`ファイルで別のポートを指定
   ```env
   SERVER_PORT=8081
   ```

### マイグレーションエラー

**症状:**
```
マイグレーションに失敗しました
```

**解決方法:**
1. データベースをリセット
   ```bash
   make db-reset
   ```

2. または手動でデータベースを削除して再作成
   ```sql
   DROP DATABASE gin_app;
   CREATE DATABASE gin_app;
   ```

## 開発時のTips

### コードのフォーマット

```bash
make fmt
```

### 静的解析

```bash
make lint
```

### テストの実行

```bash
make test
```

### ホットリロード（開発時）

[Air](https://github.com/cosmtrek/air)を使用すると、ファイル変更時に自動でリロードできます。

```bash
# Airのインストール
go install github.com/cosmtrek/air@latest

# 実行
air
```

### データベースの確認

pgAdminを使用する場合:

1. http://localhost:5050 にアクセス
2. Email: admin@example.com, Password: admin でログイン
3. サーバーを追加
   - Name: gin-app
   - Host: postgres
   - Port: 5432
   - Username: postgres
   - Password: postgres

## 本番環境へのデプロイ

### ビルド

```bash
make build-prod
```

### 環境変数の設定

本番環境では以下を必ず変更してください:

```env
GIN_MODE=release
JWT_SECRET=<強力なランダムな文字列>
APP_ENV=production
DB_SSLMODE=require
```

### Dockerイメージの作成

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bin/api cmd/api/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/api .
EXPOSE 8080
CMD ["./api"]
```

## 次のステップ

- [API ドキュメント](./API.md)を参照してエンドポイントを確認
- [アーキテクチャ](./ARCHITECTURE.md)を理解してコードベースに貢献
- テストを追加してコードの品質を向上
