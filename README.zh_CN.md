# Coze Studio 轻量工作流版

基于 Coze Studio，保留独立工作流、数据库和插件。后端为 Go，前端为 React/TypeScript，运行依赖为 MySQL、Redis、MinIO。

- 工作流创建、编辑、调试、发布、执行、记录、子工作流、循环、批处理与暂停恢复。
- HTTP、Python、大模型、问答、条件、局部变量、JSON、文本和数据库节点。
- 数据库表结构、记录管理和表格导入导出。
- 插件工具定义、调试、发布、认证和工作流调用。
- 账号、空间、权限、模型配置、文件上传、个人令牌及工作流 OpenAPI。

知识库/RAG、Agent/App、Chatflow、持久会话、全局变量、商城、统一搜索和提示词资源库的实现已删除；不提供完整模式开关。大模型节点内的提示词编辑仍可用。

## 启动

首次部署：

```bash
cp docker/.env.example docker/.env
make web
```

访问 http://localhost:8888。已有环境请保留现有 `.env` 和数据卷。模型供应商、插件凭据及代码执行环境按需配置。

新库创建保留功能需要的表；现有数据库不执行删表、删列或自动 schema 对齐。

详见[运行与兼容说明](docs/workflow-only.md)、[删除与验收报告](docs/workflow-trim-report.md)。

## 开发与验证

```bash
node common/scripts/install-run-rush.js install
node common/scripts/install-run-rush.js build --to @coze-studio/app
cd backend
GOTOOLCHAIN=go1.24.0 go build -ldflags='-checklinkname=0' .
```

Go 数据库集成测试使用 `WORKFLOW_TEST_MYSQL_DSN` 指向隔离测试库；不要指向业务数据库。固定版本 IDL 生成命令和测试命令见验收报告。

基于 [coze-dev/coze-studio](https://github.com/coze-dev/coze-studio)，遵循 [Apache 2.0](LICENSE)；第三方依赖的许可证见 [licenses](licenses)。
