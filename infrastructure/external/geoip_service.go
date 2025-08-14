package external

import "1litw/domain"

// geoIPService implements the domain.GeoIPService interface with mock data.
type geoIPService struct{}

// NewGeoIPService creates a new mock GeoIP service.
func NewGeoIPService() domain.GeoIPService {
	return &geoIPService{}
}

// CountryCode returns a mock country code for a given IP address.
// In a real implementation, this would query a GeoIP database or an external API.
func (s *geoIPService) CountryCode(ipAddress string) (string, error) {
	// Return a non-real code for local/private IPs.
	if ipAddress == "127.0.0.1" || ipAddress == "::1" {
		return "localhost", nil
	}
	// Return a dummy value for any other IP.
	return "XX", nil
}
