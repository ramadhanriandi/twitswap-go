package request

type StartStreaming struct {
	Rule string `json:"rule" binding:"required"`
}

type StopStreaming struct {
	RuleId string `json:"rule_id" binding:"required"`
}
