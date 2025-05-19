package jupyter

const (
	OutTypeStdout        = "stdout"
	OutTypeStderr        = "stderr"
	OutTypeResult        = "result"
	OutTypeEndOfExection = "end_of_execution"
)

type OutputStdout struct {
	Type      string `json:"type,omitempty"`
	Timestamp any    `json:"timestamp,omitempty"`
	Data      any    `json:"data,omitempty"`
}

func newOutputStdout(ts, data any) OutputStdout {
	return OutputStdout{
		Type:      OutTypeStdout,
		Timestamp: ts,
		Data:      data,
	}
}

type OutputStderr struct {
	Type      string `json:"type,omitempty"`
	Timestamp any    `json:"timestamp,omitempty"`
	Data      any    `json:"data,omitempty"`
}

func newOutputStderr(ts, data any) OutputStderr {
	return OutputStderr{
		Type:      OutTypeStderr,
		Timestamp: ts,
		Data:      data,
	}
}

type OutputResult struct {
	Type string         `json:"type,omitempty"`
	Data map[string]any `json:"data,omitempty"`
}

func newOutputResult(data map[string]any) OutputResult {
	return OutputResult{
		Type: OutTypeResult,
		Data: data,
	}
}

type OutputError struct {
	Type      string `json:"type,omitempty"`
	Name      string `json:"name,omitempty"`
	Value     string `json:"value,omitempty"`
	Traceback string `json:"traceback,omitempty"`
}

func newOutputError(name, value, traceback string) OutputError {
	return OutputError{
		Name:      name,
		Value:     value,
		Traceback: traceback,
	}
}

type EndOfExecution struct {
	Type string `json:"type,omitempty"`
}

func newEndOfExecution() EndOfExecution {
	return EndOfExecution{OutTypeEndOfExection}
}
