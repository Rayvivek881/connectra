package utilities

type OpenSearchCount struct {
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

type OpenSearchQuery struct {
	Must    map[string]any `json:"must,omitempty"`
	MustNot map[string]any `json:"must_not,omitempty"`
}

type TextMatchQuery struct {
	Must    []TextMatchStruct `json:"must,omitempty"`
	MustNot []TextMatchStruct `json:"must_not,omitempty"`
}

type WhereStruct struct {
	TextMatch    TextMatchQuery   `json:"text_matches"`
	KeywordMatch OpenSearchQuery `json:"keyword_match"`
	RangeQuery   OpenSearchQuery `json:"range_query"`
}

type CompanyConfig struct {
	Populate      bool     `json:"populate,omitempty"`
	SelectColumns []string `json:"select_columns,omitempty"`
}

type VQLQuery struct { // vivek Query Language
	Where WhereStruct `json:"where"`

	OrderBy       []FilterOrder  `json:"order_by,omitempty"`
	Cursor        []string       `json:"cursor,omitempty"`
	SelectColumns []string       `json:"select_columns,omitempty"`
	CompanyConfig *CompanyConfig `json:"company_config,omitempty"`

	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type InsertFileJobData struct {
	FileS3Key    string `json:"s3_key"`
	FileS3Bucket string `json:"s3_bucket"`
}

type ExportFileJobData struct {
	FileS3Bucket string   `json:"s3_bucket"`
	Service      string   `json:"service"`
	VQL          VQLQuery `json:"vql"`
}
