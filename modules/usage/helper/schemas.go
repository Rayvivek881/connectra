package helper

// FeatureUsageItem represents usage for a single feature
type FeatureUsageItem struct {
	Used  int `json:"used"`
	Limit int `json:"limit"`
}

// CurrentUsageResponse maps feature names to usage data
// Matches frontend expectation: Record<Feature, { used: number; limit: number }>
type CurrentUsageResponse map[string]FeatureUsageItem

// TrackUsageRequest for tracking feature usage
type TrackUsageRequest struct {
	Feature string `json:"feature" binding:"required"`
	Amount  int    `json:"amount,omitempty" binding:"omitempty,gte=1"` // Default: 1
}

// TrackUsageResponse for track usage response
type TrackUsageResponse struct {
	Feature string `json:"feature"`
	Used    int    `json:"used"`
	Limit   int    `json:"limit"`
	Success bool   `json:"success"`
}
