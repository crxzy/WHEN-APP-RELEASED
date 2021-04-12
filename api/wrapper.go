package api

// Wrapper response wrapper
type Wrapper struct {
	Data   interface{} `json:"data,omitempty"`
	Msg    string      `json:"msg,omitempty"`
	Status int         `json:"status"`
}
