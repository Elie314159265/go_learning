# API ドキュメント

## ベースURL

```
http://localhost:8080/api/v1
```

## 認証

ほとんどのエンドポイントは JWT 認証が必要です。ログイン後に取得したトークンをリクエストヘッダーに含めてください。

```
Authorization: Bearer <your_token_here>
```

## レスポンス形式

### 成功レスポンス

```json
{
  "message": "成功メッセージ",
  "data": { ... }
}
```

### エラーレスポンス

```json
{
  "error": "エラーメッセージ",
  "details": { ... }
}
```

## エンドポイント一覧

---

## 認証

### ユーザー登録

```
POST /auth/register
```

**リクエストボディ:**

```json
{
  "username": "testuser",
  "email": "test@example.com",
  "password": "Password123",
  "first_name": "Test",
  "last_name": "User"
}
```

**レスポンス (201 Created):**

```json
{
  "message": "ユーザー登録が完了しました",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "first_name": "Test",
    "last_name": "User",
    "role": "user",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

### ログイン

```
POST /auth/login
```

**リクエストボディ:**

```json
{
  "username": "testuser",
  "password": "Password123"
}
```

**レスポンス (200 OK):**

```json
{
  "message": "ログインに成功しました",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": { ... }
}
```

---

## ユーザー

### プロフィール取得

```
GET /users/profile
```

**認証:** 必要

**レスポンス (200 OK):**

```json
{
  "id": 1,
  "username": "testuser",
  "email": "test@example.com",
  "first_name": "Test",
  "last_name": "User",
  "role": "user",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### プロフィール更新

```
PUT /users/profile
```

**認証:** 必要

**リクエストボディ:**

```json
{
  "email": "newemail@example.com",
  "first_name": "NewFirst",
  "last_name": "NewLast"
}
```

### ユーザー一覧取得

```
GET /users?page=1&page_size=10
```

**認証:** 必要（管理者のみ）

**クエリパラメータ:**
- `page`: ページ番号（デフォルト: 1）
- `page_size`: 1ページあたりの件数（デフォルト: 10）

**レスポンス (200 OK):**

```json
{
  "users": [ ... ],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

---

## 商品

### 商品一覧取得

```
GET /products?page=1&page_size=10&category=Electronics&search=phone
```

**認証:** 不要

**クエリパラメータ:**
- `page`: ページ番号
- `page_size`: 1ページあたりの件数
- `category`: カテゴリーフィルター
- `search`: 検索キーワード
- `active_only`: アクティブな商品のみ（デフォルト: true）

**レスポンス (200 OK):**

```json
{
  "products": [
    {
      "id": 1,
      "name": "Sample Product",
      "description": "This is a sample product",
      "price": 1000.00,
      "stock": 50,
      "sku": "SKU001",
      "category": "Electronics",
      "image_url": "https://example.com/image.jpg",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 100,
  "page": 1,
  "page_size": 10,
  "total_pages": 10
}
```

### 商品詳細取得

```
GET /products/:id
```

**認証:** 不要

### 商品作成

```
POST /products
```

**認証:** 必要（管理者のみ）

**リクエストボディ:**

```json
{
  "name": "Sample Product",
  "description": "This is a sample product",
  "price": 1000.00,
  "stock": 50,
  "sku": "SKU001",
  "category": "Electronics",
  "image_url": "https://example.com/image.jpg"
}
```

### 商品更新

```
PUT /products/:id
```

**認証:** 必要（管理者のみ）

### 商品削除

```
DELETE /products/:id
```

**認証:** 必要（管理者のみ）

### カテゴリー一覧取得

```
GET /products/categories
```

**認証:** 不要

**レスポンス (200 OK):**

```json
{
  "categories": ["Electronics", "Books", "Clothing"]
}
```

---

## 注文

### 注文作成

```
POST /orders
```

**認証:** 必要

**リクエストボディ:**

```json
{
  "items": [
    {
      "product_id": 1,
      "quantity": 2
    },
    {
      "product_id": 2,
      "quantity": 1
    }
  ],
  "shipping_address": "東京都渋谷区...",
  "billing_address": "東京都渋谷区..."
}
```

**レスポンス (201 Created):**

```json
{
  "message": "注文を作成しました",
  "order": {
    "id": 1,
    "user_id": 1,
    "order_number": "ORD20240101123456",
    "status": "pending",
    "total_amount": 3000.00,
    "shipping_address": "東京都渋谷区...",
    "billing_address": "東京都渋谷区...",
    "order_items": [ ... ],
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 注文一覧取得

```
GET /orders?page=1&page_size=10
```

**認証:** 必要

### 注文詳細取得

```
GET /orders/:id
```

**認証:** 必要

### 注文キャンセル

```
POST /orders/:id/cancel
```

**認証:** 必要

### 注文ステータス更新

```
PATCH /orders/:id/status
```

**認証:** 必要（管理者のみ）

**リクエストボディ:**

```json
{
  "status": "shipped"
}
```

**利用可能なステータス:**
- `pending`: 保留中
- `confirmed`: 確認済み
- `shipped`: 発送済み
- `delivered`: 配達完了
- `cancelled`: キャンセル

---

## エラーコード

| ステータスコード | 説明 |
|----------------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 400 | リクエストが不正 |
| 401 | 認証が必要 |
| 403 | アクセス権限がない |
| 404 | リソースが見つからない |
| 409 | 競合（重複など） |
| 429 | リクエスト数が多すぎる |
| 500 | サーバーエラー |

## レート制限

- **制限**: 1分間に100リクエスト
- **超過時**: 429 Too Many Requests
