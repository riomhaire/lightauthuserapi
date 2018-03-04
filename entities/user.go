package entities

type User struct {
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Enabled  bool     `json:"enabled,omitempty"`
	Roles    []string `json:"roles,omitempty"`
	Claim1   string   `json:"claim1,omitempty"`
	Claim2   string   `json:"claim2,omitempty"`
}
