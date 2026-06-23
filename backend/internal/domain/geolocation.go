package domain

// GeoLocation holds IP-derived geographic metadata.
type GeoLocation struct {
	IP          string
	CountryCode string
	CountryName string
	Region      string
	City        string
	Timezone    string
	Currency    string
	IsEU        bool
}
