package domain

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Domain       string `json:"domain"`
	ExpiresIn    int    `json:"expires_in"`
}
