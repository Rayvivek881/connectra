package utilities

import (
	"testing"
	"vivek-ray/constants"
)

func TestVQLQuery_ToElasticsearchQuery(t *testing.T) {
	tests := []struct {
		name         string
		query        VQLQuery
		forCount     bool
		sourceFields []string
		wantQuery    bool // Just check if query is generated
	}{
		{
			name: "empty query",
			query: VQLQuery{
				Where: WhereStruct{},
			},
			forCount:     false,
			sourceFields: []string{},
			wantQuery:    true,
		},
		{
			name: "text match query",
			query: VQLQuery{
				Where: WhereStruct{
					TextMatch: TextMatchQuery{
						Must: []TextMatchStruct{
							{
								TextValue:  "test",
								FilterKey:  "name",
								SearchType: constants.SearchTypeShuffle,
								Fuzzy:      true,
							},
						},
					},
				},
			},
			forCount:     false,
			sourceFields: []string{},
			wantQuery:    true,
		},
		{
			name: "keyword match query",
			query: VQLQuery{
				Where: WhereStruct{
					KeywordMatch: ElasticQuery{
						Must: map[string]any{
							"country": []string{"USA", "Canada"},
						},
					},
				},
			},
			forCount:     false,
			sourceFields: []string{},
			wantQuery:    true,
		},
		{
			name: "range query",
			query: VQLQuery{
				Where: WhereStruct{
					RangeQuery: ElasticQuery{
						Must: map[string]any{
							"employees_count": map[string]any{
								"gte": 50,
								"lte": 1000,
							},
						},
					},
				},
			},
			forCount:     false,
			sourceFields: []string{},
			wantQuery:    true,
		},
		{
			name: "combined query",
			query: VQLQuery{
				Where: WhereStruct{
					TextMatch: TextMatchQuery{
						Must: []TextMatchStruct{
							{
								TextValue:  "software",
								FilterKey:  "name",
								SearchType: constants.SearchTypeShuffle,
							},
						},
					},
					KeywordMatch: ElasticQuery{
						Must: map[string]any{
							"country": []string{"USA"},
						},
					},
					RangeQuery: ElasticQuery{
						Must: map[string]any{
							"employees_count": map[string]any{
								"gte": 100,
							},
						},
					},
				},
				OrderBy: []FilterOrder{
					{
						OrderBy:        "employees_count",
						OrderDirection: "desc",
					},
				},
				Page:  1,
				Limit: 25,
			},
			forCount:     false,
			sourceFields: []string{"name", "employees_count"},
			wantQuery:    true,
		},
		{
			name: "count query",
			query: VQLQuery{
				Where: WhereStruct{
					KeywordMatch: ElasticQuery{
						Must: map[string]any{
							"country": []string{"USA"},
						},
					},
				},
			},
			forCount:     true,
			sourceFields: []string{},
			wantQuery:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.ToElasticsearchQuery(tt.forCount, tt.sourceFields)
			
			if !tt.wantQuery {
				if result != nil {
					t.Error("Expected nil query")
				}
				return
			}
			
			if result == nil {
				t.Error("Expected non-nil query")
				return
			}
			
			// Basic structure validation
			if _, ok := result["query"]; !ok {
				t.Error("Expected 'query' key in result")
			}
		})
	}
}

func TestVQLQuery_isEmpty(t *testing.T) {
	tests := []struct {
		name  string
		query VQLQuery
		want  bool
	}{
		{
			name: "empty query",
			query: VQLQuery{
				Where: WhereStruct{},
			},
			want: true,
		},
		{
			name: "query with text match",
			query: VQLQuery{
				Where: WhereStruct{
					TextMatch: TextMatchQuery{
						Must: []TextMatchStruct{
							{TextValue: "test", FilterKey: "name"},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "query with keyword match",
			query: VQLQuery{
				Where: WhereStruct{
					KeywordMatch: ElasticQuery{
						Must: map[string]any{
							"country": []string{"USA"},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "query with range query",
			query: VQLQuery{
				Where: WhereStruct{
					RangeQuery: ElasticQuery{
						Must: map[string]any{
							"employees_count": map[string]any{"gte": 50},
						},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.isEmpty()
			if result != tt.want {
				t.Errorf("isEmpty() = %v, want %v", result, tt.want)
			}
		})
	}
}
