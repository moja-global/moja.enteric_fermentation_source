package cookies

type Cred struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURIs []string `json:"redirect_uris"`
	AuthURI      string   `json:"auth_uri"`
	TokenURI     string   `json:"token_uri"`
}
type SecretFile struct {
	Web       *Cred `json:"web"`
	Installed *Cred `json:"installed"`
}
