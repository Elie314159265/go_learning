# Go Learning Repository

Go言語を体系的に学ぶためのリポジトリ

## ディレクトリ構成

```
go_learning/
├── basics/          # 基礎文法・型・制御構文
├── concurrency/     # Goroutines 関連
├── interfaces/      # インターフェースとポリモーフィズム
├── testing/         # テスト駆動開発
├── errors/          # エラーハンドリング
├── stdlib/          # 標準ライブラリ
├── web/             # Webフレームワーク
│   ├── stdlib-http/ # 標準ライブラリ
│   ├── echo/        # Echo フレームワーク
│   └── gin/         # Gin フレームワーク
├── algorithms/      # アルゴリズムとデータ構造
└── projects/        # 実践プロジェクト
```

## 学習順序の推奨

1. **basics/** - Go言語の基礎を理解する
2. **stdlib/** - 標準ライブラリに慣れる
3. **errors/** - エラーハンドリングを学ぶ
4. **interfaces/** - インターフェースの概念を理解する
5. **testing/** - テストの書き方を学ぶ
6. **concurrency/** - Goの並行処理を習得する
7. **algorithms/** - アルゴリズムを実装して練習する
8. **web/** - Webアプリケーション開発を学ぶ
9. **projects/** - 実践的なプロジェクトを構築する

## 使い方

各ディレクトリには `README.md` とサンプルコードが含まれています。

```bash
# サンプルコードの実行
cd basics
go run hello.go

# テストの実行
cd testing
go test -v
```

## 学習リソース

- [A Tour of Go](https://go.dev/tour/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go公式ドキュメント](https://go.dev/doc/)

## 環境構築

```bash
# Goのバージョン確認
go version

# モジュールの初期化（プロジェクトごと）
go mod init github.com/yourusername/projectname

# 依存関係のインストール
go mod tidy
```

## ライセンス

MIT License

---
最終更新: 2025-12-03
