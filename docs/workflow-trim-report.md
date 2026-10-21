# 轻量工作流版：源码裁剪与验收报告

日期：2026-10-20。目标：仅保留独立工作流、数据库、插件及必要的账号、空间、权限、模型、上传、令牌和 OpenAPI。没有完整模式开关，没有升级框架或依赖大版本。

## 删除范围

| 层次 | 本轮物理删除与保留边界 |
| --- | --- |
| 后端业务 | 删除知识库/RAG、Agent/App、Chatflow、持久会话、全局变量、商城、统一搜索、提示词资源库的应用层、领域实现、跨域接口、专属配置和测试。 |
| 工作流 | 删除对应节点、注册项和执行分支，保留循环内变量、普通节点输出引用、子工作流、批处理、问答与暂停恢复。移除 AppVariables 执行状态和无用持久化上下文。 |
| 基础设施 | 删除 ES、Milvus、etcd、NSQ、RocketMQ、OceanBase 向量能力，以及 embedding、索引、检索、知识库文档处理和搜索事件总线；不保留空事件总线。 |
| 共享能力 | 表格解析、导入组件和用户标识处理归入数据库/通用上传；插件选择器归入工作流；工作流消息结构与 Agent 实现解耦。 |
| 接口 | 裁剪 Thrift 服务、路由、handler、Go 生成类型、前端 API 方法及类型导出。保留实际引用的共享 DTO 和保留接口字段、路径、节点编号；旧能力不再作为兼容服务存在。 |
| 前端 | 删除旧页面、节点表单、弹窗、状态管理、资源包和专属素材；清理聚合导出、动态引用、Rush 清单、依赖和锁文件。移除资源切换 Chatflow、商城发布、会话绑定、插件 Agent 使用量等旧入口及实现。 |
| 部署 | 删除旧中间件服务、Atlas 自动 schema 对齐、旧向量/事件总线部署文档；Docker/Helm 初始化 SQL 一致，新库创建 23 张基础表，业务数据表按需创建。 |

仍出现的 `bot`、`app` 等历史名称可能是兼容接口路径、共享数据类型、包名或局部组件名称，不代表仍存在 Agent/App 产品实现。知识库/会话/全局变量的少量输入字段仅用于拒绝旧画布，不能用于执行这些能力。

发布新仓库前补查资源引用，另删除 36 张无引用的旧功能图片，前端完整 Rush 构建复验通过；保留应用实际引用的图片资源，编译产物、安装依赖及运行数据不进入新仓库。

## 数据保留与升级

验证期间原 `http://localhost:8888` 使用旧构建；新库和升级副本在独立端口验证，未删除旧库的表、列或数据。

- 原数据库 55 张表复制到隔离库，恢复后逐表记录数与原库一致。
- 副本上验证旧管理员访问管理配置、三个已有测试插件及其工作流调用。
- 已发布文本处理、四则运算、订单查询插件均返回预期结果；除零、订单不存在的异常按预期返回。
- 新版启动不执行 Atlas、不做破坏性 schema 对齐；旧库的非工作流表保留。
- 切换完成后，现有 55 张表逐表记录数与 `CHECKSUM TABLE` 内容校验值均与切换前一致；原管理员权限和三个插件详情正常。
- 旧二进制、静态资源、配置及本轮开始前源码保留在 `bin/trim/backup/`；数据库快照仅位于受限的本地运行目录。

## 验证结果

