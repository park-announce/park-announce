package entity

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	jwt.StandardClaims
}

type ResponseStatus struct {
	IsSucccess bool   `json:"issuccess"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Stack      string `json:"stack"`
}

type GoogleUser struct {
	Id           string `json:"sub"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsValidEmail bool   `json:"email_verified"`

	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Picture   string `json:"picture"`

	jwt.StandardClaims
}

type JWTToken struct {
	Header    string
	Payload   string
	Signature string
}

type GoogleJWTHeader struct {
	Alg string `json:"alg"`
	KID string `json:"kid"`
}

type GoogleOpenIDConfiguration struct {
	Issuer  string `json:"issuer"`
	JwksUri string `json:"jwks_uri"`
}

type GoogleOpenIDOAuthCertResponse struct {
	Keys []*GoogleOpenIDOAuthCertKey `json:"keys"`
}

type GoogleOpenIDOAuthCertKey struct {
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type Token struct {
	AccessToken string `json:"token"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SocketMessage struct {
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	Data          interface{} `json:"data"`
}

type ClientKafkaResponseMessage struct {
	ClientId string      `json:"client_id"`
	ApiId    string      `json:"api_id"`
	Data     interface{} `json:"data"`
}

type ClientKafkaRequestMessage struct {
	ClientId      string      `json:"client_id"`
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	ApiId         string      `json:"api_id"`
	Data          interface{} `json:"data"`
}

type ClientSocketResponseMessage struct {
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	Data          interface{} `json:"data"`
}
