package contract

type GetGoogleOAuthTokenRequest struct {
	Token string `json:"token"`
}

type GetGoogleOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
