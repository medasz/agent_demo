# Agent Demo - 基于AI的系统操作助手

一个使用Go语言开发的命令行AI助手，能够通过自然语言指令操作文件系统和执行系统命令。

## 功能特性

### 核心功能
- 🤖 **智能对话**：基于AI模型的自然语言交互（支持DeepSeek、OpenAI、Ollama等）
- 📁 **文件操作**：支持读取、列出、编辑和创建文件
- 💻 **系统命令**：支持执行任意系统命令并获取结果
- 🔧 **工具调用**：AI自动调用工具函数完成用户请求
- ⚡ **轻量高效**：Go语言编写，运行速度快
- 🔧 **可扩展**：易于添加新的工具函数

### 新增功能
- 🔄 **对话历史管理**：自动维护对话上下文，支持多轮对话
- 🧠 **智能工具选择**：AI自动判断何时需要调用工具以及调用哪个工具
- 🛡️ **错误处理机制**：内置基本的错误处理和异常捕获
- 💬 **交互式对话**：自然的命令行交互界面，支持退出命令
- 🔄 **工具结果反馈**：工具执行结果自动反馈给AI进行下一步决策
- 📝 **文件内容替换**：支持在文件中查找并替换特定文本内容
- 📂 **目录浏览**：支持列出指定目录下的所有文件和子目录
- 🖥️ **系统命令执行**：支持执行任意系统命令并返回标准输出、标准错误和退出码

### 技术特性
- 🏗️ **模块化设计**：工具函数独立封装，便于维护和扩展
- 🔌 **API集成**：无缝集成多种AI API，支持自定义API端点
- 📊 **JSON参数解析**：自动解析工具调用的JSON参数
- 🔄 **循环对话处理**：支持连续的工具调用和对话交互
- 🎯 **精确参数验证**：工具调用时进行参数验证和错误检查

## 功能详细说明

### 1. 智能对话系统
- **上下文感知**：AI能够记住之前的对话内容，实现连贯的多轮对话
- **意图识别**：自动识别用户意图，判断是否需要调用工具
- **自然语言理解**：支持中文和英文的自然语言指令

### 2. 文件操作功能
- **文件读取**：读取任意文本文件的内容，支持相对路径和绝对路径
- **目录浏览**：列出指定目录下的所有文件和子目录
- **文件编辑**：在文件中查找并替换特定文本内容
- **文件创建**：创建新文件或覆盖现有文件内容

### 3. 工具调用机制
- **自动工具选择**：AI根据用户请求自动选择合适的工具
- **参数自动提取**：从用户指令中自动提取工具所需参数
- **结果智能处理**：工具执行结果自动整合到对话流中
- **错误自动处理**：工具执行失败时自动提供错误信息

### 4. 系统特性
- **环境变量配置**：通过环境变量配置API密钥和端点
- **跨平台支持**：支持Windows、Linux、macOS等操作系统
- **命令行交互**：简洁的命令行界面，支持退出命令
- **依赖管理**：使用Go Modules进行依赖管理

## 工具函数

| 工具名称 | 功能描述 | 参数 |
|---------|---------|------|
| `read_file` | 读取文件内容 | `path`: 文件路径 |
| `list_file` | 列出目录中的文件 | `path`: 目录路径 |
| `edit_file` | 编辑文件（替换文本） | `path`: 文件路径, `old`: 要替换的文本, `new`: 新文本 |
| `create_file` | 创建或覆盖文件 | `path`: 文件路径, `content`: 文件内容 |
| `run_command` | 执行系统命令 | `command`: 命令名称, `args`: 参数数组, `workdir`: 工作目录（可选） |

## 快速开始

### 环境要求

- Go 1.25.6 或更高版本
- DeepSeek API密钥

### 安装步骤

#### 1. 克隆项目
```bash
git clone <your-repo-url>
cd agent_demo
```

#### 2. 安装Go依赖
```bash
go mod tidy
```

#### 3. 配置AI API
根据您使用的AI服务选择相应的配置：

##### 选项A：使用DeepSeek API
```bash
export API_KEY="your-deepseek-api-key"
export BASE_URL="https://api.deepseek.com/v1"
# export MODEL_NAME="deepseek-chat"  # 可选，指定模型名称
```