| 检查 | 结果与边界 |
| --- | --- |
| Go | 所有保留包编译、全仓 Go 单元/集成测试通过；数据库测试连接隔离的 `coze_trim_unit`。新增回归测试确认未知路由返回 404，保留路由仍需鉴权。 |
| 前端 | TypeScript build 检查、136 项目 Rush 构建、修改文件 lint 和关键包 lint 通过。7 个相关包执行测试，11 个测试文件、87 项测试通过；其中两个数据库包没有单元测试文件，数据库通过后端及运行时集成验证。 |
| 依赖 | `go list -deps .` 不包含移除中间件 SDK；工作区依赖图、Rush install/update、Go tidy 通过。移除 332 条已无源码引用的前端生产依赖声明。 |
| 部署 | 后端和前端 Docker 构建通过；Compose 两套配置解析通过；Helm lint/template 通过；未在 Kubernetes 集群部署。 |
| 新环境 | 仅 MySQL、Redis、MinIO；注册、登录、空间、上传/下载、模型配置及跨账号权限隔离通过。 |
| 工作流 | 创建、保存、复制、校验、调试、发布、同步 OpenAPI、流式执行、执行记录、Python、HTTP、子工作流、数据库 SQL、插件调用通过。循环、批处理、问答及暂停恢复通过；另验证已发布 OpenAPI 的 Interrupt → stream_resume → Done。 |
| 数据库 | 表结构/元数据、记录新增/查询/修改/删除、CSV 解析预览及导入、表格模板导出下载通过。 |
| 插件 | 创建、参数定义、调试、发布、工作流调用及异常处理通过；插件授权交互有单元测试覆盖。未连接外部 OAuth 服务进行真实授权。 |
| 旧接口/画布 | 旧知识库、Agent、搜索、Chat 等接口返回 404；不支持节点、知识库绑定、会话历史、全局变量引用明确报错。 |
| 浏览器 | 登录、三类资源列表、数据库结构/记录、插件详情、节点菜单及插件选择、Python 调试结果和工作流发布通过；保留页面未发现请求已删除接口。 |
| 大模型 | 使用本地 OpenAI 协议测试替身，验证模型创建、LLM 节点调试/发布/OpenAPI 调用；没有真实供应商凭据，未验证真实供应商输出、计费或外部网络可靠性。 |

测试日志保存在本地忽略目录 `bin/trim/`：`go-tests-isolated.log`、`middleware-tests.log`、`frontend-tests.log`、`frontend-lint.log`、`lint-changes.log`、`frontend-types.log`、`rush-build-final.log`、`docker-build-final.log`、`docker-frontend-final.log`、`helm-lint.log`、`compose-config.log`、`workflow-smoke.log`、`stream-smoke.log`、`stream-resume-smoke.log`、`extended-smoke.log`、`control-flow-smoke.log`、`question-smoke.log`、`upgrade-smoke.log`、`permission-smoke.log`；切换后验收为 `live-after-cutover.json`，浏览器截图为 `browser-live-delivery.png` 和 `browser-workflow-pass.png`。

流式脚本最初将结束事件写成小写 `done`，实际协议为 `Done`；已修正断言并单独复验通过，见 `stream-smoke.log`。数据库导出验证覆盖原有的表格模板下载，不表示新增了全量记录导出入口。

## 代码量

使用同一个 `scripts/measure-code.py` 统计原始提交 `fefb05ff`、本轮开始前源码归档和最终工作区。包含新增文件；不计安装依赖、构建产物、本地运行目录、数据库文件及其他 Git 忽略产物。统计文本物理行，保留正常格式和必要测试，不使用压缩代码来提高比例。

<!-- SIZE_TABLE_START -->
| 文本行数分类 | 原始 fefb05ff | 本轮开始前 | 本轮完成后 | 本轮净减少 | 较原始降幅 |
| --- | ---: | ---: | ---: | ---: | ---: |
| 后端业务源码 | 145,486 | 145,486 | 79,558 | 65,928 | 45.32% |
| 前端业务源码 | 742,479 | 741,691 | 302,671 | 439,020 | 59.24% |
| 后端测试 | 25,499 | 26,188 | 13,619 | 12,569 | 46.59% |
| 前端测试 | 73,005 | 73,200 | 37,858 | 35,342 | 48.14% |
| 后端生成代码 | 536,550 | 536,573 | 228,477 | 308,096 | 57.42% |
| 前端生成代码 | 277,281 | 277,281 | 30,231 | 247,050 | 89.10% |
| 其他源码（含 SQL/脚本） | 12,699 | 12,612 | 8,212 | 4,400 | 35.33% |
| 文档 | 19,969 | 20,108 | 12,802 | 7,306 | 35.89% |
| 部署配置 | 10,915 | 8,114 | 2,712 | 5,402 | 75.15% |
| 其他文本（含锁文件/IDL/配置） | 233,396 | 233,396 | 172,722 | 60,674 | 26.00% |
| 其他生成代码 | 21 | 21 | 21 | 0 | 0.00% |
| **全仓文本合计** | 2,077,300 | 2,074,670 | 888,883 | 1,185,787 | 57.21% |

