package util

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"strings"

	"github.com/park-announce/pa-api/client"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/types"
	"github.com/sethvargo/go-password/password"

	b64 "encoding/base64"

	"github.com/gin-gonic/gin"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

func CheckErr(err error) {
	if err != nil {
		log.Println("error :", err)
		panic(err)
	}
}

func HandleErr(ctx *gin.Context, err interface{}) {
	exp := &types.ExceptionMessage{}
	_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
	responseSatus := PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
	ctx.JSON(http.StatusBadRequest, responseSatus)
}

func PrepareResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func PrepareResponseStatus(err interface{}) entity.ResponseStatus {
	return entity.ResponseStatus{
		IsSucccess: false,
		Message:    fmt.Sprint(err),
	}
}

func PrepareResponseStatusWithMessage(isSucccess bool, message string, code string, stack string) entity.ResponseStatus {
	return entity.ResponseStatus{
		IsSucccess: isSucccess,
		Message:    message,
		Code:       code,
		Stack:      stack,
	}
}

func GenerateRandomNumber(max int) int {
	return rand.Intn(max)
}

func GenerateGuid() string {
	return uuid.NewV4().String()
}

func DecodeBase64(encodeData string) ([]byte, error) {
	encodeData = strings.Replace(encodeData, "-", "+", -1)
	encodeData = strings.Replace(encodeData, "_", "/", -1)
	res, _ := b64.RawStdEncoding.DecodeString(encodeData)
	//res, _ := b64.StdEncoding.DecodeString(encodeData)

	return res, nil
}

func GetBase64PayloadFromJWT(jwtToken string) string {
	res := strings.Split(jwtToken, ".")[1]
	return res
}

func GetGoogleIdTokenSignKey(httpClient client.IHttpClient, idToken string) (string, error) {

	var googleOpenIDOAuthCertKey *entity.GoogleOpenIDOAuthCertKey
	googleJWTHeader, err := GetGoogleIdTokenHeaderInfo(idToken)
	if err != nil {
		log.Println("error :", err)
		return "", err
	}

	googleOpenIDOAuthCertKey, err = GetGoogleOpenIDOAuthCertKey(httpClient, googleJWTHeader)

	if err != nil {
		log.Println("error :", err)
		return "", err
	}

	if googleOpenIDOAuthCertKey == nil {
		return "", types.NewBusinessException("google idtoken sign key exception", "exp.google.id.token.sign.key")
	}

	return googleOpenIDOAuthCertKey.N, nil
}

func GetGoogleIdTokenHeaderInfo(idToken string) (*entity.GoogleJWTHeader, error) {
	jwtToken := GetJWTTokenInfo(idToken)
	googleJWTHeader := &entity.GoogleJWTHeader{}

	headerBytes, err := DecodeBase64(jwtToken.Header)

	if err != nil {
		log.Println("error :", err)
		return nil, err
	}

	err = json.Unmarshal(headerBytes, googleJWTHeader)

	if err != nil {
		log.Println("error :", err)
		return nil, err
	}

	return googleJWTHeader, nil
}

func GetJWTTokenInfo(jwtToken string) *entity.JWTToken {
	segments := strings.Split(jwtToken, ".")
	return &entity.JWTToken{Header: segments[0], Payload: segments[1], Signature: segments[2]}
}

func GetGoogleOpenIDConfiguration(httpClient client.IHttpClient) (*entity.GoogleOpenIDConfiguration, error) {
	conf := &entity.GoogleOpenIDConfiguration{}
	err := httpClient.Get("https://accounts.google.com/.well-known/openid-configuration").EndStruct(conf)
	return conf, err
}

func GetGoogleOpenIDOAuthCertKey(httpClient client.IHttpClient, jwtHeader *entity.GoogleJWTHeader) (*entity.GoogleOpenIDOAuthCertKey, error) {

	var googleOpenIDOAuthCertKey *entity.GoogleOpenIDOAuthCertKey
	googleOpenIDOAuthCertResponse, err := GetGoogleOpenIDOAuthCerts(httpClient)

	if err != nil {
		log.Println("error :", err)
		return nil, err
	}

	googleOpenIDOAuthCertKey = FindGoogleOpenIDOAuthCertKey(googleOpenIDOAuthCertResponse.Keys, jwtHeader)

	return googleOpenIDOAuthCertKey, nil
}

func GetGoogleOpenIDOAuthCerts(httpClient client.IHttpClient) (*entity.GoogleOpenIDOAuthCertResponse, error) {
	conf, err := GetGoogleOpenIDConfiguration(httpClient)
	if err != nil {
		log.Println("error :", err)
		return nil, err
	}
	certResponse := &entity.GoogleOpenIDOAuthCertResponse{}
	err = httpClient.Get(conf.JwksUri).EndStruct(certResponse)
	if err != nil {
		log.Println("error :", err)
		return nil, err
	}
	return certResponse, err
}

func FindGoogleOpenIDOAuthCertKey(certList []*entity.GoogleOpenIDOAuthCertKey, jwtHeader *entity.GoogleJWTHeader) *entity.GoogleOpenIDOAuthCertKey {
	var foundCert *entity.GoogleOpenIDOAuthCertKey
	if len(certList) > 0 {
		for _, cert := range certList {
			if cert.Alg == jwtHeader.Alg && cert.Kid == jwtHeader.KID {
				foundCert = cert
				break
			}
		}
	}
	return foundCert
}

func IsOneOf(list []string, item string) (error, bool) {

	if !slices.Contains(list, item) {
		return types.NewBusinessException("invalid item", "exp.util.not_in_list"), false
	}

	return nil, true
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GeneratePasswordHash(pwd string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error :", err)
		return "", err
	}

	return string(res), nil
}

func GenerateSecurePassword(length, numDigits, numSymbols int, noUpper, allowRepeat bool) (string, error) {
	res, err := password.Generate(length, numDigits, numSymbols, noUpper, allowRepeat)
	if err != nil {
		log.Println("error :", err)
		log.Println(err)
		return "", nil
	}
	return res, nil
}
