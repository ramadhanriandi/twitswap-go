package response

import (
	"database/sql"
	"time"
)

type StartStreaming struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	StartTime time.Time `json:"start_time"`
	RuleID    string    `json:"rule_id"`
	Rule      string    `json:"rule"`
}

type GetLatestStreaming struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	StartTime time.Time    `json:"start_time"`
	EndTime   sql.NullTime `json:"end_time"`
	RuleID    string       `json:"rule_id"`
	Rule      string       `json:"rule"`
}

type GetAllStreaming struct {
	ID        int64        `json:"id"`
	Name      string       `json:"name"`
	StartTime time.Time    `json:"start_time"`
	EndTime   sql.NullTime `json:"end_time"`
	RuleID    string       `json:"rule_id"`
	Rule      string       `json:"rule"`
}
