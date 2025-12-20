package utilities

// GeolocationData represents IP geolocation data from frontend
// This is a utility type for parsing geolocation data
// Conversion to models should be done in the service/helper layer
type GeolocationData struct {
	IP            *string  `json:"ip,omitempty"`
	Continent     *string  `json:"continent,omitempty"`
	ContinentCode *string  `json:"continent_code,omitempty"`
	Country       *string  `json:"country,omitempty"`
	CountryCode   *string  `json:"country_code,omitempty"`
	Region        *string  `json:"region,omitempty"`
	RegionName    *string  `json:"region_name,omitempty"`
	City          *string  `json:"city,omitempty"`
	District      *string  `json:"district,omitempty"`
	Zip           *string  `json:"zip,omitempty"`
	Lat           *float64 `json:"lat,omitempty"`
	Lon           *float64 `json:"lon,omitempty"`
	Timezone      *string  `json:"timezone,omitempty"`
	Offset        *int     `json:"offset,omitempty"`
	Currency      *string  `json:"currency,omitempty"`
	ISP           *string  `json:"isp,omitempty"`
	Org           *string  `json:"org,omitempty"`
	ASName        *string  `json:"asname,omitempty"`
	Reverse       *string  `json:"reverse,omitempty"`
	Device        *string  `json:"device,omitempty"`
	Proxy         *bool    `json:"proxy,omitempty"`
	Hosting       *bool    `json:"hosting,omitempty"`
}

