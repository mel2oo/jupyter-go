package jupyter

import "github.com/google/uuid"

const (
	SessionNotebook = "notebook"
	SessionConsole  = "console"
	SessionTerminal = "terminal"
	SessionFile     = "file"
)

type Session struct {
	ID       string   `json:"id,omitempty"`
	Path     string   `json:"path,omitempty"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type,omitempty"`
	Kernel   Kernel   `json:"kernel,omitempty"`
	Notebook Notebook `json:"notebook,omitempty"`
}

type Kernel struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	LastActivity   string  `json:"last_activity,omitempty"`
	ExecutionState string  `json:"execution_state,omitempty"`
	Connections    float64 `json:"connections,omitempty"`
}

type Notebook struct {
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

// ws message

const (
	MsgExecuteRequest    = "execute_request"
	MsgCompleteRequest   = "complete_request"
	MsgKernelInfoRequest = "kernel_info_request"
)

type JupyterMessage struct {
	Header       MessageHeader  `json:"header"`
	ParentHeader map[string]any `json:"parent_header"`
	Metadata     map[string]any `json:"metadata"`
	Content      map[string]any `json:"content"`
}

type MessageHeader struct {
	Date     string `json:"date,omitempty"`
	MsgID    string `json:"msg_id"`
	MsgType  string `json:"msg_type"`
	Session  string `json:"session"`
	Username string `json:"username"`
	Version  string `json:"version"`
}

func NewMessageTmpl(sessionID string) *JupyterMessage {
	return &JupyterMessage{
		Header: MessageHeader{
			MsgID:    uuid.New().String(),
			Session:  sessionID,
			Username: "go-client",
			Version:  "5.0",
		},
		ParentHeader: map[string]any{},
		Metadata:     map[string]any{},
		Content:      map[string]any{},
	}
}
