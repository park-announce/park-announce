package entity

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Id          string `db:"id" json:"id"`
	FirstName   string `db:"first_name" json:"first_name"`
	LastName    string `db:"last_name" json:"last_name"`
	Email       string `db:"email" json:"email"`
	Picture     string `db:"picture" json:"picture"`
	Status      int16  `db:"status" json:"status"`
	MobilePhone string `db:"mobile_phone" json:"mobile_phone"`
	CityCode    int16  `db:"city_code" json:"city_code"`
	jwt.StandardClaims
}

type CorporationUser struct {
	Id            string `db:"id" json:"id"`
	Password      string `db:"password" json:"password"`
	Email         string `db:"email" json:"email"`
	CorporationId string `db:"corpoaration_id" json:"corpoaration_id"`
	RoleId        string `db:"role_id" json:"role_id"`
}

type CorporationUserRole struct {
	Id   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
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
	Timeout       int64       `json:"timeout"`
	Data          interface{} `json:"data"`
}

type ClientKafkaResponseMessage struct {
	ClientId string      `json:"client_id"`
	ApiId    string      `json:"api_id"`
	Data     interface{} `json:"data"`
}

type ClientKafkaRequestMessage struct {
	ClientId        string      `json:"client_id"`
	Operation       string      `json:"operation"`
	TransactionId   string      `json:"transaction_id"`
	ApiId           string      `json:"api_id"`
	TransactionTime int64       `json:"transaction_time"`
	Timeout         int64       `json:"timeout"`
	Data            interface{} `json:"data"`
}

type ClientSocketResponseMessage struct {
	Operation     string      `json:"operation"`
	TransactionId string      `json:"transaction_id"`
	Data          interface{} `json:"data"`
}

type IEntity interface {
	Do()
}

func (user User) Do()                               {}
func (corporationUser CorporationUser) Do()         {}
func (corporationUserRole CorporationUserRole) Do() {}
