package response

type GetRules struct {
	Data []GetRulesData `json:"data"`
	Meta GetRulesMeta   `json:"meta"`
}

type GetRulesData struct {
	Id    string  `json:"id"`
	Value string  `json:"value"`
	Tag   *string `json:"tag"`
}

type GetRulesMeta struct {
	Sent string `json:"sent"`
}

type PostRules struct {
	Data []PostRulesData `json:"data"`
	Meta PostRulesMeta   `json:"meta"`
}

type PostRulesData struct {
	Value string `json:"value"`
	Id    string `json:"id"`
}

type PostRulesMeta struct {
	Sent    string           `json:"sent"`
	Summary PostRulesSummary `json:"summary"`
}

type PostRulesSummary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Valid      int `json:"valid"`
	Invalid    int `json:"invalid"`
}

type DeleteRules struct {
	Meta DeleteRulesMeta `json:"meta"`
}

type DeleteRulesMeta struct {
	Sent    string             `json:"sent"`
	Summary DeleteRulesSummary `json:"summary"`
}

type DeleteRulesSummary struct {
	Deleted    int `json:"deleted"`
	NotDeleted int `json:"not_deleted"`
}
