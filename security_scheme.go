package spec

import (
	"encoding/json"

	"github.com/casualjim/go-swagger/jsonpointer"
	"github.com/casualjim/go-swagger/swag"
)

const (
	basic       = "basic"
	apiKey      = "apiKey"
	oauth2      = "oauth2"
	implicit    = "implicit"
	password    = "password"
	application = "application"
	accessCode  = "accessCode"
)

// BasicAuth creates a basic auth security scheme
func BasicAuth() *SecurityScheme {
	return &SecurityScheme{securitySchemeProps: securitySchemeProps{Type: basic}}
}

// APIKeyAuth creates an api key auth security scheme
func APIKeyAuth(fieldName, valueSource string) *SecurityScheme {
	return &SecurityScheme{securitySchemeProps: securitySchemeProps{Type: apiKey, Name: fieldName, In: valueSource}}
}

// OAuth2Implicit creates an implicit flow oauth2 security scheme
func OAuth2Implicit(authorizationURL string) *SecurityScheme {
	return &SecurityScheme{securitySchemeProps: securitySchemeProps{
		Type:             oauth2,
		Flow:             implicit,
		AuthorizationURL: authorizationURL,
	}}
}

// OAuth2Password creates a password flow oauth2 security scheme
func OAuth2Password(tokenURL string) *SecurityScheme {
	return &SecurityScheme{securitySchemeProps: securitySchemeProps{
		Type:     oauth2,
		Flow:     password,
		TokenURL: tokenURL,
	}}
}

// OAuth2Application creates an application flow oauth2 security scheme
func OAuth2Application(tokenURL string) *SecurityScheme {
	return &SecurityScheme{securitySchemeProps: securitySchemeProps{
		Type:     oauth2,
		Flow:     application,
		TokenURL: tokenURL,
	}}
}

// OAuth2AccessToken creates an access token flow oauth2 security scheme
func OAuth2AccessToken(authorizationURL, tokenURL string) *SecurityScheme {
	return &SecurityScheme{securitySchemeProps: securitySchemeProps{
		Type:             oauth2,
		Flow:             accessCode,
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
	}}
}

type securitySchemeProps struct {
	Description      string            `json:"description,omitempty"`
	Type             string            `json:"type"`
	Name             string            `json:"name,omitempty"`             // api key
	In               string            `json:"in,omitempty"`               // api key
	Flow             string            `json:"flow,omitempty"`             // oauth2
	AuthorizationURL string            `json:"authorizationUrl,omitempty"` // oauth2
	TokenURL         string            `json:"tokenUrl,omitempty"`         // oauth2
	Scopes           map[string]string `json:"scopes,omitempty"`           // oauth2
}

// AddScope adds a scope to this security scheme
func (s *securitySchemeProps) AddScope(scope, description string) {
	if s.Scopes == nil {
		s.Scopes = make(map[string]string)
	}
	s.Scopes[scope] = description
}

// SecurityScheme allows the definition of a security scheme that can be used by the operations.
// Supported schemes are basic authentication, an API key (either as a header or as a query parameter)
// and OAuth2's common flows (implicit, password, application and access code).
//
// For more information: http://goo.gl/8us55a#securitySchemeObject
type SecurityScheme struct {
	vendorExtensible
	securitySchemeProps
}

// JSONLookup implements an interface to customize json pointer lookup
func (s SecurityScheme) JSONLookup(token string) (interface{}, error) {
	if ex, ok := s.Extensions[token]; ok {
		return &ex, nil
	}

	r, _, err := jsonpointer.GetForToken(s.securitySchemeProps, token)
	return r, err
}

// MarshalJSON marshal this to JSON
func (s SecurityScheme) MarshalJSON() ([]byte, error) {
	b1, err := json.Marshal(s.securitySchemeProps)
	if err != nil {
		return nil, err
	}
	b2, err := json.Marshal(s.vendorExtensible)
	if err != nil {
		return nil, err
	}
	return swag.ConcatJSON(b1, b2), nil
}

// UnmarshalJSON marshal this from JSON
func (s *SecurityScheme) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.securitySchemeProps); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &s.vendorExtensible); err != nil {
		return err
	}
	return nil
}
