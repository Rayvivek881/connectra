package utilities

import (
	"reflect"
	"vivek-ray/constants"
)

func (q *NQLQuery) isEmpty() bool {
	queries := []ElasticQuery{
		q.Where.RangeQuery,
		q.Where.KeywordMatch,
	}
	for _, query := range queries {
		if len(query.Must) > 0 || len(query.MustNot) > 0 {
			return false
		}
	}
	return len(q.Where.TextMatch.Must) == 0 && len(q.Where.TextMatch.MustNot) == 0
}

func isSliceOrArray(value any) bool {
	rv := reflect.ValueOf(value)
	return (rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array) && rv.Len() > 0
}

func buildTextQueries(conditions []TextMatchStruct, isMust bool) []map[string]any {
	result := make([]map[string]any, 0)
	if len(conditions) == 0 {
		return result
	}

	queryMap := make(map[string][]map[string]any)
	for _, condition := range conditions {
		if _, ok := queryMap[condition.FilterKey]; !ok {
			queryMap[condition.FilterKey] = make([]map[string]any, 0)
		}
		switch condition.SearchType {
		case constants.SearchTypeExact:
			queryMap[condition.FilterKey] = append(queryMap[condition.FilterKey], map[string]any{
				"match_phrase": map[string]any{
					condition.FilterKey: map[string]any{
						"query": condition.TextValue,
						"slop":  condition.Slop,
					},
				},
			})
		case constants.SearchTypeShuffle:
			queryMap[condition.FilterKey] = append(queryMap[condition.FilterKey], map[string]any{
				"match": map[string]any{
					condition.FilterKey: map[string]any{
						"query":     condition.TextValue,
						"operator":  InlineIf(condition.Operator != "", condition.Operator, "and"),
						"fuzziness": InlineIf(condition.Fuzzy, "AUTO", 0),
					},
				},
			})
		case constants.SearchTypeSubstring:
			queryMap[condition.FilterKey] = append(queryMap[condition.FilterKey], map[string]any{
				"match": map[string]any{
					condition.FilterKey + ".ngram": map[string]any{
						"query":    condition.TextValue,
						"operator": InlineIf(condition.Operator != "", condition.Operator, "and"),
					},
				},
			})
		}
	}
	for _, queries := range queryMap {
		if len(queries) == 0 {
			continue
		}
		if isMust {
			result = append(
				result,
				map[string]any{
					"bool": map[string]any{
						"should":               queries,
						"minimum_should_match": 1,
					},
				},
			)
		} else {
			result = append(result, queries...)
		}
	}
	return result
}

func buildRangeQueries(conditions map[string]any) []map[string]any {
	queries := make([]map[string]any, 0)
	if len(conditions) == 0 {
		return queries
	}
	for key, value := range conditions {
		queries = append(queries, map[string]any{
			"range": map[string]any{
				key: value,
			},
		})
	}
	return queries
}

func buildKeywordQueries(conditions map[string]any) []map[string]any {
	queries := make([]map[string]any, 0)
	if len(conditions) == 0 {
		return queries
	}
	for key, value := range conditions {
		if isSliceOrArray(value) {
			queries = append(queries, map[string]any{
				"terms": map[string]any{
					key: value,
				},
			})
		} else {
			queries = append(queries, map[string]any{
				"term": map[string]any{
					key: value,
				},
			})
		}
	}
	return queries
}

func (q *NQLQuery) buildBoolQuery() map[string]any {
	mustQuery := buildTextQueries(q.Where.TextMatch.Must, true)
	mustNotQuery := buildTextQueries(q.Where.TextMatch.MustNot, false)

	mustNotKeywordQueries := buildKeywordQueries(q.Where.KeywordMatch.MustNot)
	if len(mustNotKeywordQueries) > 0 {
		mustNotQuery = append(mustNotQuery, mustNotKeywordQueries...)
	}

	filterQuery := buildRangeQueries(q.Where.RangeQuery.Must)
	keywordQueries := buildKeywordQueries(q.Where.KeywordMatch.Must)
	if len(keywordQueries) > 0 {
		filterQuery = append(filterQuery, keywordQueries...)
	}

	boolQuery := make(map[string]any)
	if len(mustQuery) > 0 {
		boolQuery["must"] = mustQuery
	}
	if len(mustNotQuery) > 0 {
		boolQuery["must_not"] = mustNotQuery
	}
	if len(filterQuery) > 0 {
		boolQuery["filter"] = filterQuery
	}
	return boolQuery
}

func (q *NQLQuery) addPagination(resultQuery map[string]any) {
	if len(q.SearchAfter) > 0 {
		resultQuery["search_after"] = q.SearchAfter
	}
	if q.Page > 0 {
		resultQuery["from"] = (q.Page - 1) * q.Limit
	}
	resultQuery["size"] = InlineIf(q.Limit > 0, q.Limit, constants.DefaultPageSize)
}

func (q *NQLQuery) addSort(resultQuery map[string]any) {
	if len(q.OrderBy) > 0 {
		sort := make([]map[string]any, 0, len(q.OrderBy))
		for _, order := range q.OrderBy {
			if order.OrderBy != "" {
				direction := "asc"
				if order.OrderDirection == "desc" {
					direction = "desc"
				}
				sort = append(sort, map[string]any{
					order.OrderBy: map[string]any{
						"order": direction,
					},
				})
			}
		}
		if len(sort) > 0 {
			resultQuery["sort"] = sort
		}
	}
}

func (q *NQLQuery) ToElasticsearchQuery(forCount bool) map[string]any {
	resultQuery := map[string]any{"_source": []string{"id"}}
	if !forCount {
		q.addPagination(resultQuery)
		q.addSort(resultQuery)
	}

	boolQuery := q.buildBoolQuery()
	if q.isEmpty() || len(boolQuery) == 0 {
		resultQuery["query"] = map[string]any{
			"match_all": map[string]any{},
		}
		return resultQuery
	}
	resultQuery["query"] = map[string]any{"bool": boolQuery}
	return resultQuery
}
