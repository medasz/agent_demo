# AI Agent Demo - 智能系统操作助手

一个基于 Go 语言开发的命令行 AI 助手，能够通过自然语言指令操作文件系统和执行系统命令。

## 🚀 项目概述

本项目演示了如何将 AI 能力与系统操作工具结合，实现通过自然语言控制文件系统和执行命令的功能。项目使用 Go 语言开发，采用**模块化架构**，支持多种 AI 后端（DeepSeek、OpenAI、Ollama 等）。

项目代码**全面覆盖单元测试**，核心模块均有对应的测试文件，确保代码质量与可维护性。

### 项目亮点

| 特性 | 说明 |
|------|------|
| 🤖 自然语言驱动 | 用中文或英文与 AI 对话，自动调用工具完成文件操作、命令执行等任务 |
| 🧩 模块化设计 | agent / config / llm / session / tool 五层解耦，接口清晰 |
| 🔌 多 AI 后端 | 兼容任意 OpenAI 格式的 API（DeepSeek、OpenAI、Ollama 等） |
| ⚡ 流式响应 | 实时显示 AI 输出，体验更流畅 |
| 🧪 全面测试 | 核心模块全覆盖，含 Mock 测试 |
| 🧠 RAG 实验模块 | 基于 chromem-go + Ollama 的本地知识库检索增强生成 |

## 📁 项目结构

```
agent_demo/
├── main.go                        # 主程序入口：组装各模块并启动 REPL 循环
├── go.mod / go.sum                # Go 模块定义与依赖校验
├── README.md                      # 项目说明文档
├── .gitignore                     # Git 忽略配置
│
├── internal/                      # 内部核心包
│   ├── agent/                     # Agent 核心调度器
│   │   ├── service.go             #   消息循环、工具调用调度（最多 20 轮）
│   │   └── service_test.go        #   单元测试（Mock LLM + Mock Executor）
│   │
│   ├── config/                    # 配置管理（环境变量）
│   │   ├── config.go              #   API_KEY / BASE_URL / MODEL_NAME 加载与校验
│   │   └── config_test.go         #   单元测试（正常加载 / 缺失报错）
│   │
│   ├── llm/                       # LLM 客户端
│   │   └── client.go              #   基于 go-openai 的流式 Chat Completion 封装
│   │
│   ├── session/                   # 对话会话管理
│   │   ├── session.go             #   消息追加与不可变副本返回
│   │   └── session_test.go        #   单元测试（追加 / 防篡改）
│   │
│   └── tool/                      # 工具系统
│       ├── registry.go            #   工具注册表：5 个工具的 OpenAI Tool 定义
│       ├── registry_test.go       #   单元测试（定义完整性）
│       ├── local.go               #   5 个工具函数的具体实现
│       ├── executor.go            #   工具执行器：JSON 参数解析与分发
│       └── executor_test.go       #   单元测试（读取文件 / 未知工具）
│
└── rag/                           # RAG 检索增强生成（独立实验模块）
    └── rag.go                     # 基于 chromem-go + Ollama 的本地 RAG 示例
```

## 🧱 核心架构详解

### 1️⃣ `internal/agent/` — Agent 核心调度器

```
用户输入 → Session.Append(用户消息)
         ↓
   ┌─ LLM.CompleteStream(消息列表, 工具列表) ──┐
   │         ↓ 流式输出 onContent              │
   │   返回 AI 响应消息                         │
   │         ↓                                 │
   │ 有 ToolCalls? ──是──→ 逐个执行工具         │
   │         │              ↓                  │
   │         否         Session.Append(结果)    │
   │         │              └──→ 继续循环       │
   │         ↓               (最多 20 轮)       │
   │   返回最终回答                             │
   └──────────────────────────────────────────┘
```

- **`Agent` 结构体**：封装 `LLMClient`、`ToolExecutor`、`Session`、工具列表
- **`Run()` 方法**：接收用户输入，循环调用 LLM 直到无工具调用或超过最大迭代次数（20 次）
- **最大工具迭代次数**：`maxToolIterations = 20`，防止无限循环

### 2️⃣ `internal/config/` — 配置管理

从环境变量加载配置，自动去除值两侧空白字符：

| 环境变量 | 必需 | 说明 | 默认值 |
|----------|------|------|--------|
| `API_KEY` | ✅ | AI API 密钥 | — |
| `MODEL_NAME` | ✅ | AI 模型名称 | — |
| `BASE_URL` | ❌ | API 基础 URL（需含协议头） | OpenAI 默认地址 |

### 3️⃣ `internal/llm/` — LLM 客户端

基于 `github.com/sashabaranov/go-openai` 库的流式 API 封装：

- **`CompleteStream()`**：流式接收 AI 响应，通过 `onContent` 回调实时输出
- **增量组装 Tool Calls**：流式模式下 Tool Calls 分块到达，客户端逐步拼接 `Function.Name` 和 `Function.Arguments`

