package jupyter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	cli, err := NewClient(WithBaseURL("http://192.168.134.142:8888"))
	assert.NoError(t, err)

	res, err := cli.GetVersion(context.Background())
	assert.NoError(t, err)

	println(res)
}