##### 选项B：使用OpenAI API
```bash
export API_KEY="your-openai-api-key"
export BASE_URL="https://api.openai.com/v1"
# export MODEL_NAME="gpt-4"  # 可选，指定模型名称
```

##### 选项C：使用本地Ollama
```bash
# 首先安装并启动Ollama
curl -fsSL https://ollama.com/install.sh | sh
ollama serve &
ollama pull deepseek-coder  # 或其他支持的模型

# 然后配置环境变量
export API_KEY="ollama"  # 可以是任意值
export BASE_URL="http://localhost:11434/v1"
export MODEL_NAME="deepseek-coder"  # 指定使用的模型
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

### 故障排除

#### 常见错误

##### 错误1：`unsupported protocol scheme ""`
**问题描述**：
```
CreateChatCompletion error: Post "/chat/completions": unsupported protocol scheme ""
```

**原因**：`BASE_URL` 环境变量未设置或缺少协议头
**解决方案**：
```bash
# 确保BASE_URL包含协议头
export BASE_URL="https://api.deepseek.com/v1"  # 正确
# export BASE_URL="api.deepseek.com/v1"        # 错误，缺少协议头
```

##### 错误2：`API_KEY is required`
**问题描述**：API密钥未设置
**解决方案**：
```bash
export API_KEY="your-api-key-here"
```

##### 错误3：连接Ollama失败
**问题描述**：无法连接到本地Ollama服务
**解决方案**：
```bash
# 确保Ollama服务正在运行
ollama serve &

# 检查服务状态
curl http://localhost:11434/api/tags
```

### 使用示例

启动程序后，您可以通过自然语言与AI交互：

#### 示例1：基本文件操作
```
You > 请读取当前目录下的所有文件
Agent > 正在调用工具...
Agent > 当前目录下有: [.git .gitignore go.mod go.sum main.go README.md agent_demo.exe]

You > 请读取main.go文件的内容
Agent > 正在调用工具...
Agent > 这是main.go文件的内容...

You > 创建一个名为test.txt的文件，内容为"Hello World"
Agent > 正在调用工具...
Agent > File created successfully.

You > exit
Bye!
```

#### 示例2：多轮对话和文件编辑
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

#### 示例3：复杂操作组合
```
You > 查看项目结构，然后创建一个docs目录
Agent > 正在调用工具...
Agent > 当前目录下有: [.git .gitignore go.mod go.sum main.go README.md agent_demo.exe]

You > 在docs目录中创建一个使用指南文件
Agent > 正在调用工具...
Agent > File created successfully.

Agent > 已创建docs目录并在其中创建了使用指南文件。
```

#### 示例4：系统命令执行
```
You > 查看当前目录的文件列表
Agent > 正在调用工具...
Agent > {"stdout":"[.git .gitignore go.mod go.sum main.go README.md agent_demo.exe]\r\n","stderr":"","exit_code":0}

Agent > 当前目录下有: .git, .gitignore, go.mod, go.sum, main.go, README.md, agent_demo.exe

You > 查看系统信息
Agent > 正在调用工具...
Agent > {"stdout":"\r\nWindows IP \ufffd\ufffd\ufffd\ufffd\r\n\r\n   \ufffd\ufffd\ufffd\ufffd\ufffd\ufffd  . . . . . . . . . . . . . : SHAG-ZJH296-PC\r\n   \ufffd\ufffd DNS \ufffd\ufffd׺ . . . . . . . . . . . : angeek.com.cn\r\n   \ufffdڵ\ufffd\ufffd\ufffd\ufffd\ufffd  . . . . . . . . . . . . : \ufffd\ufffd\ufffd\ufffd\r\n   IP ·\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd . . . . . . . . . . : \ufffd\ufffd\r\n   WINS \ufffd\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd\ufffd . . . . . . . . . : \ufffd\ufffd\r\n   DNS \ufffd\ufffd׺\ufffd\ufffd\ufffd\ufffd\ufffdб\ufffd  . . . . . . . . : angeek.com.cn\r\n\r\n","stderr":"","exit_code":0}

Agent > 系统信息已获取，主机名为SHAG-ZJH296-PC，域名为angeek.com.cn。

#### 示例5：错误处理示例
```
You > 读取一个不存在的文件
Agent > 正在调用工具...
Agent > open nonexistent.txt: no such file or directory

Agent > 抱歉，文件不存在。请检查文件路径是否正确。
```