### 4️⃣ `internal/session/` — 会话管理

- 维护 `[]openai.ChatCompletionMessage` 消息列表
- `Messages()` 返回副本，防止外部篡改（测试已验证）
- 支持用户消息、助手消息、工具调用结果消息三种类型

### 5️⃣ `internal/tool/` — 工具系统

注册了 5 个系统操作工具：

| 工具名称 | 功能 | 参数 |
|----------|------|------|
| `read_file` | 读取文本文件内容 | `path: string` |
| `list_file` | 列出目录下的文件和子目录 | `path: string` |
| `edit_file` | 在文件中查找并替换文本 | `path: string`, `old: string`, `new: string` |
| `create_file` | 创建新文件或覆盖已有文件 | `path: string`, `content: string` |
| `run_command` | 执行系统命令，返回 stdout/stderr/exit_code | `command: string`, `workdir?: string`, `args?: string[]` |

**架构分层**：
- **`registry.go`**：定义工具的 OpenAI Tool 结构（名称、描述、JSON Schema 参数）
- **`local.go`**：5 个工具函数的纯 Go 实现
- **`executor.go`**：根据工具名称解析 JSON 参数并分发到对应函数，统一错误处理

### 6️⃣ `rag/` — RAG 检索增强生成（实验模块）

独立的 RAG 示例，演示如何为 AI 接入本地知识库：

- **向量数据库**：基于 `chromem-go` 的持久化本地向量数据库（`./db` 目录）
- **嵌入模型**：Ollama 的 `bge-m3:567m` 模型
- **LLM 模型**：`qwen3:0.6b` 模型进行增强问答
- **系统提示**：通过 Go `html/template` 动态生成，根据是否有上下文决定提示策略
- **后处理**：自动去除模型回复中的 `＜think＞...＜/think＞` 思考过程标记

> ⚠️ 当前为独立实验模块，未集成到 Agent 主流程中。

## 🚀 快速开始

### 环境要求

- Go 1.25.6 或更高版本
- AI API 密钥（DeepSeek、OpenAI 等）

### 安装步骤

#### 1. 克隆项目

```bash
git clone <your-repo-url>
cd agent_demo
```

#### 2. 安装依赖

```bash
go mod tidy
```

#### 3. 配置 AI API

**使用 DeepSeek API：**
```bash
export API_KEY="your-deepseek-api-key"
export BASE_URL="https://api.deepseek.com/v1"
export MODEL_NAME="deepseek-chat"
```

**使用 OpenAI API：**
```bash
export API_KEY="your-openai-api-key"
export BASE_URL="https://api.openai.com/v1"
export MODEL_NAME="gpt-4o"
```

**使用本地 Ollama：**
```bash
# 安装并启动 Ollama
curl -fsSL https://ollama.com/install.sh | sh
ollama serve &
ollama pull deepseek-coder

# 配置环境变量（API_KEY 需为非空值）
export API_KEY="ollama"
export BASE_URL="http://localhost:11434/v1"
export MODEL_NAME="deepseek-coder"
```

#### 4. 运行程序

```bash
go run main.go
```

#### 5. 构建可执行文件（可选）

```bash
go build -o agent-demo
./agent-demo
```

### 运行测试

```bash
# 运行所有测试
go test -v ./...

# 运行特定模块测试
go test -v ./internal/agent/...
go test -v ./internal/config/...
go test -v ./internal/session/...
go test -v ./internal/tool/...
```

## 📖 使用示例

### 文件读取与目录浏览

```
You > 请读取当前目录下的所有文件
Agent > 正在调用工具...
Agent > [.git .gitignore go.mod go.sum main.go README.md]

You > 请读取 main.go 文件的内容
Agent > 正在调用工具...
Agent > package main
         import ...
```

### 文件创建与编辑

```
You > 创建一个名为 test.txt 的文件，内容为 "Hello World"
Agent > File created successfully.

You > 将 test.txt 中的 "World" 替换为 "AI Agent"
Agent > File edited successfully.
```

### 系统命令执行

```
You > 查看当前工作目录
Agent > {"stdout":"/home/user/projects/agent_demo\n","stderr":"","exit_code":0}

You > 查看 Go 版本
Agent > {"stdout":"go version go1.25.6 linux/amd64\n","stderr":"","exit_code":0}
```

### 退出程序

```
You > exit
Bye!
```

## 🧪 测试体系

### 测试覆盖总览

| 测试文件 | 测试内容 | 关键验证点 |
|----------|----------|-----------|
| `agent/service_test.go` | Agent 核心循环 | 无工具调用直接回答、多轮工具调用、流式内容回调、会话消息数 |
| `config/config_test.go` | 配置加载 | 正常加载（含空白修剪）、缺失 API_KEY 和 MODEL_NAME 报错 |
| `session/session_test.go` | 会话管理 | 消息追加正确性、返回副本防外部篡改 |
| `tool/executor_test.go` | 工具执行 | 读取文件内容、未知工具错误处理 |
| `tool/registry_test.go` | 工具注册 | 5 个工具定义完整性、Function 非空校验 |

