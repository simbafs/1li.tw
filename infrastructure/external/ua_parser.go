package external

import (
	"1litw/domain"

	"github.com/ua-parser/uap-go/uaparser"
)

var _ domain.UAParserService = (*uaParserService)(nil)

// uaParserService implements the domain.UAParserService interface using uap-go.
type uaParserService struct {
	parser *uaparser.Parser
}

// NewUAParserService creates a new User-Agent parsing service.
// It panics if the regexes.yaml file cannot be found or parsed,
// as this is a critical configuration for the service to function.
func NewUAParserService() domain.UAParserService {
	parser := uaparser.NewFromSaved()
	return &uaParserService{parser: parser}
}

// Parse extracts OS and browser information from a User-Agent string.
func (s *uaParserService) Parse(userAgent string) *domain.UAParserResult {
	client := s.parser.Parse(userAgent)
	return &domain.UAParserResult{
		OSName:      client.Os.Family,
		BrowserName: client.UserAgent.Family,
	}
}