## 项目结构

```
.
├── main.go          # 主程序文件 - 包含AI助手核心逻辑和工具函数
├── go.mod           # Go模块配置文件 - 定义项目依赖
├── go.sum           # 依赖校验文件 - 确保依赖版本一致性
├── .gitignore       # Git忽略配置 - 忽略IDE配置和可执行文件
├── README.md        # 项目说明文档 - 使用指南和功能说明
├── .git/            # Git版本控制目录
└── .idea/           # IDE配置文件目录（可选）

### 主要文件说明

#### main.go
- **AI对话循环**：处理用户输入和AI响应的主循环
- **工具函数定义**：包含5个核心工具函数的实现
- **API客户端配置**：支持多种AI API后端
- **工具调用处理**：解析AI的工具调用请求并执行相应操作

#### 工具函数
1. **read_file** - 读取文件内容
2. **list_file** - 列出目录内容
3. **edit_file** - 编辑文件（文本替换）
4. **create_file** - 创建或覆盖文件
5. **run_command** - 执行系统命令

## 配置说明

### API配置

程序使用以下环境变量：

- `API_KEY`: **必需** - AI API密钥（DeepSeek、OpenAI或Ollama）
- `BASE_URL`: **必需** - API基础URL，必须包含协议头（http://或https://）
- `MODEL_NAME`: **可选** - AI模型名称，默认为空（使用API默认模型）

### 依赖管理

项目使用Go Modules进行依赖管理：

```bash
# 安装依赖
go mod tidy

# 查看依赖
go list -m all
```

## 开发指南

### 添加新工具

要添加新的工具函数，需要修改以下部分：

#### 1. 实现工具函数
在`main.go`中添加新的工具函数实现，例如：
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

### 代码结构说明

#### 主循环流程
1. 读取用户输入
2. 将用户消息添加到对话历史
3. 调用AI API获取响应
4. 如果AI返回工具调用，执行相应工具
5. 将工具执行结果反馈给AI
6. 重复步骤3-5直到AI返回最终回答
7. 显示AI回答并等待下一个用户输入

#### 工具调用机制
- 每个工具都有明确定义的JSON Schema参数
- AI自动从用户指令中提取参数
- 工具执行结果以JSON格式返回
- 错误信息也会反馈给AI进行下一步决策

## 注意事项

### 安全性
1. **API密钥安全**：请勿将API密钥提交到版本控制系统，使用环境变量管理
2. **文件权限**：程序需要适当的文件系统权限来操作文件，注意权限控制
3. **命令执行**：`run_command`工具可以执行任意系统命令，使用时需谨慎

### 使用限制
1. **上下文长度**：对话历史会不断累积，可能达到AI模型的上下文限制
2. **工具调用深度**：支持多轮工具调用，但可能受AI模型限制
3. **文件大小**：大文件操作可能影响性能

### 兼容性
1. **操作系统**：支持Windows、Linux、macOS，但某些系统命令可能不兼容
2. **AI模型**：支持多种AI模型，但不同模型的工具调用能力可能不同
3. **文件编码**：文件操作使用UTF-8编码，其他编码可能有问题

### 性能考虑
1. **API延迟**：云端API调用可能有网络延迟
2. **本地模型**：使用本地Ollama时，模型加载和推理需要足够内存
3. **并发处理**：当前版本为单线程，不支持并发用户

## 许可证

[添加您的许可证信息]

## 贡献指南

欢迎提交Issue和Pull Request来改进这个项目！

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 联系方式

[添加您的联系方式或项目链接]

---

## 技术架构

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

## 应用场景

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

## 未来规划

### 计划功能
1. **Web界面**：添加Web UI支持
2. **插件系统**：支持第三方工具插件
3. **会话管理**：支持多个独立会话
4. **权限控制**：细粒度的工具权限管理
5. **批量操作**：支持批量文件处理

### 性能优化
1. **缓存机制**：减少重复API调用
2. **异步处理**：支持并发工具执行
3. **流式响应**：实时显示AI响应

---

**提示**：这是一个演示项目，展示了如何将AI能力与文件系统和系统命令操作结合。在实际使用中，请确保对AI的操作权限进行适当限制，特别是在生产环境中使用时。

## 版本历史

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