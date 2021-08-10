package response

type GetVisualizationByRuleID struct {
	TweetAnnotations []TweetAnnotation `json:"tweet_annotations"`
	TweetDomains     []TweetDomain     `json:"tweet_domains"`
}

type TweetAnnotation struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TweetDomain struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}
