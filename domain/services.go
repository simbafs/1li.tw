package domain

// UAParserResult holds the structured data from a User-Agent string.
type UAParserResult struct {
	OSName      string
	BrowserName string
}

// UAParserService defines the contract for a service that can parse a User-Agent string.
type UAParserService interface {
	Parse(userAgent string) *UAParserResult
}

// GeoIPService defines the contract for a service that can look up the country from an IP address.
type GeoIPService interface {
	CountryCode(ipAddress string) (string, error)
}
