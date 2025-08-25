package applog

type LogEntry struct {
	Level     string `json:"level"`
	Msg       string `json:"msg"`
	Time      string `json:"time"`
	Path      string `json:"path,omitempty"`
	Log_scope string `json:"log_scope"`
	TraceID   string `json:"trace_id,omitempty"`
	Recover   any    `json:"recover,omitempty"`
	Stack     string `json:"stack,omitempty"`
	IPAddress string `json:"ip,omitempty"`
	Duration  string `json:"duration,omitempty"`
}
