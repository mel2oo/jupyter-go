package jupyter

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
	MsgExecuteRequest = "execute_request"

	MsgError         = "error"
	MsgStream        = "stream"
	MsgStatus        = "status"
	MsgDisplayData   = "display_data"
	MsgExecuteInput  = "execute_input"
	MsgExecuteReply  = "execute_reply"
	MsgExecuteResult = "execute_result"
)

type JupyterRequestMessage struct {
	Header       MessageHeader  `json:"header"`
	ParentHeader MessageHeader  `json:"parent_header"`
	Metadata     map[string]any `json:"metadata"`
	Content      map[string]any `json:"content"`
	Buffers      []any          `json:"buffers"`
}

type JupyterResponseMessage struct {
	MsgID        string         `json:"msg_id"`
	MsgType      string         `json:"msg_type"`
	Header       MessageHeader  `json:"header"`
	ParentHeader MessageHeader  `json:"parent_header"`
	Metadata     map[string]any `json:"metadata"`
	Content      struct {
		Name           string         `json:"name,omitempty"`
		Text           string         `json:"text,omitempty"`
		Data           map[string]any `json:"data,omitempty"`
		Ename          string         `json:"ename,omitempty"`
		Evalue         string         `json:"evalue,omitempty"`
		Traceback      []string       `json:"traceback,omitempty"`
		ExecutionState string         `json:"execution_state,omitempty"`
		Status         string         `json:"status,omitempty"`
	} `json:"content,omitempty"`
	Buffers []any  `json:"buffers"`
	Channel string `json:"channel"`
}

type MessageHeader struct {
	MsgID    string `json:"msg_id"`
	MsgType  string `json:"msg_type"`
	Username string `json:"username"`
	Session  string `json:"session"`
	Date     string `json:"date"`
	Version  string `json:"version"`
}
