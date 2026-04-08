# Agent Demo - 基于DeepSeek的AI文件操作助手

一个使用Go语言开发的命令行AI助手，能够通过自然语言指令操作文件系统。

## 功能特性

### 核心功能
- 🤖 **智能对话**：基于DeepSeek AI模型的自然语言交互
- 📁 **文件操作**：支持读取、列出、编辑和创建文件
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

### 技术特性
- 🏗️ **模块化设计**：工具函数独立封装，便于维护和扩展
- 🔌 **API集成**：无缝集成DeepSeek API，支持自定义API端点
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

## 快速开始

### 环境要求

- Go 1.25.6 或更高版本
- DeepSeek API密钥

### 安装步骤

1. **克隆项目**
   ```bash
   git clone <your-repo-url>
   cd agent_demo
   ```

2. **设置环境变量**
   ```bash
   # 设置DeepSeek API密钥
   export API_KEY="your-deepseek-api-key"
   
   # 设置API基础URL（可选，默认为DeepSeek官方API）
   export BASE_URL="https://api.deepseek.com"
   ```

3. **运行程序**
   ```bash
   go run main.go
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

#### 示例4：错误处理示例
```
You > 读取一个不存在的文件
Agent > 正在调用工具...
Agent > open nonexistent.txt: no such file or directory

Agent > 抱歉，文件不存在。请检查文件路径是否正确。
```

## 项目结构

```
.
├── main.go          # 主程序文件
├── go.mod           # Go模块配置文件
├── go.sum           # 依赖校验文件
├── .gitignore       # Git忽略配置
├── README.md        # 项目说明文档
├── agent_demo.exe   # Windows可执行文件
├── .git/            # Git版本控制目录
└── .idea/           # IDE配置文件目录
```

## 配置说明

### API配置

程序使用以下环境变量：

- `API_KEY`: **必需** - DeepSeek API密钥
- `BASE_URL`: **可选** - API基础URL，默认为DeepSeek官方API

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

1. **在main.go中添加工具函数实现**
2. **在CreateChatCompletionRequest的Tools数组中注册新工具**
3. **在toolCall switch语句中添加对应的处理逻辑**

### 构建与部署

```bash
# 构建可执行文件
go build -o agent-demo

# 交叉编译（Linux）
GOOS=linux GOARCH=amd64 go build -o agent-demo-linux

# 交叉编译（Windows）
GOOS=windows GOARCH=amd64 go build -o agent-demo.exe
```

## 注意事项

1. **API密钥安全**：请勿将API密钥提交到版本控制系统
2. **文件权限**：程序需要适当的文件系统权限来操作文件
3. **错误处理**：程序包含基本的错误处理，但建议在生产环境中添加更完善的错误处理机制
4. **工具拼写**：已修正`create_file`工具名称的拼写

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

**提示**：这是一个演示项目，展示了如何将AI能力与文件系统操作结合。在实际使用中，请确保对AI的操作权限进行适当限制。