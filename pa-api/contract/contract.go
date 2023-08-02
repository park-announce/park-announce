package contract

type GetGoogleOAuthTokenRequest struct {
	Token      string `json:"token"`
	ClientType string `json:"client_type"`
}

type GetGoogleOAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IdToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type CorporationLocationUpdateRequest struct {
	Count         int32  `json:"count"`
	CorporationId string `json:"corporation_id"`
}

type CorporationOAuthTokenRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type CorporationUserInsertRequest struct {
	Email         string `json:"email"`
	CorporationId string `json:"corporation_id"`
}
