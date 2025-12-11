package utilities

type ElasticCount struct {
	Count int64 `json:"count"`
}

type FilterOrder struct {
	OrderBy        string `json:"order_by"`
	OrderDirection string `json:"order_direction"`
}

type TextMatchStruct struct {
	TextValue  string `json:"text_value"`
	FilterKey  string `json:"filter_key"`
	SearchType string `json:"search_type"`
	Slop       int    `json:"slop,omitempty"`
	Operator   string `json:"operator,omitempty"`
	Fuzzy      bool   `json:"fuzzy,omitempty"`
}

type ElasticQuery struct {
	Must    map[string]any `json:"must,omitempty"`
	MustNot map[string]any `json:"must_not,omitempty"`
}

type TextMatchQuery struct {
	Must    []TextMatchStruct `json:"must,omitempty"`
	MustNot []TextMatchStruct `json:"must_not,omitempty"`
}

type WhereStruct struct {
	TextMatch    TextMatchQuery `json:"text_matches"`
	KeywordMatch ElasticQuery   `json:"keyword_match"`
	RangeQuery   ElasticQuery   `json:"range_query"`
}

type NQLQuery struct {
	Where       WhereStruct   `json:"where"`
	OrderBy     []FilterOrder `json:"order_by,omitempty"`
	SearchAfter []string      `json:"search_after,omitempty"`

	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type InsertFileJobData struct {
	FileS3Key    string `json:"file_s3_key"`
	FileS3Bucket string `json:"file_s3_bucket"`
	FileType     string `json:"file_type"`
}
