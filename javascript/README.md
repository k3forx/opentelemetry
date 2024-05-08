# OpenTelemetry in JavaScript

## Setup Example App

```bash
❯ npm init -y

❯ npm install typescript \
  ts-node \
  @types/node \
  express \
  @types/express

❯ npx tsc --init
```

- launch App

```bash
❯ npx ts-node app.ts
Listening for requests on http://localhost:8080
```

## Instrumentation

- Install the Node SDK and autoinstrumentations package

```bash
❯ npm install @opentelemetry/sdk-node \
  @opentelemetry/api \
  @opentelemetry/auto-instrumentations-node \
  @opentelemetry/sdk-metrics \
  @opentelemetry/sdk-trace-node
```
