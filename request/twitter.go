package request

type PostRules struct {
	Add []PostRulesValue `json:"add"`
}

type PostRulesValue struct {
	Value string `json:"value"`
}

type DeleteRules struct {
	Delete DeleteRulesBody `json:"delete"`
}

type DeleteRulesBody struct {
	Ids []string `json:"ids"`
}
