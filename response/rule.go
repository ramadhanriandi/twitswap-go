package response

type GetVisualizationByRuleID struct {
	TweetAnnotations  []TweetAnnotation  `json:"tweet_annotations"`
	TweetDomains      []TweetDomain      `json:"tweet_domains"`
	TweetGeolocations []TweetGeolocation `json:"tweet_geolocations"`
}

type TweetAnnotation struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TweetDomain struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TweetGeolocation struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
