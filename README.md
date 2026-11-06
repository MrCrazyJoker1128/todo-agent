# TODO Agent

 ![demo](./docs/demo.gif)
 带 AI 副驾的任务管理 + 知识库问答应用。后端 Go（GoFrame），前端 Vue 3，Agent 基于 eino（ReAct），知识库走 ES 向量检索（RAG）。

## 功能

- **任务管理**：日历/列表双视图，创建、改期、改优先级、完成、取消，软删除可找回。
- **AI 副驾**：流式对话，自然语言建任务（"明天下午开会"自动解析时间），人机验证卡确认后才入库；日报/周报/月报总结。
- **知识库**：多库隔离，支持 html 地址抓取、md 文件上传（拖拽/点击，异步入库 + 状态轮询），库内/全库向量检索；聊天时自动 RAG（前置检索拼进 system，多轮指代、闲聊直答）。
- **多轮对话**：近 30 轮历史回放（确认指令转述、取消指令丢弃、think 块剥离），断裂工具调用对自动修复。
- **设置**：对话模型 Key（AES-256-GCM 入库加密存储）+ 嵌入模型地址可配。

## 架构

```
frontend (Vue 3 + Element Plus, :3000 开发 / dist 由后端托管)
    │  HTTP /api
service (GoFrame, :8080)
    ├── PG 17（业务库：用户/任务/聊天/库目录/系统 KEK；扩展 pgvector + zhparser）
    ├── Redis 7（token 缓存）
    ├── Elasticsearch 8（knowledge_chunks 向量索引，content_vector 4096 维）
    └── Ollama（对话模型走 OpenAI 兼容接口；嵌入默认 qwen3-embedding:8b，可在设置页换）
```

agent 链路：`RunAgentStream` → 前置 RAG 检索（`retrieveForChat`，闲聊/历史回放类跳过）→ system（含库目录 `<retrieved>` 资料块）+ 近 30 轮历史 → ReAct → 工具（create/cancel/query/update_task、summarize、search_knowledge、add_knowledge）→ 人机验证 envelope → 落库。

## 快速启动（Docker，一键）

前置：只要 Docker Desktop/Orbstack。PG / ES / Redis / Ollama 全由 compose 新建，
宿主机什么都不用装。数据落四个命名卷（`pgdata/esdata/redisdata/ollama`），
`docker compose down` 不丢，加 `-v` 才清。

```bash
# 1. 配环境（端口/密码按需改；宿主机 5432/9200/6379/11434 被占就改 .env 里对应 PORT）
cp .env.example .env

# 2. 配后端连接：config.docker.yaml 里 PG 密码改成和 .env 里 PG_PASSWORD 一致的值
#    （默认都是 postgres123，不改也能跑；ES/Redis/ollama 地址已指向容器服务名，不用动）
vi service/manifest/config/config.docker.yaml

# 3. 启动（首次构建 app 镜像 + 初始化 PG 库表 + 拉 ES/Redis/ollama 镜像，十几分钟）
docker compose up -d --build

# 4. 灌嵌入模型（ollama 卷是空的，一次即可，几个 G，看网速）
docker compose exec ollama ollama pull qwen3-embedding:8b

# 5. 浏览器开 http://localhost:8080，用演示账号直接登录：
#    邮箱 demo@todo-agent.local，密码 123456（compose 首次启动自动建好）
# 6. 设置页配对话模型 Key（聊天用，需自己申请 MiniMax/OpenAI 兼容 Key）
#    KEK 已预置默认值，开箱即用；生产请在设置页轮换
#    嵌入模型默认 qwen3-embedding:8b，直连 compose 内 ollama，不用改
```

说明：

- KEK（模型 Key 加密钥匙）首次在设置页配置，存库（`db_system_configs`），容器不需要传环境变量。
- 已有本地 PG 想复用：`config.docker.yaml` 里 PG host 改 `host.docker.internal`（容器看宿主机的通道），端口避开冲突即可；compose 的 postgres 服务可照旧留着不用。
- 停：`docker compose down`；清数据卷重来：`docker compose down -v`。

## 本地开发（不走 Docker）

```bash
# 后端（需本地 PG/Redis/ES/ollama，配置见 service/manifest/config/config.example.yaml）
cd service && cp manifest/config/config.example.yaml manifest/config/config.yaml  # 改密码
go run main.go            # 或 gf run main.go（热重载），:8080

# 前端
cd frontend && npm i && npm run dev  # :3000，/dev-api 代理到 :8080/api
```

数据库初始化：`sql/schema.sql` 建表 + 按顺序执行 `sql/migration_*.sql`（全部幂等，可重复执行）。

前端构建：`npm run build`（`vue-tsc -b && vite build`，干净构建无 error 才算过）。

## 配置说明

| 项 | 位置 | 说明 |
|---|---|---|
| 后端端口/中间件地址 | `service/manifest/config/config.yaml`（本地文件，不入库；模板见 `config.example.yaml`/`config.docker.yaml`） | `database.default.link` 填 `pgsql:用户:密码@tcp(主机:端口)/库名`；Redis 要 `address: "host:port"` 单字段（GoFrame 只认这个）；ES 填 `elasticsearch.addresses` |
| 模型 Key / 嵌入模型 | 设置页（按用户存库） | API Key AES-256-GCM 加密存储；嵌入默认本地 ollama `qwen3-embedding:8b`，换模型后历史向量需重建 |
| KEK | 设置页 → 系统 KEK（存 `db_system_configs`） | 丢了则所有用户 Key 解不开；换 KEK 前先重导 Key |

## 仓库结构

```
service/            # GoFrame 后端（api 定义 / controller / service / dao / utility）
  internal/service/chat/      # agent 主循环：前置 RAG、历史回放、工具执行、人机 envelope
  internal/tools/             # eino 工具：search_knowledge / create_task / query_tasks 等
  internal/service/modelconfig/  # 用户模型+嵌入配置解析
frontend/           # Vue 3：TaskPanel 日历列表 / KnowledgeBase 多库管理 / AIAssistant 对话 / settings
sql/                # schema + 幂等 migration（含 compose 初始化脚本 init-compose.sh）
docs/               # 需求/设计文档
compose.yaml / Dockerfile.local  # 一键启动
```

## License

MIT
