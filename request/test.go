package request

type StartTesting struct {
	Count int64 `json:"count" binding:"required"`
}
