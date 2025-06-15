package jupyter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Channel struct {
	ws         *websocket.Conn
	sessionID  string
	executions sync.Map
	done       chan struct{}
	errChan    chan error
}

func (s *SessionService) Connect(ctx context.Context, sessionID, kernelID string) (*Channel, error) {
	// 设置连接超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	ws, _, err := websocket.DefaultDialer.DialContext(ctx,
		fmt.Sprintf("ws://%s/api/kernels/%s/channels?session_id=%s",
			s.client.baseURL.Host, kernelID, sessionID), nil)
	if err != nil {
		return nil, err
	}

	c := &Channel{
		ws:         ws,
		sessionID:  sessionID,
		executions: sync.Map{},
		done:       make(chan struct{}),
		errChan:    make(chan error, 1),
	}

	return c, c.receiveMessage()
}

func (c *Channel) Close() error {
	close(c.done)
	return c.ws.Close()
}

// code execute
func (c *Channel) CodeExecute(ctx context.Context, code string) ([]*Output, error) {
	msgID := uuid.NewString()
	execution := &Exection{queue: make(chan *Output, 10)}

	c.executions.Store(msgID, execution)
	defer c.executions.Delete(msgID)

	if err := c.ws.WriteJSON(c.newExecuteRequest(msgID, code)); err != nil {
		return nil, err
	}

	res := make([]*Output, 0)
	for {
		select {
		case <-ctx.Done():
			return res, ctx.Err()
		case err := <-c.errChan:
			return res, err
		case val, ok := <-execution.queue:
			if !ok {
				return res, nil
			}
			if val.Type == OutTypeEndOfExection {
				return res, nil
			}
			res = append(res, val)
		}
	}
}

func (c *Channel) CodeExecuteStream(ctx context.Context, code string) (*Exection, error) {
	msgID := uuid.NewString()
	execution := &Exection{queue: make(chan *Output, 10)}

	c.executions.Store(msgID, execution)
	defer c.executions.Delete(msgID)

	if err := c.ws.WriteJSON(c.newExecuteRequest(msgID, code)); err != nil {
		return nil, err
	}

	return execution, nil
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
			"store_history":    true,
			"user_expressions": map[string]any{},
			"allow_stdin":      false,
			"stop_on_error":    true,
		},
		Channel: "shell",
	}
}

func (c *Channel) receiveMessage() error {
	go func() {
		defer close(c.errChan)

		for {
			select {
			case <-c.done:
				return
			default:
				_, raw, err := c.ws.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						c.errChan <- fmt.Errorf("websocket read error: %w", err)
					}
					return
				}

				var msg JupyterResponseMessage
				if err := json.Unmarshal(raw, &msg); err != nil {
					c.errChan <- fmt.Errorf("unmarshal message error: %w", err)
					return
				}

				c.processMessage(&msg)
			}
		}
	}()

	return nil
}

func (c *Channel) processMessage(msg *JupyterResponseMessage) {
	v, ok := c.executions.Load(msg.ParentHeader.MsgID)
	if !ok {
		return
	}

	execution, ok := v.(*Exection)
	if !ok {
		return
	}

	switch msg.MsgType {
	case MsgError:
		if execution.errored {
			return
		}

		execution.put(newOutputError(msg.Content.Ename, msg.Content.Evalue,
			strings.Join(msg.Content.Traceback, "\n")))
		execution.setErrored()

	case MsgStream:
		if msg.Content.Name == "stdout" {
			execution.put(newOutputStdout(msg.Header.Date, msg.Content.Text))
			return
		}

		if msg.Content.Name == "stderr" {
			execution.put(newOutputStderr(msg.Header.Date, msg.Content.Text))
		}

	case MsgDisplayData, MsgExecuteResult:
		execution.put(newOutputResult(msg.Content.Data))

	case MsgStatus:
		if msg.Content.ExecutionState == "busy" {
			execution.setInputAccepted()
		}

		if msg.Content.ExecutionState == "idle" {
			if execution.inputAccepted {
				execution.put(newEndOfExecution())
			}
		}

		if msg.Content.ExecutionState == "error" {
			execution.put(newOutputError(msg.Content.Ename, msg.Content.Evalue,
				strings.Join(msg.Content.Traceback, "\n")))
			execution.put(newEndOfExecution())
		}

	case MsgExecuteReply:
		if msg.Content.Status == "error" {
			if execution.errored {
				return
			}

			execution.put(newOutputError(msg.Content.Ename, msg.Content.Evalue,
				strings.Join(msg.Content.Traceback, "\n")))
			execution.setErrored()
			return
		}

		if msg.Content.Status == "abort" {
			execution.put(newOutputError("aborted", "execution was aborted", ""))
			return
		}

	case MsgExecuteInput:
		execution.setInputAccepted()
	}
}

type Exection struct {
	queue         chan *Output
	inputAccepted bool
	errored       bool
}

func (e *Exection) put(v *Output) { e.queue <- v }

func (e *Exection) setInputAccepted() { e.inputAccepted = true }

func (e *Exection) setErrored() { e.errored = true }

func (e *Exection) Recv() (*Output, error) {
	msg := <-e.queue
	switch msg.Type {
	case OutTypeEndOfExection:
		return nil, io.EOF
	default:
		return msg, nil
	}
}
