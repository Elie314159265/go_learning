# Concurrency

Goの並行処理を学ぶディレクトリ

## トピック

- Goroutines（軽量スレッド）
- Channels（goroutine間の通信）
- Select文
- Buffered channels
- WaitGroups
- Mutex（排他制御）

## 実行方法

```bash
go run goroutines.go
go run channels.go
```

## 重要な原則

> "Do not communicate by sharing memory; instead, share memory by communicating."
