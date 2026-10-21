# Coze Studio — workflow edition

A focused Coze Studio checkout for standalone workflows, databases and plugins. The Go backend and React/TypeScript frontend run with MySQL, Redis and MinIO.

Retained capabilities include workflow editing, debugging, publishing, execution history, subworkflows, loops, batches, pause/resume, HTTP/Python/LLM/question nodes, database CRUD/import/export, plugin management/authentication, accounts, spaces, models, uploads, personal tokens and workflow OpenAPI.

Knowledge/RAG, Agent/App, Chatflow, persistent conversations, global variables, marketplace, unified search and prompt-library implementations have been removed. There is no full/light mode switch. Inline LLM prompt editing remains available.

## Run

```bash
cp docker/.env.example docker/.env
make web
```

Open http://localhost:8888. Preserve existing environment files and data volumes when upgrading. Fresh databases use the reduced schema; existing tables and columns are not automatically dropped or reconciled.

See the [Chinese README](README.zh_CN.md), [deployment notes](docs/workflow-only.md) and [validation report](docs/workflow-trim-report.md).

## Build

```bash
node common/scripts/install-run-rush.js install
node common/scripts/install-run-rush.js build --to @coze-studio/app
cd backend
GOTOOLCHAIN=go1.24.0 go build -ldflags='-checklinkname=0' .
```

Database integration tests require an isolated database via `WORKFLOW_TEST_MYSQL_DSN`. Never run them against business data.

Derived from [coze-dev/coze-studio](https://github.com/coze-dev/coze-studio), under [Apache 2.0](LICENSE). Third-party notices are retained in [licenses](licenses).
