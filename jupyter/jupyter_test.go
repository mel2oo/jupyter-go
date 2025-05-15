package jupyter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetVersion(t *testing.T) {
	cli, err := NewClient(WithBaseURL("http://127.0.0.1:8888"))
	assert.NoError(t, err)

	_, err = cli.GetVersion(context.TODO())
	assert.NoError(t, err)
}
