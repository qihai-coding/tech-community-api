package models

import "time"

// CodeSnippet 代码片段结构体
type CodeSnippet struct {
	ID          uint      `json:"id" db:"id"`
	UserID      uint      `json:"user_id" db:"user_id"`
	Title       string    `json:"title" db:"title"`
	Language    string    `json:"language" db:"language"`
	Code        string    `json:"code" db:"code"`
	Description string    `json:"description" db:"description"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	ShareToken  *string   `json:"share_token,omitempty" db:"share_token"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CodeExecution 代码执行记录结构体
type CodeExecution struct {
	ID            uint      `json:"id" db:"id"`
	SnippetID     *uint     `json:"snippet_id,omitempty" db:"snippet_id"`
	UserID        uint      `json:"user_id" db:"user_id"`
	Language      string    `json:"language" db:"language"`
	Code          string    `json:"code" db:"code"`
	Stdin         string    `json:"stdin" db:"stdin"`
	Output        string    `json:"output" db:"output"`
	Error         string    `json:"error" db:"error"`
	ExecutionTime *int      `json:"execution_time,omitempty" db:"execution_time"` // 毫秒
	MemoryUsage   *int64    `json:"memory_usage,omitempty" db:"memory_usage"`     // 字节
	Status        string    `json:"status" db:"status"`                           // success, error, timeout
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// CodeCollaboration 代码协作会话结构体
type CodeCollaboration struct {
	ID           uint      `json:"id" db:"id"`
	SnippetID    uint      `json:"snippet_id" db:"snippet_id"`
	SessionToken string    `json:"session_token" db:"session_token"`
	ActiveUsers  string    `json:"active_users" db:"active_users"` // JSON 字符串
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
}

// ExecuteCodeRequest 执行代码请求
type ExecuteCodeRequest struct {
	Language string `json:"language" binding:"required"`
	Code     string `json:"code" binding:"required"`
	Stdin    string `json:"stdin"`
	SaveAs   string `json:"save_as"` // 可选：保存代码片段的标题
}

// ExecuteCodeResponse 执行代码响应
type ExecuteCodeResponse struct {
	Output        string `json:"output"`
	Error         string `json:"error,omitempty"`
	ExecutionTime int    `json:"execution_time"` // 毫秒
	MemoryUsage   int64  `json:"memory_usage"`   // 字节
	Status        string `json:"status"`         // success, error, timeout
	SnippetID     *uint  `json:"snippet_id,omitempty"`
}

// SaveSnippetRequest 保存代码片段请求
type SaveSnippetRequest struct {
	Title       string `json:"title" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// UpdateSnippetRequest 更新代码片段请求
type UpdateSnippetRequest struct {
	Title       string `json:"title"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsPublic    *bool  `json:"is_public"`
}

// ShareSnippetResponse 分享代码片段响应
type ShareSnippetResponse struct {
	ShareToken string `json:"share_token"`
	ShareURL   string `json:"share_url"`
}

// LanguageInfo 支持的语言信息
type LanguageInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	PistonName  string `json:"piston_name"`  // 兼容前端字段，值与语言 ID 保持一致
	DefaultCode string `json:"default_code"` // 默认代码模板
	Judge0ID    int    `json:"-"`            // Judge0 语言 ID（仅后端使用）
}

// Judge0SubmissionRequest Judge0 提交请求
type Judge0SubmissionRequest struct {
	SourceCode    string  `json:"source_code"`
	LanguageID    int     `json:"language_id"`
	Stdin         string  `json:"stdin,omitempty"`
	CPUTimeLimit  float64 `json:"cpu_time_limit,omitempty"`
	WallTimeLimit float64 `json:"wall_time_limit,omitempty"`
	MemoryLimit   int     `json:"memory_limit,omitempty"` // KB
}

// Judge0SubmissionResponse Judge0 提交响应
type Judge0SubmissionResponse struct {
	Stdout        *string `json:"stdout"`
	Stderr        *string `json:"stderr"`
	CompileOutput *string `json:"compile_output"`
	Message       *string `json:"message"`
	Error         *string `json:"error"`
	Time          string  `json:"time"`
	Memory        int64   `json:"memory"` // KB
	Token         string  `json:"token"`
	Status        struct {
		ID          int    `json:"id"`
		Description string `json:"description"`
	} `json:"status"`
}

// Judge0ErrorResponse Judge0 API 错误响应
type Judge0ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

// CodeSnippetListItem 代码片段列表项（简化版）
type CodeSnippetListItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Language  string    `json:"language"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CodeSnippetWithUser 代码片段及用户信息
type CodeSnippetWithUser struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	Language    string    `json:"language"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	ShareToken  *string   `json:"share_token,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
