package service

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/park-announce/pa-api/contract"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"

	jwt "github.com/dgrijalva/jwt-go"
)

var mobileClientTypes = []string{"ios", "android"}
var webClientTypes = []string{"web"}
var clientTypes = append(mobileClientTypes, webClientTypes...)

var clientIds = map[string]string{"web": os.Getenv("PA_API_WEB_GOOGLE_CLIENT_ID"), "ios": os.Getenv("PA_API_IOS_GOOGLE_CLIENT_ID"), "android": os.Getenv("PA_API_ANDROID_GOOGLE_CLIENT_ID")}

func (service *UserService) GetGoogleOAuthCodeResponse(code string, clientType string) (*entity.Token, error) {
	var err error
	var tokenstring string

	err, valid := util.IsOneOf(webClientTypes, clientType)

	if !valid {
		return nil, types.NewBusinessException("google oauth2 client_type exception", "exp.google.oauth2.clint_type")
	}

	data := url.Values{}
	data.Add("code", code)
	data.Add("client_id", os.Getenv("PA_API_WEB_GOOGLE_CLIENT_ID"))
	data.Add("client_secret", os.Getenv("PA_API_WEB_GOOGLE_CLIENT_SECRET"))
	data.Add("redirect_uri", os.Getenv("PA_API_WEB_GOOGLE_REDIRECT_URI"))
	data.Add("grant_type", "authorization_code")

	googleTokenResponse := &contract.GetGoogleOAuthTokenResponse{}
	err = service.httpClient.PostUrlEncoded("https://www.googleapis.com/oauth2/v4/token", data).EndStruct(googleTokenResponse)

	if err != nil {
		fmt.Println("err -> ", err)
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	user := entity.GoogleUser{}

	jwt.ParseWithClaims(googleTokenResponse.IdToken, &user, func(token *jwt.Token) (interface{}, error) {

		key, err := util.GetGoogleIdTokenSignKey(service.httpClient, googleTokenResponse.IdToken)

		if err != nil {
			return nil, err
		}

		return []byte(key), nil
	})

	if user.Audience != clientIds[clientType] {
		return nil, types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.aud")
	}

	// Embed User information to `token`
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Picture:   user.Picture,
	})
	// token -> string. Only server knows this secret (foobar).
	tokenstring, err = token.SignedString([]byte(os.Getenv("PA_API_JWT_KEY")))

	if err != nil {
		fmt.Println("err -> ", err.Error())
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}

func (service *UserService) GetGoogleOAuthTokenResponse(idToken string, clientType string) (*entity.Token, error) {
	var err error
	var tokenstring string

	err, valid := util.IsOneOf(clientTypes, clientType)

	if !valid {
		return nil, types.NewBusinessException("google oauth2 client_type exception", "exp.google.oauth2.clint_type")
	}

	user := entity.GoogleUser{}

	jwt.ParseWithClaims(idToken, &user, func(token *jwt.Token) (interface{}, error) {

		key, err := util.GetGoogleIdTokenSignKey(service.httpClient, idToken)

		if err != nil {
			return nil, err
		}

		return []byte(key), nil
	})

	//TODO:token'in expire olup olmadigi kontrol edilecek.

	if user.Audience != clientIds[clientType] {
		return nil, types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.aud")
	}

	//check is user already register

	dbUser, err := service.userRepository.GetByMail("User", user.Email, "select id from pa_users where email = $1")

	if err != nil {
		return nil, types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if dbUser == nil {
		return nil, types.NewBusinessException("user not found exception", "exp.user.notfound")
	}

	//TODO:status alani kontrol edilecek, yasakli kullanici ise hata verilecek.

	userId := dbUser.(*entity.User).Id

	// Embed User information to `token`
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id:        userId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Picture:   user.Picture,
	})

	// token -> string. Only server knows this secret (foobar).
	tokenstring, err = jwtToken.SignedString([]byte(os.Getenv("PA_API_JWT_KEY")))

	if err != nil {
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}

func (service *UserService) GetGoogleOAuthRegisterResponse(idToken string, clientType string) (*entity.Token, error) {
	var err error
	var tokenstring string

	err, valid := util.IsOneOf(clientTypes, clientType)

	if !valid {
		return nil, types.NewBusinessException("google oauth2 client_type exception", "exp.google.oauth2.clint_type")
	}

	user := entity.GoogleUser{}

	jwt.ParseWithClaims(idToken, &user, func(token *jwt.Token) (interface{}, error) {

		key, err := util.GetGoogleIdTokenSignKey(service.httpClient, idToken)

		if err != nil {
			return nil, err
		}

		return []byte(key), nil
	})

	if user.Audience != clientIds[clientType] {
		return nil, types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.aud")
	}

	//check is user already register

	dbUser, err := service.userRepository.GetByMail("User", user.Email, "select id from pa_users where email = $1")

	if err != nil {
		return nil, types.NewBusinessException("db exception", "exp.db.query.error")
	}

	var userId string = util.GenerateGuid()
	if dbUser == nil {
		_, err := service.userRepository.Insert("insert into pa_users(id,email) values($1,$2)", userId, user.Email)

		if err != nil {
			return nil, types.NewBusinessException("db exception", "exp.db.query.error")
		}
	} else {
		userId = dbUser.(*entity.User).Id
	}

	// Embed User information to `token`
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id:        userId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Picture:   user.Picture,
	})
	// token -> string. Only server knows this secret (foobar).
	tokenstring, err = jwtToken.SignedString([]byte(os.Getenv("PA_API_JWT_KEY")))

	if err != nil {
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}

func (service *UserService) SendOtp(email string) (string, error) {

	//TODO:generate authorizatin code and write it to redis and send as mail
	authorizationCode := 1234

	guid := util.GenerateGuid()
	_, err := service.redisClient.Set(fmt.Sprintf("preregister-mail|%s", guid), email, time.Minute*15)

	if err != nil {
		return "", types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	_, err = service.redisClient.Set(fmt.Sprintf("preregister-otp|%s", guid), authorizationCode, time.Minute*3)

	if err != nil {
		return "", types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	return guid, nil

}

func (service *UserService) ValidateOtp(guid string, email string, otp string) error {

	sendedMail, err := service.redisClient.Get(fmt.Sprintf("preregister-mail|%s", guid))

	if err != nil {
		return types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	if sendedMail != email {
		return types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	sendOtp, err := service.redisClient.Get(fmt.Sprintf("preregister-otp|%s", guid))

	if sendOtp != otp {
		return types.NewBusinessException("invalid authorization code", "exp.register.invalid_otp")
	}

	return nil

}

func (service *UserService) Register(guid, email, firstName, lastName, mobilePhone, password string, city int16) error {

	sendedMail, err := service.redisClient.Get(fmt.Sprintf("preregister-mail|%s", guid))

	if err != nil {
		return types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	if sendedMail != email {
		return types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	userId := util.GenerateGuid()
	//TODO :check is there any existing record with this email
	_, err = service.userRepository.Insert("insert into  pa_users (id,email,first_name,last_name,status,mobile_phone,city_code) values($1,$2,$3,$4,$5,$6,$7);", userId, email, firstName, lastName, 1, mobilePhone, city)

	if err != nil {
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	hashedPassword, err := util.GeneratePasswordHash(password)

	if err != nil {
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	_, err = service.userRepository.Insert("insert into  pa_user_passwords (id,user_id,password) values($1,$2,$3);", util.GenerateGuid(), userId, hashedPassword)

	if err != nil {
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	return nil

}
