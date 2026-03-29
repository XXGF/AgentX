package conversation

// Role 消息角色枚举（对应 Java 的 Role）
type Role string

const (
	RoleUser      Role = "USER"
	RoleSystem    Role = "SYSTEM"
	RoleAssistant Role = "ASSISTANT"
	RoleSummary   Role = "SUMMARY" // 只存在历史消息表中，实际发送给大模型时会被转换
)

// MessageType 消息类型枚举（对应 Java 的 MessageType）
type MessageType string

const (
	MessageTypeText              MessageType = "TEXT"
	MessageTypeToolCall          MessageType = "TOOL_CALL"
	MessageTypeTaskExec          MessageType = "TASK_EXEC"
	MessageTypeTaskStatusLoading MessageType = "TASK_STATUS_TO_LOADING"
	MessageTypeTaskStatusFinish  MessageType = "TASK_STATUS_TO_FINISH"
	MessageTypeTaskSplitFinish   MessageType = "TASK_SPLIT_FINISH"
	MessageTypeRagRetrievalStart    MessageType = "RAG_RETRIEVAL_START"
	MessageTypeRagRetrievalProgress MessageType = "RAG_RETRIEVAL_PROGRESS"
	MessageTypeRagRetrievalEnd      MessageType = "RAG_RETRIEVAL_END"
	MessageTypeRagThinkingStart     MessageType = "RAG_THINKING_START"
	MessageTypeRagThinkingProgress  MessageType = "RAG_THINKING_PROGRESS"
	MessageTypeRagThinkingEnd       MessageType = "RAG_THINKING_END"
	MessageTypeRagAnswerStart       MessageType = "RAG_ANSWER_START"
	MessageTypeRagAnswerProgress    MessageType = "RAG_ANSWER_PROGRESS"
	MessageTypeRagAnswerEnd         MessageType = "RAG_ANSWER_END"
)
