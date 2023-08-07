package service

import (
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"
)

func (service *CorporationService) GetCorporationToken(password string, email string) (*entity.Token, error) {

	corporationUserFromDb, err := service.corporationRepository.QueryX("CorporationUser", "select id,password,email from pa_corporation_users where email = $1 and status = $2;", email, 1)

	if err != nil {
		log.Println("error :", err)
		return nil, types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationUserFromDb == nil {
		return nil, types.NewBusinessException("user not found exception", "exp.user.notfound")
	}

	corporationUser := corporationUserFromDb.(*entity.CorporationUser)

	if !util.CheckPasswordHash(password, corporationUser.Password) {
		return nil, types.NewBusinessException("invalid credentials", "exp.user.invalid_credentials")
	}

	userId := corporationUser.Id

	// Embed User information to `token`
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entity.User{
		Id:    userId,
		Email: corporationUser.Email,
	})

	// token -> string. Only server knows this secret (foobar).
	tokenstring, err := jwtToken.SignedString([]byte(os.Getenv("PA_API_JWT_KEY")))

	if err != nil {
		log.Println("error :", err)
		return nil, types.NewBusinessException("google oauth2 token exception", "exp.google.oauth2.token")
	}

	tokenData := &entity.Token{AccessToken: tokenstring}

	return tokenData, err
}

func (service *CorporationService) UpdateCorporationLocation(user entity.User, id string, corporationId string, count int32) error {

	corporationUserFromDb, err := service.corporationRepository.QueryX("CorporationUser", "select * from pa_corporation_users where id = $1 and corporation_id = $2;", user.Id, corporationId)

	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationUserFromDb == nil {
		log.Printf("corporationUserFromDb is nil. id : %s, corporation_id : %s", user.Id, corporationId)
		return types.NewBusinessException("user not found exception", "exp.user.notfound")
	}

	corporationUser := corporationUserFromDb.(*entity.CorporationUser)

	corporationRoleFromDb, err := service.corporationRepository.QueryX("CorporationUserRole", "select * from pa_corporation_roles where id = $1;", corporationUser.RoleId)
	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationRoleFromDb == nil {
		log.Printf("corporationRoleFromDb is nil. role id : %s", corporationUser.RoleId)
		return types.NewBusinessException("user not found exception", "exp.userrole.notfound")
	}

	corporationUserRole := corporationRoleFromDb.(*entity.CorporationUserRole)

	if corporationUserRole.Name != "admin" {
		log.Printf("corporationUserRole.Name is not admin. role id : %s", corporationUser.RoleId)
		return types.NewBusinessException("user not found exception", "exp.userrole.invalid")
	}

	return service.corporationRepository.UpdateX("update pa_corporation_locations set available_location_count = $1 where id = $2 and corporation_id = $3;", count, id, corporationId)
}

func (service *CorporationService) InsertCorporationUser(user entity.User, email string, corporationId string) error {

	corporationUserFromDb, err := service.corporationRepository.QueryX("CorporationUser", "select * from pa_corporation_users where id = $1 and corporation_id = $2;", user.Id, corporationId)

	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationUserFromDb == nil {
		log.Printf("corporationUserFromDb is nil. id : %s, corpotation_id : %s", user.Id, corporationId)
		return types.NewBusinessException("user not found exception", "exp.user.notfound")
	}

	corporationUser := corporationUserFromDb.(*entity.CorporationUser)

	corporationRoleFromDb, err := service.corporationRepository.QueryX("CorporationUserRole", "select * from pa_corporation_roles where id = $1;", corporationUser.RoleId)
	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationRoleFromDb == nil {
		log.Printf("corporationRoleFromDb is nil. role id : %s", corporationUser.RoleId)
		return types.NewBusinessException("user not found exception", "exp.userrole.notfound")
	}

	corporationUserRole := corporationRoleFromDb.(*entity.CorporationUserRole)

	if corporationUserRole.Name != "admin" {
		log.Printf("corporationUserRole.Name is not admin. role id : %s", corporationUser.RoleId)
		return types.NewBusinessException("user not found exception", "exp.userrole.invalid")
	}

	corporationUserFromDb, err = service.corporationRepository.QueryX("CorporationUser", "select * from pa_corporation_users where email = $1;", email)

	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationUserFromDb != nil {
		log.Printf("corporationUserFromDb.Name is not nil. email : %s", email)
		return types.NewBusinessException("user already exist exception", "exp.user.already_exist")
	}

	//generate random password and send via mail
	pwd, err := util.GenerateSecurePassword(64, 10, 10, false, false)
	log.Printf("generated password for %s : %s", email, pwd)
	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	hashedPassword, err := util.GeneratePasswordHash(pwd)

	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	roleName := "user"
	corporationRoleFromDb, err = service.corporationRepository.QueryX("CorporationUserRole", "select * from pa_corporation_roles where name = $1;", roleName)
	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	if corporationRoleFromDb == nil {
		log.Printf("corporationRoleFromDb is nil. role name : %s", roleName)
		return types.NewBusinessException("user not found exception", "exp.userrole.notfound")
	}

	corporationUserRole = corporationRoleFromDb.(*entity.CorporationUserRole)

	err = service.corporationRepository.InsertX("insert into  pa_corporation_users (id,email,password,status,corporation_id,role_id) values($1,$2,$3,$4,$5,$6);", util.GenerateGuid(), email, hashedPassword, 1, corporationId, corporationUserRole.Id)

	if err != nil {
		log.Println("error :", err)
		return types.NewBusinessException("db exception", "exp.db.query.error")
	}

	return nil
}
