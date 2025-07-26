# gRPC Gateway Distributed System

このプロジェクトは、gRPC Gateway を使用したマイクロサービス アーキテクチャのサンプル実装です。

## アーキテクチャ

```
┌─────────────────┐    HTTP/REST    ┌─────────────────┐
│   HTTP Client   │ ──────────────→ │  Gateway Service │ :8080
└─────────────────┘                 │ (gRPC Gateway)  │
                                    └─────────────────┘
                                             │ gRPC
                         ┌───────────────────┼───────────────────┐
                         │                   │                   │
                         ▼                   ▼                   ▼
              ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
              │  Book Service   │  │ Author Service  │  │  MySQL Database │
              │     :9091       │  │     :9092       │  │     :3307       │
              │   (gRPC Only)   │  │   (gRPC Only)   │  │                 │
              └─────────────────┘  └─────────────────┘  └─────────────────┘
                         │                   ▲                   │
                         │       gRPC        │                   │
                         └───────────────────┘                   │
                                                                 │
                                        ┌─────────────────┐      │
                                        │     Jaeger      │      │
                                        │   :16687 UI     │ ◄────┘
                                        │ (Tracing/Logs)  │
                                        └─────────────────┘
```

## サービス詳細

### Gateway Service (ポート: 8080)

- **役割**: gRPC Gateway（HTTP ⇄ gRPC 変換）
- **機能**: HTTP REST API を受け付けて、バックエンド gRPC サービスに転送

### Book Service (ポート: 9091)

- **役割**: Book 専用 gRPC サーバー
- **機能**: 書籍情報の取得、Author Service への gRPC 呼び出し

### Author Service (ポート: 9092)

- **役割**: Author 専用 gRPC サーバー
- **機能**: 著者情報の取得

### MySQL Database (ポート: 3307)

- **役割**: データ永続化
- **共有**: 全サービスが同一データベースにアクセス

### Jaeger (ポート: 16687)

- **役割**: 分散トレーシングと観測性
- **機能**: サービス間の通信トレース、パフォーマンス監視

## API エンドポイント

### Books API

```bash
# 書籍情報を取得
curl -X GET "http://localhost:8080/v1/books/{id}"
```

### Authors API

```bash
# 著者情報を取得
curl -X GET "http://localhost:8080/v1/authors/{id}"
```

## 開発環境のセットアップ

### 前提条件

- Docker & Docker Compose
- Go 1.24.1+
- buf CLI

### 起動方法

1. **サービスをビルド:**

   ```bash
   cd services
   make build
   ```

2. **サービスを起動:**

   ```bash
   make up
   ```

3. **API をテスト:**

   ```bash
   make test
   ```

4. **Jaeger UI を開く:**

   ```bash
   make jaeger
   # または直接アクセス: http://localhost:16687
   ```

5. **サービスを停止:**
   ```bash
   make down
   ```

## 利用可能なコマンド

```bash
cd services

# サービスをビルド
make build

# サービスを起動
make up

# サービスを停止
make down

# 全データを削除してクリーンアップ
make clean

# サービス状態を確認
make status

# ログを表示
make logs

# API テスト実行
make test

# Jaeger UI アクセス情報を表示
make jaeger
```

## 監視と観測性

### Jaeger UI

- **URL**: http://localhost:16687
- **機能**: 分散トレーシング、サービス依存関係の可視化、パフォーマンス分析

#### トレースの確認方法

1. API リクエストを実行 (`make test`)
2. Jaeger UI にアクセス
3. サービス選択（gateway, book-service, author-service）
4. トレースを検索・分析

## Protocol Buffers

Protocol Buffers の定義は `proto/` ディレクトリにあります：

- `proto/book/v1/book.proto` - Book Service の API 定義
- `proto/author/v1/author.proto` - Author Service の API 定義

### コード生成

```bash
# Protocol Buffers からコード生成
buf generate --include-imports
```

## 技術スタック

- **API Gateway**: gRPC Gateway
- **Microservices**: gRPC サーバー (Go)
- **Protocol Buffers**: API 定義とコード生成
- **Database**: MySQL 8.3
- **Containerization**: Docker & Docker Compose
- **Live Reload**: Air
- **Observability**: OpenTelemetry + Jaeger
- **Build Tool**: buf

## テスト例

### Book API テスト

```bash
curl -s "http://localhost:8080/v1/books/1" | jq .
```

レスポンス:

```json
{
  "book": {
    "id": "1",
    "title": "容疑者Xの献身",
    "authorName": "東野圭吾"
  }
}
```

### Author API テスト

```bash
curl -s "http://localhost:8080/v1/authors/1" | jq .
```

レスポンス:

```json
{
  "author": {
    "id": "1",
    "name": "東野圭吾",
    "bio": ""
  }
}
```

## 分散トレーシング

API リクエストにより、以下のトレースが生成されます：

1. **Gateway → Book Service** (gRPC)
2. **Book Service → Author Service** (gRPC)
3. **各サービス → MySQL** (SQL クエリ)

すべてのトレースは Jaeger UI で確認できます。
