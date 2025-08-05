package exchangeCredentials

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
type RenewAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
