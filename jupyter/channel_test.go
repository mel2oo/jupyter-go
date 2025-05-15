package jupyter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodeExecute(t *testing.T) {
	cli, err := NewClient(WithBaseURL("http://127.0.0.1:8888"))
	assert.NoError(t, err)

	s1, err := cli.Sessions.Create(context.TODO(), &Session{
		Name:   "test1",
		Type:   SessionNotebook,
		Path:   "/home/workspace/s1",
		Kernel: Kernel{Name: "python3"},
	})
	assert.NoError(t, err)

	ch, err := cli.Sessions.Connect(context.TODO(), s1.Kernel.ID, s1.ID)
	assert.NoError(t, err)

	go ch.Recv()

	assert.NoError(t, ch.CodeExecute(context.TODO(), s1.ID, `print("print hello!")

def handler():
    return {"key": "function hello1"}

def handler1():
    return "function hello2"

handler()
handler()`))

	select {}
}
