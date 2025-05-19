package jupyter

import (
	"context"
	"fmt"
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

	res, err := ch.CodeExecute(context.TODO(), `# coding=utf-8

import types
import base64
import runtime


def handler():
	return {"key": "function hello1"}

handler()`)
	assert.NoError(t, err)

	for _, r := range res {
		switch v := r.(type) {
		case OutputError:
			fmt.Println(v.Traceback)
		default:
			fmt.Println(v)
		}
	}

	list, err := cli.Sessions.List(context.TODO())
	assert.NoError(t, err)

	for _, v := range list {
		assert.NoError(t, cli.Sessions.Delete(context.TODO(), v.ID))
	}
}
