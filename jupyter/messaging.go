package jupyter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Channel struct {
	ws         *websocket.Conn
	kernelID   string
	sessionID  string
	executions map[string]*Exection
}

func (s *SessionService) Connect(ctx context.Context, kernelID, sessionID string) (*Channel, error) {
	ws, _, err := websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://localhost:8888/api/kernels/%s/channels?session_id=%s",
			kernelID, sessionID), nil)
	if err != nil {
		return nil, err
	}

	c := &Channel{
		ws:         ws,
		kernelID:   kernelID,
		sessionID:  sessionID,
		executions: make(map[string]*Exection),
	}

	return c, c.receiveMessage()
}

func (c *Channel) Close() error {
	return c.ws.Close()
}

// code execute
func (c *Channel) CodeExecute(ctx context.Context, code string) ([]any, error) {
	msgID := uuid.NewString()
	exection := &Exection{queue: make(chan any, 10)}

	c.executions[msgID] = exection
	defer delete(c.executions, msgID)

	if err := c.ws.WriteJSON(c.newExecuteRequest(msgID, code)); err != nil {
		return nil, err
	}

	res := make([]any, 0)
	for {
		val, err := exection.Recv()
		if err != nil || err == io.EOF {
			return res, nil
		}

		res = append(res, val)
	}
}

func (c *Channel) CodeExecuteStream(ctx context.Context, code string) (*Exection, error) {
	msgID := uuid.NewString()
	exection := &Exection{queue: make(chan any, 10)}

	c.executions[msgID] = exection
	defer delete(c.executions, msgID)

	if err := c.ws.WriteJSON(c.newExecuteRequest(msgID, code)); err != nil {
		return nil, err
	}

	return exection, nil
}

func (c *Channel) newExecuteRequest(msgID, code string) *JupyterRequestMessage {
	return &JupyterRequestMessage{
		Header: MessageHeader{
			MsgID:    msgID,
			MsgType:  MsgExecuteRequest,
			Username: "go-client",
			Session:  c.sessionID,
			Version:  "5.3",
			Date:     time.Now().UTC().Format(time.RFC3339),
		},
		ParentHeader: MessageHeader{},
		Metadata: map[string]any{
			"trusted":      true,
			"deletedCells": []any{},
			"recordTiming": false,
			"cellId":       uuid.NewString(),
		},
		Content: map[string]any{
			"code":             code,
			"slient":           true,
			"store_history":    false,
			"user_expressions": map[string]any{},
			"allow_stdin":      false,
			"stop_on_error":    true,
		},
	}
}

func (c *Channel) receiveMessage() error {
	go func() {
		for {
			_, raw, err := c.ws.ReadMessage()
			if err != nil {
				return
			}

			var msg JupyterResponseMessage
			if err := json.Unmarshal(raw, &msg); err != nil {
				continue
			}

			c.processMessage(&msg)
		}
	}()

	return nil
}

func (c *Channel) processMessage(msg *JupyterResponseMessage) {
	execution, ok := c.executions[msg.ParentHeader.MsgID]
	if !ok {
		return
	}

	switch msg.MsgType {
	case MsgError:
		if execution.errored {
			return
		}

		name, _ := msg.Content["ename"].(string)
		value, _ := msg.Content["evalue"].(string)
		traceback, _ := msg.Content["traceback"].(string)

		execution.put(newOutputError(name, value, traceback))
		execution.setErrored()

	case MsgStream:
		name, ok := msg.Content["name"].(string)
		if !ok {
			return
		}

		if name == "stdout" {
			execution.put(newOutputStdout(msg.Header.Date, msg.Content["text"]))
			return
		}

		if name == "stderr" {
			execution.put(newOutputStderr(msg.Header.Date, msg.Content["text"]))
		}

	case MsgDisplayData, MsgExecuteResult:
		data, ok := msg.Content["data"].(map[string]any)
		if !ok {
			return
		}

		execution.put(newOutputResult(data))

	case MsgStatus:
		state, ok := msg.Content["execution_state"].(string)
		if !ok {
			return
		}

		if state == "busy" {
			execution.setInputAccepted()
		}

		if state == "idle" {
			if execution.inputAccepted {
				execution.put(newEndOfExecution())
			}
		}

		if state == "error" {
			name, _ := msg.Content["ename"].(string)
			value, _ := msg.Content["evalue"].(string)
			traceback, _ := msg.Content["traceback"].(string)

			execution.put(newOutputError(name, value, traceback))
			execution.put(newEndOfExecution())
		}

	case MsgExecuteReply:
		status, ok := msg.Content["status"].(string)
		if !ok {
			return
		}

		if status == "error" {
			if execution.errored {
				return
			}

			name, _ := msg.Content["ename"].(string)
			value, _ := msg.Content["evalue"].(string)
			traceback, _ := msg.Content["traceback"].(string)

			execution.put(newOutputError(name, value, traceback))
			execution.setErrored()
			return
		}

		if status == "abort" {
			execution.put(newOutputError("aborted", "execution was aborted", ""))
			return
		}

	case MsgExecuteInput:
		execution.setInputAccepted()
	}
}

type Exection struct {
	queue         chan any
	inputAccepted bool
	errored       bool
}

func (e *Exection) put(v any) { e.queue <- v }

func (e *Exection) setInputAccepted() { e.inputAccepted = true }

func (e *Exection) setErrored() { e.errored = true }

func (e *Exection) Recv() (any, error) {
	msg := <-e.queue
	switch val := msg.(type) {
	case EndOfExecution:
		return nil, io.EOF
	default:
		return val, nil
	}
}
