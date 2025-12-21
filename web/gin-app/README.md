# Gin Web Application

Go言語の[Gin](https://gin-gonic.com/)フレームワークを使用したRESTful APIアプリケーションです。ユーザー管理、商品管理、注文管理機能を備えたECサイトのバックエンドAPIとして設計されています。

## 機能

- **ユーザー認証**: JWT認証によるユーザー登録・ログイン
- **ユーザー管理**: プロフィール管理、ロールベースのアクセス制御
- **商品管理**: 商品のCRUD操作、カテゴリー管理
- **注文管理**: 注文の作成、キャンセル、ステータス管理
- **ミドルウェア**: CORS、レートリミット、認証、ロギング
- **データベース**: GORM を使用したPostgreSQL接続
- **グレースフルシャットダウン**: 安全なサーバー停止

## プロジェクト構造

```
web/gin-app/
├── cmd/
│   └── api/
│       └── main.go                 # エントリーポイント
├── internal/
│   ├── config/
│   │   └── config.go              # 設定管理
│   ├── database/
│   │   └── database.go            # データベース接続
│   ├── handlers/
│   │   ├── user_handler.go        # ユーザーハンドラー
│   │   ├── product_handler.go     # 商品ハンドラー
│   │   └── order_handler.go       # 注文ハンドラー
│   ├── middleware/
│   │   ├── auth.go                # 認証ミドルウェア
│   │   ├── cors.go                # CORSミドルウェア
│   │   ├── logger.go              # ロギングミドルウェア
│   │   └── rate_limiter.go        # レートリミッター
│   ├── models/
│   │   ├── user.go                # ユーザーモデル
│   │   ├── product.go             # 商品モデル
│   │   └── order.go               # 注文モデル
│   ├── router/
│   │   └── router.go              # ルーター設定
│   └── utils/
│       ├── jwt.go                 # JWT処理
│       ├── response.go            # レスポンスヘルパー
│       └── validator.go           # バリデーション
├── docs/                          # ドキュメント
├── .env.example                   # 環境変数の例
├── .gitignore
├── go.mod
└── README.md
```

## 必要な環境

- Go 1.21以上
- PostgreSQL 12以上

## セットアップ

### 1. リポジトリのクローン

```bash
cd web/gin-app
```

### 2. 依存パッケージのインストール

```bash
go mod download
```

### 3. 環境変数の設定

`.env.example`をコピーして`.env`を作成し、必要な値を設定します。

```bash
cp .env.example .env
```

`.env`ファイルを編集:

```env
# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=gin_app

# JWT設定
JWT_SECRET=your-secret-key-here

# サーバー設定
SERVER_PORT=8080
GIN_MODE=debug
```

### 4. データベースのセットアップ

PostgreSQLでデータベースを作成:

```sql
CREATE DATABASE gin_app;
```

アプリケーションは起動時に自動的にテーブルをマイグレーションします。

### 5. アプリケーションの起動

```bash
go run cmd/api/main.go
```

サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

### 認証

| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| POST | `/api/v1/auth/register` | ユーザー登録 | 不要 |
| POST | `/api/v1/auth/login` | ログイン | 不要 |

### ユーザー

| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | `/api/v1/users/profile` | プロフィール取得 | 必要 |
| PUT | `/api/v1/users/profile` | プロフィール更新 | 必要 |
| GET | `/api/v1/users` | ユーザー一覧 | 管理者のみ |
| GET | `/api/v1/users/:id` | ユーザー詳細 | 管理者のみ |
| DELETE | `/api/v1/users/:id` | ユーザー削除 | 管理者のみ |

### 商品

| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| GET | `/api/v1/products` | 商品一覧 | 不要 |
| GET | `/api/v1/products/:id` | 商品詳細 | 不要 |
| GET | `/api/v1/products/categories` | カテゴリー一覧 | 不要 |
| POST | `/api/v1/products` | 商品作成 | 管理者のみ |
| PUT | `/api/v1/products/:id` | 商品更新 | 管理者のみ |
| DELETE | `/api/v1/products/:id` | 商品削除 | 管理者のみ |

### 注文

| メソッド | エンドポイント | 説明 | 認証 |
|---------|---------------|------|------|
| POST | `/api/v1/orders` | 注文作成 | 必要 |
| GET | `/api/v1/orders` | 注文一覧 | 必要 |
| GET | `/api/v1/orders/:id` | 注文詳細 | 必要 |
| POST | `/api/v1/orders/:id/cancel` | 注文キャンセル | 必要 |
| PATCH | `/api/v1/orders/:id/status` | ステータス更新 | 管理者のみ |

### その他

| メソッド | エンドポイント | 説明 |
|---------|---------------|------|
| GET | `/health` | ヘルスチェック |
| GET | `/api/v1/docs` | APIドキュメント |

## 使用例

### ユーザー登録

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

### ログイン

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "Password123"
  }'
```

レスポンスでトークンが返されます:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...}
}
```

### 認証が必要なエンドポイントへのアクセス

```bash
curl -X GET http://localhost:8080/api/v1/users/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 商品作成（管理者のみ）

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Sample Product",
    "description": "This is a sample product",
    "price": 1000.00,
    "stock": 50,
    "sku": "SKU001",
    "category": "Electronics"
  }'
```

### 注文作成

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {
        "product_id": 1,
        "quantity": 2
      }
    ],
    "shipping_address": "東京都渋谷区...",
    "billing_address": "東京都渋谷区..."
  }'
```

## ミドルウェア

### 認証ミドルウェア

- JWTトークンによる認証
- `Authorization: Bearer <token>` ヘッダーを検証
- ユーザー情報をコンテキストに設定

### CORSミドルウェア

- クロスオリジンリクエストを許可
- 開発環境用の設定が含まれています

### レートリミッター

- IPアドレスごとにリクエスト数を制限
- デフォルト: 1分間に100リクエスト

### ロギングミドルウェア

- リクエストとレスポンスの詳細をログ出力
- リクエストID生成

## データベースモデル

### User（ユーザー）

- ID, Username, Email, Password（ハッシュ化）
- FirstName, LastName, Role, IsActive
- 作成日時、更新日時、削除日時（ソフトデリート）

### Product（商品）

- ID, Name, Description, Price, Stock
- SKU, Category, ImageURL, IsActive
- 作成日時、更新日時、削除日時

### Order（注文）

- ID, UserID, OrderNumber, Status
- TotalAmount, ShippingAddress, BillingAddress
- 作成日時、更新日時、削除日時

### OrderItem（注文明細）

- ID, OrderID, ProductID, Quantity
- Price（注文時の価格）, Subtotal

## セキュリティ

- パスワードは bcrypt でハッシュ化
- JWT による認証
- ロールベースのアクセス制御
- レートリミッターによるDDoS対策
- 入力値のバリデーション

## ライセンス

MIT License

## 参考リンク

- [Gin Framework](https://gin-gonic.com/)
- [GORM](https://gorm.io/)
- [JWT](https://jwt.io/)
- [PostgreSQL](https://www.postgresql.org/)