本轮减少 **1,185,787 行（57.16%）**；相对原始版本累计减少 **1,188,417 行（57.21%）**。

| 数量指标 | 原始 | 本轮开始前 | 本轮完成后 | 本轮净减少 |
| --- | ---: | ---: | ---: | ---: |
| 全部文件（含二进制素材） | 15,305 | 15,301 | 7,071 | 8,230 |
| 文本文件 | 15,008 | 15,004 | 6,946 | 8,058 |
| Rush 工作区项目 | 259 | 259 | 136 | 123 |
| package.json 文件 | 269 | 269 | 146 | 123 |
| Go require（含间接依赖） | 310 | 310 | 189 | 121 |
| 前端生产依赖声明 | 3,031 | 3,031 | 962 | 2,069 |
| 前端外部生产依赖名称（去重） | 263 | 263 | 202 | 61 |
| 前端 workspace 依赖边 | 1,820 | 1,820 | 488 | 1,332 |
| pnpm 锁文件条目 | 3,931 | 3,931 | 3,054 | 877 |
<!-- SIZE_TABLE_END -->

源码行数不包括生成代码与测试；“其他文本”含未单列的配置、IDL、锁文件和测试 JSON 等；SQL/脚本按源码统计，文档和部署配置单列。文本文件数与包含二进制素材的全部文件数分别统计。依赖声明数量不是安装后的磁盘占用；锁文件条目可包含同包不同版本/peer 组合。

## 复现命令

```bash
python3 scripts/measure-code.py --ref fefb05ff
python3 scripts/measure-code.py --archive bin/trim/backup/source-before.tar.gz
python3 scripts/measure-code.py
node common/scripts/install-run-rush.js install
node common/scripts/install-run-rush.js build --to @coze-studio/app
node frontend/apps/coze-studio/node_modules/typescript/bin/tsc -b frontend/apps/coze-studio/tsconfig.build.json --pretty false
cd backend
GOTOOLCHAIN=go1.24.0 go mod tidy
GOTOOLCHAIN=go1.24.0 go build -ldflags='-checklinkname=0' .
# WORKFLOW_TEST_MYSQL_DSN 必须指向隔离测试库，并先应用精简 schema。
GOTOOLCHAIN=go1.24.0 go test -p 1 -count=1 -ldflags=-checklinkname=0 -gcflags='all=-l -N' ./...
GOTOOLCHAIN=go1.24.0 go run github.com/cloudwego/hertz/cmd/hz@v0.9.7 update --idl ../idl/api.thrift --enable_extends --sort_router --handler_dir api/handler --model_dir api/model --exclude_file api/router/register.go
```

后端生成工具固定为 hz 0.9.7 / thriftgo 0.4.5，保留现有工具链版本。前端生成客户端以保留服务路径和实际类型引用裁剪；IDL/API 变更后应同步检查客户端与类型导出。

## 当前站点与恢复

交付地址：`http://localhost:8888`。保留原管理员、三个测试插件、现有工作流和数据库数据；隔离测试资源不写入业务库。站点使用本次已验证的二进制与静态资源。首次切换因启动脚本误选 `.env.server` 已回退，随后按原站点 `.env` 配置成功切换；前端启动就绪检查也已补齐重试。

回退仅恢复旧构建与配置，不回滚或重置业务数据：运行本地 `python3 bin/trim/rollback.py`。旧文件在 `bin/trim/backup/`，新版发布目录在 `bin/releases/workflow-20261020-trim/`。如需停止隔离验收服务，仅停止 `coze-trim-verify`、`coze-trim-upgrade` 及其测试 Compose 项目，不删除原项目的数据卷。
