package jupyter

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	cli, err := NewClient(WithBaseURL("http://127.0.0.1:8888"))
	assert.NoError(t, err)

	_, err = cli.GetVersion(context.TODO())
	assert.NoError(t, err)
}

func TestSession(t *testing.T) {
	cli, err := NewClient(WithBaseURL("http://127.0.0.1:8888"))
	assert.NoError(t, err)

	s1, err := cli.Sessions.Create(context.TODO(), &Session{
		Name:   "test1",
		Type:   SessionNotebook,
		Path:   "/home/workspace/s1",
		Kernel: Kernel{Name: "python3"},
	})
	assert.NoError(t, err)

	s2, err := cli.Sessions.Create(context.TODO(), &Session{
		Name:   "test2",
		Type:   SessionNotebook,
		Path:   "/home/workspace/s2",
		Kernel: Kernel{Name: "python3"},
	})
	assert.NoError(t, err)

	_, err = cli.Sessions.Get(context.TODO(), s1.ID)
	assert.NoError(t, err)

	_, err = cli.Sessions.Get(context.TODO(), s2.ID)
	assert.NoError(t, err)

	_, err = cli.Sessions.Update(context.TODO(), s1.ID, &Session{Name: "rename"})
	assert.NoError(t, err)

	list, err := cli.Sessions.List(context.TODO())
	assert.NoError(t, err)

	for _, v := range list {
		assert.NoError(t, cli.Sessions.Delete(context.TODO(), v.ID))
	}
}

func TestSessionBenchmark(t *testing.T) {
	cli, err := NewClient(WithBaseURL("http://127.0.0.1:8888"))
	assert.NoError(t, err)

	start := time.Now()

	for i := range 10 {
		_, err := cli.Sessions.Create(context.TODO(), &Session{
			Name:   "test",
			Type:   SessionNotebook,
			Path:   fmt.Sprintf("/home/workspace/s%d", i),
			Kernel: Kernel{Name: "python3"},
		})
		assert.NoError(t, err)

		fmt.Println(time.Since(start))
	}

	list, err := cli.Sessions.List(context.TODO())
	assert.NoError(t, err)

	for _, v := range list {
		assert.NoError(t, cli.Sessions.Delete(context.TODO(), v.ID))
	}
}
