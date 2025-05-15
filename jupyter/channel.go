package jupyter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Channel struct {
	ws *websocket.Conn
}

func (s *SessionService) Connect(ctx context.Context, kernelID, sessionID string) (*Channel, error) {
	conn, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://localhost:8888/api/kernels/%s/channels?session_id=%s",
			kernelID, sessionID), nil)
	if err != nil {
		return nil, err
	}

	return &Channel{conn}, nil
}

func (c *Channel) Close() error {
	return c.ws.Close()
}

func (c *Channel) Recv() error {
	for {
		_, raw, err := c.ws.ReadMessage()
		if err != nil {
			return err
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(raw, &msg); err != nil {
			fmt.Println("json error:", err)
			continue
		}

		msgType := msg["msg_type"]
		content := msg["content"]

		if msgType == "stream" {
			c := content.(map[string]interface{})
			fmt.Print("Print:", c["text"])
		} else if msgType == "execute_result" {
			c := content.(map[string]interface{})
			fmt.Println("Result:", c["data"])
		} else if msgType == "error" {
			c := content.(map[string]interface{})
			fmt.Println("Error:", c["evalue"])
		}
	}
}

// code execute
// content:
// code	string 要执行的代码，如 "1 + 1"
// silent bool 若为 true，抑制 iopub 输出，仍有 execute_reply
// store_history bool 是否存入历史记录，默认 true
// user_expressions	object 用户定义表达式（key -> expr），将执行并返回
// allow_stdin bool 是否允许请求 stdin 输入
// stop_on_error bool 若为 true，出错即终止执行（默认 false）
func (c *Channel) CodeExecute(ctx context.Context, sessionID string, code string) error {
	msg := NewMessageTmpl(sessionID)
	msg.Header.MsgType = MsgExecuteRequest
	msg.Content["code"] = code
	msg.Content["slient"] = true
	msg.Content["store_history"] = false

	return c.ws.WriteJSON(msg)
}
