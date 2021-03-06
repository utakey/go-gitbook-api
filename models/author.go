package models

// Author data as returned by the API
// Author only contains public information
type Author struct {
	Type     string     `json:"type"`
	Name     string     `json:"name"`
	Username string     `json:"username"`
	Urls     AuthorUrls `json:"urls"`
}
type AuthorUrls struct {
	Profile string `json:"profile"`
	Avatar  string `json:"avatar"`
}
