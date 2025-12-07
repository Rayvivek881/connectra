package constants

var (
	CompanyIndex      = "companies_index"
	ContactIndex      = "contacts_index"
	SearchTypeExact   = "exact"
	SearchTypeShuffle = "shuffle"

	DefaultPageSize = 25
)

var MatchAllQuery = map[string]any{
	"query": map[string]any{
		"match_all": map[string]any{},
	},
}
