package service

import (
	"fmt"
	"os"
	"time"

	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"

	jwt "github.com/dgrijalva/jwt-go"
)

var mobileClientTypes = []string{"ios", "android"}
var webClientTypes = []string{"web"}
var clientTypes = append(mobileClientTypes, webClientTypes...)
var clientIds = map[string]string{"web": os.Getenv("PA_API_WEB_GOOGLE_CLIENT_ID"), "ios": os.Getenv("PA_API_IOS_GOOGLE_CLIENT_ID"), "android": os.Getenv("PA_API_ANDROID_GOOGLE_CLIENT_ID")}

type UserStatusType int16

const (
	Passive UserStatusType = 0
	Active  UserStatusType = 1
	Blocked UserStatusType = 2
)

func (service *UserService) GetGoogleOAuthTokenResponse(idToken string, clientType string) (*entity.Token, error) {
	var err error
	var tokenstring string

	err, valid := util.IsOneOf(clientTypes, clientType)

	if !valid {
		return nil, types.NewBusinessException("google oauth2 client_type exception", "exp.google.oauth2.clint_type")
	}

	googleUser := entity.GoogleUser{}

	jwt.ParseWithClaims(idToken, &googleUser, func(token *jwt.Token) (interface{}, error) {

		key, err := util.GetGoogleIdTokenSignKey(service.httpClient, idToken)

		if err != nil {
			return nil, err
		}

		return []byte(key), nil
	})

	if googleUser.Audience != clientIds[clientType] {
		return nil, types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.invalid_aud")
	}

	//check is token expired
	now := time.Now().Unix()
	if now > googleUser.ExpiresAt {
		return nil, types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.token_expired")
	}

	//check is user registered
	dbUser, err := service.userRepository.QueryX("User", googleUser.Email, "select id from pa_users where email = $1")

	if err != nil {
		return nil, types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if dbUser == nil {
		return nil, types.NewBusinessException("user not found exception", "exp.user.notfound")
	}

	user := dbUser.(*entity.User)

	//check user status
	if user.Status != int16(Active) {
		return nil, types.NewBusinessException("user not found exception", "exp.user.notfound")
	}

	userId := user.Id

	// Embed User information to `token`
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id:        userId,
		FirstName: googleUser.FirstName,
		LastName:  googleUser.LastName,
		Email:     googleUser.Email,
		Picture:   googleUser.Picture,
	})

	// token -> string. Only server knows this secret (foobar).
	tokenstring, err = jwtToken.SignedString([]byte(os.Getenv("PA_API_JWT_KEY")))

	if err != nil {
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}

func (service *UserService) GetGuidForGoogleRegistration(idToken string, clientType string) (string, error) {

	err, valid := util.IsOneOf(clientTypes, clientType)

	if !valid {
		return "", types.NewBusinessException("google oauth2 client_type exception", "exp.google.oauth2.clint_type")
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
		return "", types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.aud")
	}

	//check is token expired
	now := time.Now().Unix()
	if now > user.ExpiresAt {
		return "", types.NewBusinessException("google oauth2 aud exception", "exp.google.oauth2.token_expired")
	}

	//check is user already register
	dbUser, err := service.userRepository.QueryX("User", user.Email, "select id from pa_users where email = $1")

	if err != nil {
		return "", types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if dbUser != nil {
		return "", types.NewBusinessException("user not found exception", "exp.user.is_already_registered")
	}

	guid := util.GenerateGuid()

	_, err = service.redisClient.Set(fmt.Sprintf("preregister-guid|%s", guid), user.Email, time.Minute*15)

	if err != nil {
		return "", types.NewBusinessException("register with google", "exp.redis.otp.error")
	}

	return guid, err
}

func (service *UserService) SendOtp(email string) (string, error) {

	//check is email exist

	userFromDb, err := service.userRepository.QueryX("User", "select * from pa_users where email = $1;", email)

	if err != nil {
		return "", types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if userFromDb != nil {
		return "", types.NewBusinessException("user already exist exception", "exp.user.already_exist")
	}

	//TODO:generate authorizatin code and write it to redis and send as mail
	authorizationCode := 1234

	guid := util.GenerateGuid()
	_, err = service.redisClient.Set(fmt.Sprintf("preregister-guid|%s", guid), email, time.Minute*15)

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

	sendedMail, err := service.redisClient.Get(fmt.Sprintf("preregister-guid|%s", guid))

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

	sendedMail, err := service.redisClient.Get(fmt.Sprintf("preregister-guid|%s", guid))

	if err != nil {
		return types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	if sendedMail != email {
		return types.NewBusinessException("register with mail", "exp.redis.otp.error")
	}

	userId := util.GenerateGuid()
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
