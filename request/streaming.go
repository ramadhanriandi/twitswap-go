package request

type StartStreaming struct {
	Name string `json:"name" binding:"required"`
	Rule string `json:"rule" binding:"required"`
}
