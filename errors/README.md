# Errors

Goのエラーハンドリングを学ぶディレクトリ

## トピック

- エラーの基本（errorインターフェース）
- カスタムエラー型
- Sentinel errors（定義済みエラー）
- エラーのラップ（`fmt.Errorf` with `%w`）
- `errors.Is()` と `errors.As()`
- panic と recover
- エラーハンドリングのベストプラクティス

## 実行方法

```bash
go run custom_errors.go
```

## ポイント

Goでは例外の代わりにエラー値を返します。明示的なエラーチェックが推奨されます。
