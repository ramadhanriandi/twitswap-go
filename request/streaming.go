package request

type StartStreaming struct {
	Name string `json:"name" binding:"required"`
	Rule string `json:"rule" binding:"required"`
}

type StopStreaming struct {
	ID int64 `json:"id" binding:"required"`
}
