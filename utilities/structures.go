package utilities

import "github.com/uptrace/bun"

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

type DefaultFilters struct {
	Page    int           `json:"page,omitempty"`
	Limit   int           `json:"limit,omitempty"`
	OrderBy []FilterOrder `json:"order_by,omitempty"`
}

func (f *DefaultFilters) ToWhere(query *bun.SelectQuery) *bun.SelectQuery {
	if f.Page > 0 {
		query = query.Offset((f.Page - 1) * f.Limit)
	}
	if f.Limit > 0 {
		query = query.Limit(f.Limit)
	}
	if len(f.OrderBy) > 0 {
		for _, order := range f.OrderBy {
			query = query.Order(order.OrderBy + " " + order.OrderDirection)
		}
	}
	return query
}

type VQLQuery struct { // vivek Query Language
	Where WhereStruct `json:"where"`

	Cursor        []string       `json:"cursor,omitempty"`
	SelectColumns []string       `json:"select_columns,omitempty"`
	CompanyConfig *CompanyConfig `json:"company_config,omitempty"`

	DefaultFilters
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
