package response

import "time"

type GetVisualizationByRuleID struct {
	TweetAnnotations  []TweetAnnotation  `json:"tweet_annotations"`
	TweetDomains      []TweetDomain      `json:"tweet_domains"`
	TweetGeolocations []TweetGeolocation `json:"tweet_geolocations"`
	TweetHashtags     []TweetHashtag     `json:"tweet_hashtags"`
	TweetLanguages    TweetLanguage      `json:"tweet_languages"`
	TweetMetrics      TweetMetric        `json:"tweet_metrics"`
	TweetPopularities []TweetPopularity  `json:"tweet_popularities"`
	TweetTypes        TweetType          `json:"tweet_types"`
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

type TweetMetric struct {
	Cumulative TweetMetricCumulative `json:"cumulative"`
	Interval   []TweetMetricInterval `json:"interval"`
}

type TweetMetricInterval struct {
	Like      int64     `json:"like"`
	Reply     int64     `json:"reply"`
	Retweet   int64     `json:"retweet"`
	Quote     int64     `json:"quote"`
	CreatedAt time.Time `json:"created_at"`
}

type TweetMetricCumulative struct {
	Like    int64 `json:"like"`
	Reply   int64 `json:"reply"`
	Retweet int64 `json:"retweet"`
	Quote   int64 `json:"quote"`
}

type TweetPopularity struct {
	TweetID string `json:"tweet_id"`
	Count   int64  `json:"count"`
}

type TweetType struct {
	Cumulative TweetTypeCumulative `json:"cumulative"`
	Interval   []TweetTypeInterval `json:"interval"`
}

type TweetTypeInterval struct {
	Tweet     int64     `json:"tweet"`
	Retweet   int64     `json:"retweet"`
	Quote     int64     `json:"quote"`
	Reply     int64     `json:"reply"`
	CreatedAt time.Time `json:"created_at"`
}

type TweetTypeCumulative struct {
	Tweet   int64 `json:"tweet"`
	Retweet int64 `json:"retweet"`
	Quote   int64 `json:"quote"`
	Reply   int64 `json:"reply"`
}