### 测试设计特点

- **Mock 模式**：Agent 测试使用 `fakeLLM` 和 `fakeExecutor` 模拟 LLM 和工具执行器，不依赖真实 API
- **临时目录**：工具执行器测试使用 `t.TempDir()` 创建临时文件，测试后自动清理
- **环境变量隔离**：配置测试使用 `t.Setenv()`，仅在测试函数内生效

## 🛠️ 开发指南

### 添加新工具

只需修改 `internal/tool/` 包中的 **3 个位置** + **更新测试**：

#### 步骤 1：实现工具函数（`local.go`）

```go
func NewTool(param1 string, param2 int) (string, error) {
    // 工具函数实现
    return "result", nil
}
```

#### 步骤 2：注册工具定义（`registry.go`）

在 `defaultTools()` 函数中添加新的 Tool 定义：

```go
{
    Type: openai.ToolTypeFunction,
    Function: &openai.FunctionDefinition{
        Name:        "new_tool",
        Description: "新工具的功能描述",
        Strict:      false,
        Parameters: json.RawMessage(`{
            "type":"object",
            "properties": {
                "param1": {"type":"string"},
                "param2": {"type":"integer"}
            },
            "required": ["param1", "param2"]
        }`),
    },
}
```

#### 步骤 3：添加工具调用处理（`executor.go`）

在 `Execute()` 方法的 `switch` 中添加对应的处理逻辑：

```go
case "new_tool":
    var args struct {
        Param1 string `json:"param1"`
        Param2 int    `json:"param2"`
    }
    if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
        return "", fmt.Errorf("invalid arguments for %s: %w", name, err)
    }
    result, err = NewTool(args.Param1, args.Param2)
```

#### 步骤 4：更新注册表测试（`registry_test.go`）

在工具名称列表中追加新工具：

```go
for _, name := range []string{"read_file", "list_file", ..., "new_tool"} {
```

### 交叉编译

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o agent-demo-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o agent-demo.exe

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o agent-demo-macos
```

## ⚠️ 注意事项

### 安全性

1. **API 密钥安全**：请勿将 API 密钥提交到版本控制系统
2. **文件权限**：程序需要适当的文件系统权限
3. **命令执行**：`run_command` 工具可执行任意系统命令，生产环境中需谨慎使用

### 使用限制

1. **上下文长度**：长对话可能达到 AI 模型的上下文窗口限制
2. **工具调用深度**：最多支持 20 轮工具调用迭代
3. **文件大小**：大文件操作可能影响性能
4. **模型兼容性**：不同模型的工具调用能力可能不同，建议使用支持 Function Calling 的模型

### 兼容性

- **操作系统**：支持 Windows、Linux、macOS
- **AI 模型**：兼容任何 OpenAI 格式的 API
- **文件编码**：文件操作使用 UTF-8 编码

## 🔮 未来规划

### 计划功能

- [ ] **Web 界面**：添加 Web UI 支持
- [ ] **插件系统**：支持第三方工具插件
- [ ] **多会话管理**：支持多个独立会话
- [ ] **权限控制**：细粒度的工具权限管理
- [ ] **RAG 集成**：将 `rag/` 模块整合到 Agent 主流程中，实现知识库增强
- [ ] **上下文压缩**：长对话时压缩历史消息

### 性能优化

- [ ] **缓存机制**：减少重复 API 调用
- [ ] **异步执行**：支持并发工具执行
- [ ] **流式工具结果**：大文件读取时支持流式返回

## 📝 许可证

本项目采用 MIT 许可证。

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request 来改进这个项目！

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📋 版本历史

### v1.0.0（当前版本）
- ✅ 模块化重构：拆分为 agent / config / llm / session / tool 五个独立包
- ✅ 全面单元测试覆盖核心模块
- ✅ 基础文件操作工具（读取、列出、编辑、创建）
- ✅ 系统命令执行工具（返回 stdout / stderr / exit_code）
- ✅ 多 AI 后端支持（DeepSeek、OpenAI、Ollama）
- ✅ 流式响应与实时显示
- ✅ 对话历史管理（会话不可变拷贝）
- ✅ 工具调用循环（最多 20 次迭代）
- ✅ RAG 实验模块（基于 chromem-go + Ollama）
- ✅ 工具定义含参数描述（便于 LLM 理解参数含义）

---

**提示**：这是一个演示项目，展示了如何将 AI 能力与文件系统和系统命令操作结合。在实际使用中，请确保对 AI 的操作权限进行适当限制，特别是在生产环境中使用时。
