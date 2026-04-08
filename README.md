# Agent Demo - 基于DeepSeek的AI文件操作助手

一个使用Go语言开发的命令行AI助手，能够通过自然语言指令操作文件系统。

## 功能特性

- 🤖 **智能对话**：基于DeepSeek AI模型的自然语言交互
- 📁 **文件操作**：支持读取、列出、编辑和创建文件
- 🔧 **工具调用**：AI自动调用工具函数完成用户请求
- ⚡ **轻量高效**：Go语言编写，运行速度快
- 🔧 **可扩展**：易于添加新的工具函数

## 工具函数

| 工具名称 | 功能描述 | 参数 |
|---------|---------|------|
| `read_file` | 读取文件内容 | `path`: 文件路径 |
| `list_file` | 列出目录中的文件 | `path`: 目录路径 |
| `edit_file` | 编辑文件（替换文本） | `path`: 文件路径, `old`: 要替换的文本, `new`: 新文本 |
| `crate_file` | 创建或覆盖文件 | `path`: 文件路径, `content`: 文件内容 |

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

## 项目结构

```
.
├── main.go          # 主程序文件
├── go.mod           # Go模块配置文件
├── go.sum           # 依赖校验文件
├── .gitignore       # Git忽略配置
└── README.md        # 项目说明文档
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
4. **工具拼写**：注意`crate_file`工具名称的拼写（应为`create_file`但保持原样）

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