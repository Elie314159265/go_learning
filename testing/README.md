# Testing

Goのテストを学ぶディレクトリ

## トピック

- ユニットテスト（`*_test.go`ファイル）
- Table-driven tests（テーブル駆動テスト）
- ベンチマーク
- テストカバレッジ
- モック
- サブテスト

## 実行方法

```bash
# テストの実行
go test

# 詳細表示
go test -v

# カバレッジ
go test -cover

# ベンチマーク
go test -bench=.
```

## ポイント

Goのテストは標準ライブラリに組み込まれており、追加のフレームワークは不要です。
