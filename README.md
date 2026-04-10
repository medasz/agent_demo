# AI Agent Demo - 智能系统操作助手

一个基于Go语言开发的命令行AI助手，能够通过自然语言指令操作文件系统和执行系统命令。

## 🚀 项目概述

这是一个演示项目，展示了如何将AI能力与系统操作工具结合，实现通过自然语言控制文件系统和执行命令的功能。项目使用Go语言开发，支持多种AI后端（DeepSeek、OpenAI、Ollama等）。

## ✨ 核心特性

### 🤖 智能对话
- **自然语言交互**：使用中文或英文与AI对话
- **上下文感知**：维护对话历史，支持多轮交互
- **意图识别**：自动识别用户需求并调用相应工具

### 🛠️ 工具系统
- **文件读取**：读取任意文本文件内容
- **目录浏览**：列出指定目录的文件和子目录
- **文件编辑**：在文件中查找并替换文本
- **文件创建**：创建新文件或覆盖现有文件
- **命令执行**：执行任意系统命令并获取结果

### 🔧 技术特性
- **模块化设计**：工具函数独立封装，易于扩展
- **多AI后端**：支持OpenAI兼容的API（DeepSeek、Ollama等）
- **流式响应**：实时显示AI思考过程
- **错误处理**：完善的错误处理和用户友好提示

## 📁 项目结构

```
agent_demo/
├── main.go          # 主程序文件，包含所有核心逻辑
├── go.mod           # Go模块配置文件
├── go.sum           # 依赖校验文件
├── README.md        # 项目说明文档
└── .gitignore       # Git忽略配置
```

### 主要文件说明

#### main.go
- **AI对话循环**：处理用户输入和AI响应的主循环
- **工具函数实现**：5个核心工具函数的Go实现
- **API客户端配置**：支持多种AI API后端
- **工具调用处理**：解析和执行AI的工具调用请求

#### 工具函数
1. **read_file** - 读取文件内容
2. **list_file** - 列出目录内容
3. **edit_file** - 编辑文件（文本替换）
4. **create_file** - 创建或覆盖文件
5. **run_command** - 执行系统命令

## 🚀 快速开始

### 环境要求
- Go 1.25.6 或更高版本
- AI API密钥（DeepSeek、OpenAI等）

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

#### 3. 配置AI API
根据您使用的AI服务选择相应的配置：

**使用DeepSeek API：**
```bash
export API_KEY="your-deepseek-api-key"
export BASE_URL="https://api.deepseek.com/v1"
```

**使用OpenAI API：**
```bash
export API_KEY="your-openai-api-key"
export BASE_URL="https://api.openai.com/v1"
```

**使用本地Ollama：**
```bash
# 安装并启动Ollama
curl -fsSL https://ollama.com/install.sh | sh
ollama serve &
ollama pull deepseek-coder

# 配置环境变量
export API_KEY="ollama"  # 任意值
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

## 📖 使用示例

### 基本文件操作
```
You > 请读取当前目录下的所有文件
Agent > 正在调用工具...
Agent > 当前目录下有: [.git .gitignore go.mod go.sum main.go README.md]

You > 请读取main.go文件的内容
Agent > 正在调用工具...
Agent > 这是main.go文件的内容...

You > 创建一个名为test.txt的文件，内容为"Hello World"
Agent > 正在调用工具...
Agent > File created successfully.

You > exit
Bye!
```

### 多轮对话和文件编辑
```
You > 我想修改README文件
Agent > 您想如何修改README文件？请告诉我具体要修改什么内容。

You > 在README中添加项目描述
Agent > 正在调用工具...
Agent > 已读取README.md文件内容

Agent > 我已经读取了README文件。您想在哪个位置添加项目描述？或者我可以帮您编辑特定内容。

You > 在文件开头添加"这是一个AI文件操作助手项目"
Agent > 正在调用工具...
Agent > File edited successfully.

Agent > 已成功在README文件开头添加了项目描述。
```

### 系统命令执行
```
You > 查看当前目录的文件列表
Agent > 正在调用工具...
Agent > {"stdout":"[.git .gitignore go.mod go.sum main.go README.md]\r\n","stderr":"","exit_code":0}

