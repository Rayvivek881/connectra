package models

// Role represents user role types
type Role string

const (
	RoleSuperAdmin Role = "SuperAdmin"
	RoleAdmin      Role = "Admin"
	RoleProUser    Role = "ProUser"
	RoleFreeUser   Role = "FreeUser"
)

// UserHistoryEventType represents event types for user history
type UserHistoryEventType string

const (
	EventTypeRegistration UserHistoryEventType = "registration"
	EventTypeLogin        UserHistoryEventType = "login"
)

// ActivityServiceType represents service types for user activities
type ActivityServiceType string

const (
	ServiceTypeLinkedIn ActivityServiceType = "linkedin"
	ServiceTypeEmail    ActivityServiceType = "email"
)

// ActivityActionType represents action types for user activities
type ActivityActionType string

const (
	ActionTypeSearch ActivityActionType = "search"
	ActionTypeExport ActivityActionType = "export"
)

// ActivityStatus represents status types for user activities
type ActivityStatus string

const (
	ActivityStatusSuccess ActivityStatus = "success"
	ActivityStatusFailed  ActivityStatus = "failed"
	ActivityStatusPartial ActivityStatus = "partial"
)

// FeatureType represents feature types for usage tracking
type FeatureType string

const (
	FeatureTypeAIChat           FeatureType = "AI_CHAT"
	FeatureTypeBulkExport       FeatureType = "BULK_EXPORT"
	FeatureTypeAPIKeys          FeatureType = "API_KEYS"
	FeatureTypeTeamManagement   FeatureType = "TEAM_MANAGEMENT"
	FeatureTypeEmailFinder      FeatureType = "EMAIL_FINDER"
	FeatureTypeVerifier         FeatureType = "VERIFIER"
	FeatureTypeLinkedIn         FeatureType = "LINKEDIN"
	FeatureTypeDataSearch       FeatureType = "DATA_SEARCH"
	FeatureTypeAdvancedFilters  FeatureType = "ADVANCED_FILTERS"
	FeatureTypeAISummaries      FeatureType = "AI_SUMMARIES"
	FeatureTypeSaveSearches     FeatureType = "SAVE_SEARCHES"
	FeatureTypeBulkVerification FeatureType = "BULK_VERIFICATION"
)
