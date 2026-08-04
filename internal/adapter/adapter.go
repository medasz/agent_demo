package adapter

// Message 表示来自某个渠道的用户输入
type Message struct {
	SessionID string          // 渠道的会话ID
	Content   string          // 用户输入的文本
	ReplyChan chan<- ReplyChunk // 用于将Agent的回复推回给该渠道
}

// ReplyChunk 表示推回给渠道的回复分片或状态变更
type ReplyChunk struct {
	IsThinking bool   // 是否正在思考
	Chunk      string // 流式文本片段
	IsDone     bool   // 当前回复是否结束
	Error      error  // 发生的错误
}

// Channel 表示一个输入输出渠道适配器（如终端、桌宠、微信等）
type Channel interface {
	// Start 启动该渠道，通常这是一个阻塞调用，对于 Ebiten 来说必须在主线程运行
	Start(inputBus chan<- Message) error
}