Agent > 当前目录下有: .git, .gitignore, go.mod, go.sum, main.go, README.md
```

## 🔧 配置说明

### 必需环境变量
- `API_KEY`: AI API密钥
- `BASE_URL`: API基础URL（必须包含http://或https://）

### 可选环境变量
- `MODEL_NAME`: AI模型名称（默认为API默认模型）

### 依赖管理
项目使用Go Modules进行依赖管理：
```bash
# 安装依赖
go mod tidy

# 查看依赖
go list -m all
```

## 🛠️ 开发指南

### 添加新工具
要添加新的工具函数，需要修改以下部分：

#### 1. 实现工具函数
在`main.go`中添加新的工具函数实现：
```go
func newToolFunction(param1 string, param2 int) (string, error) {
    // 工具函数实现
    return "result", nil
}
```

#### 2. 注册工具定义
在`tools`数组中添加新的工具定义：
```go
{
    Type: openai.ToolTypeFunction,
    Function: &openai.FunctionDefinition{
        Name:        "new_tool",
        Description: "Description of the new tool",
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

#### 3. 添加工具调用处理
在`toolCall switch`语句中添加对应的处理逻辑：
```go
case "new_tool":
    var args struct {
        Param1 string `json:"param1"`
        Param2 int    `json:"param2"`
    }
    json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
    result, err = newToolFunction(args.Param1, args.Param2)
```

### 构建与部署
```bash
# 构建可执行文件
go build -o agent-demo

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o agent-demo-linux

# 交叉编译（Windows）
GOOS=windows GOARCH=amd64 go build -o agent-demo.exe

# 交叉编译（macOS）
GOOS=darwin GOARCH=arm64 go build -o agent-demo-macos
```

## ⚠️ 注意事项

### 安全性
1. **API密钥安全**：请勿将API密钥提交到版本控制系统
2. **文件权限**：程序需要适当的文件系统权限
3. **命令执行**：`run_command`工具可以执行任意系统命令，使用时需谨慎

### 使用限制
1. **上下文长度**：对话历史可能达到AI模型的上下文限制
2. **工具调用深度**：支持多轮工具调用，但可能受AI模型限制
3. **文件大小**：大文件操作可能影响性能

### 兼容性
1. **操作系统**：支持Windows、Linux、macOS
2. **AI模型**：支持多种AI模型，但不同模型的工具调用能力可能不同
3. **文件编码**：文件操作使用UTF-8编码

## 📊 技术架构

### 核心组件
1. **AI客户端**：基于go-openai库，支持OpenAI兼容的API
2. **工具系统**：模块化的工具函数，每个工具独立实现
3. **对话管理器**：维护对话历史，支持多轮交互
4. **参数解析器**：自动解析JSON格式的工具参数

### 通信流程
```
用户输入 → AI分析 → 工具调用 → 执行操作 → 结果反馈 → AI响应 → 用户
```

### 扩展性设计
- **工具热插拔**：新工具只需在三个位置添加代码
- **多AI后端**：支持任何OpenAI兼容的API
- **跨平台**：Go语言原生跨平台支持

## 🎯 应用场景

### 开发辅助
- 自动化文件操作和代码管理
- 系统信息查询和配置管理
- 项目结构分析和文档生成

### 运维工具
- 服务器状态监控
- 日志文件分析
- 系统配置管理

### 个人助手
- 文件整理和分类
- 自动化任务执行
- 信息查询和整理

## 🔮 未来规划

### 计划功能
1. **Web界面**：添加Web UI支持
2. **插件系统**：支持第三方工具插件
3. **会话管理**：支持多个独立会话
4. **权限控制**：细粒度的工具权限管理

### 性能优化
1. **缓存机制**：减少重复API调用
2. **异步处理**：支持并发工具执行
3. **流式响应**：实时显示AI响应

## 📝 许可证

[添加您的许可证信息]

## 🤝 贡献指南

欢迎提交Issue和Pull Request来改进这个项目！

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 📞 联系方式

[添加您的联系方式或项目链接]

---

**提示**：这是一个演示项目，展示了如何将AI能力与文件系统和系统命令操作结合。在实际使用中，请确保对AI的操作权限进行适当限制，特别是在生产环境中使用时。

## 📋 版本历史

### v1.0.0 (当前版本)
- 基础文件操作工具（读取、列出、编辑、创建）
- 系统命令执行工具
- 多AI后端支持（DeepSeek、OpenAI、Ollama）
- 对话历史管理
- 错误处理和用户友好界面

### 后续版本计划
- v1.1.0：添加Web界面支持
- v1.2.0：实现插件系统
- v2.0.0：重构为微服务架构