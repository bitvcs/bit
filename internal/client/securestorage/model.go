package securestorage

type AccessToken struct {
	Domain       string `json:"domain"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}
