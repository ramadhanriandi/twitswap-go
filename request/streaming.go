package request

type StartStreaming struct {
	Rule string `json:"rule" binding:"required"`
}
