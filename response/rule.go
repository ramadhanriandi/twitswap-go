package response

type GetVisualizationByRuleID struct {
	TweetAnnotations  []TweetAnnotation  `json:"tweet_annotations"`
	TweetDomains      []TweetDomain      `json:"tweet_domains"`
	TweetGeolocations []TweetGeolocation `json:"tweet_geolocations"`
	TweetHashtags     []TweetHashtag     `json:"tweet_hashtags"`
	TweetLanguages    TweetLanguage      `json:"tweet_languages"`
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

type TweetHashtag struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

type TweetLanguage struct {
	En    int64 `json:"en"`
	In    int64 `json:"in"`
	Other int64 `json:"other"`
}
