# OpenTelemetry

# What is OpenTelemetry?

## What is observability?

## Why OpenTelemetry?

# Concepts

## Observability Primer

### What is Observability? (観測可能性とは何か？)

- 観測可能性によって、私たちは外側の世界からシステムを理解することができる
- アプリケーションはsignalを発する
  - トレース
  - メトリクス
  - ログ
- OpenTelemetryは、システムを観測可能にするために、アプリケーションコードが計装されるメカニズム

### Reliability & Metrics (信頼性とメトリクス)

- **Telemetry**: システムから放出される、その動作に関するデータのこと。データはトレース、メトリクス、ログの形式で提供される。

### Understanding Distributed Tracing (分散トレーシングについて理解する)

#### Logs

- **log**: サービスや他のコンポーネントから発せられるタイムスタンプのあるメッセージ

#### Spans

- **span**: 仕事や作業の単位を表現する

#### Distributed Traces

- **distributed trace (trace)**:マイクロサービスやサーバーレスアプリケーションのような複数サービス構成を通過して伝搬するようなリクエストによって取られるパスの記録のこと
  - トレースは1つ以上のスパンによって構成される

## Context Propagation (コンテキストの伝搬)

- コンテキストの伝搬があれば、シグナルがどこで生み出されようとも、お互いを紐づけられることが可能

### Context (コンテキスト)

- **Context**: コンテキストは、送受信サービス（または実行ユニット）がある信号と別の信号を関連付けるための情報を含むオブジェクトである。

### Propagation (伝搬)

# Implementation in Go

## Getting Started

### Create and launch an HTTP server

- `main.go`

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/rolldice", rolldice)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

- `rolldice.go`

```go
package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

func rolldice(w http.ResponseWriter, req *http.Request) {
	roll := 1 + rand.Intn(6)

	resp := strconv.Itoa(roll) + "\n"
	if _, err := io.WriteString(w, resp); err != nil {
		log.Panicf("Write failed: %v\n", err)
	}
}
```

- アプリケーションを走らせてみる

```bash
go run .
```

![alt text](image.png)

### Add OpenTelemetry Instrumentation

- パッケージをインストール

```bash
go get "go.opentelemetry.io/otel" \
  "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric" \
  "go.opentelemetry.io/otel/exporters/stdout/stdouttrace" \
  "go.opentelemetry.io/otel/propagation" \
  "go.opentelemetry.io/otel/sdk/metric" \
  "go.opentelemetry.io/otel/sdk/resource" \
  "go.opentelemetry.io/otel/sdk/trace" \
  "go.opentelemetry.io/otel/semconv/v1.24.0" \
  "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
```

- SDKを初期化

## Instrumentation

## Using instrumentation libraries

## Exporters

## Resources

## Sampling

## API reference

## Examples

## Registry
