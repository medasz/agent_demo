# AI Agent Demo

一个基于 Go 的命令行 Agent Demo。它会把用户输入交给 LLM，在需要时自动调用本地工具，并把工具结果回填给模型，直到得到最终回复。

## 当前能力

- 多轮对话：通过 `session` 维护上下文消息。
- 流式输出：LLM 响应会通过回调实时打印到终端。
- 工具调用循环：当模型返回 tool call 时，Agent 会执行工具并继续调用模型。
- OpenAI 兼容接口：可接 OpenAI、DeepSeek、Ollama 等兼容服务。

当前内置工具：

- `read_file`
- `list_file`
- `create_file`
- `edit_file`
- `run_command`

## 项目结构

```text
agent_demo/
├── main.go
├── go.mod
├── go.sum
├── README.md
├── internal/
│   ├── agent/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── config/
│   │   ├── config.go
│   │   └── config_test.go
│   ├── llm/
│   │   └── client.go
│   ├── session/
│   │   ├── session.go
│   │   └── session_test.go
│   └── tool/
│       ├── tool.go
│       ├── registry.go
│       ├── registry_test.go
│       ├── local.go
│       ├── read_file.go
│       ├── read_file_test.go
│       ├── list_file.go
│       ├── create_file.go
│       ├── edit_file.go
│       ├── run_command.go
│       └── file_tools_test.go
└── rag/
    └── rag.go
```

## 架构说明

### `internal/agent`

- `Agent.Run()` 负责主流程编排。
- 每轮会：
  - 把用户消息写入 Session
  - 调用 LLM
  - 如果有 tool calls，就执行工具并回填 tool message
  - 如果没有 tool calls，就返回最终内容
- 当前最大工具迭代次数是 `20`。

### `internal/config`

从环境变量加载配置：

- `API_KEY`：必填
- `MODEL_NAME`：必填
- `BASE_URL`：可选

### `internal/llm`

- 封装 `go-openai`
- 提供 `CompleteStream(...)`
- 负责流式拼接内容和 tool calls

### `internal/session`

- 管理会话消息列表
- 对外返回消息副本，避免外部直接篡改内部状态

### `internal/tool`

当前工具模块已经完成从 switch 分发到注册表分发的迁移：

- `Tool` 接口定义在 `tool.go`
- 每个工具各自实现 `Name / Definition / Execute`
- `Registry` 统一负责：
  - 注册工具
  - 提供工具定义列表
  - 根据工具名执行对应工具

`local.go` 目前保留纯文件/命令函数实现，供各工具对象复用。

## 运行方式

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

DeepSeek:

```bash
export API_KEY="your-deepseek-api-key"
export BASE_URL="https://api.deepseek.com/v1"
export MODEL_NAME="deepseek-chat"
```

OpenAI:

```bash
export API_KEY="your-openai-api-key"
export BASE_URL="https://api.openai.com/v1"
export MODEL_NAME="gpt-4o"
```

Ollama:

```bash
export API_KEY="ollama"
export BASE_URL="http://localhost:11434/v1"
export MODEL_NAME="deepseek-coder"
```

### 3. 启动

```bash
go run main.go
```

退出：

```text
exit
quit
```

## 测试

运行全部测试：

```bash
go test ./...
```

单独运行工具测试：

```bash
go test ./internal/tool/...
```

## 当前测试覆盖

- `agent/service_test.go`
  - 无工具调用场景
  - 有工具调用场景
- `config/config_test.go`
  - 配置正常加载
  - 必填配置缺失
- `session/session_test.go`
  - 消息追加
  - 返回副本而不是内部切片
- `tool/registry_test.go`
  - 工具定义完整性
  - 已注册工具执行
  - 未知工具错误
- `tool/read_file_test.go`
  - `read_file` 执行
  - 参数解析失败
- `tool/file_tools_test.go`
  - `list_file`
  - `create_file`
  - `edit_file`
  - `run_command`

## 已知限制

- `Tool` 接口里的 `Definition()` 目前仍直接返回 `openai.Tool`，还存在对 OpenAI 协议的耦合。
- `run_command` 还没有超时、白名单或安全沙箱。
- `Session` 目前只存在内存里，不支持持久化和多会话管理。
- `README` 当前只描述主 Agent 流程，`rag/` 仍是独立实验模块，未接入主链路。
