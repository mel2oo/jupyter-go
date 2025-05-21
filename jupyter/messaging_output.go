package jupyter

const (
	OutTypeError         = "error"
	OutTypeStdout        = "stdout"
	OutTypeStderr        = "stderr"
	OutTypeResult        = "result"
	OutTypeEndOfExection = "end_of_execution"
)

type Output struct {
	Type      string `json:"type,omitempty"`
	Timestamp any    `json:"timestamp,omitempty"`
	Data      any    `json:"data,omitempty"`
	Name      string `json:"name,omitempty"`
	Value     string `json:"value,omitempty"`
	Traceback string `json:"traceback,omitempty"`
}

func newOutputStdout(ts, data any) *Output {
	return &Output{
		Type:      OutTypeStdout,
		Timestamp: ts,
		Data:      data,
	}
}

func newOutputStderr(ts, data any) *Output {
	return &Output{
		Type:      OutTypeStderr,
		Timestamp: ts,
		Data:      data,
	}
}

func newOutputResult(data map[string]any) *Output {
	return &Output{
		Type: OutTypeResult,
		Data: data,
	}
}

func newOutputError(name, value, traceback string) *Output {
	return &Output{
		Type:      OutTypeError,
		Name:      name,
		Value:     value,
		Traceback: traceback,
	}
}

func newEndOfExecution() *Output {
	return &Output{
		Type: OutTypeEndOfExection,
	}
}
