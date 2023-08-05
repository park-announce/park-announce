package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/contract"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"

	"runtime/debug"
)

func (handler UserHandler) HandleOAuth2GoogleCode(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.GetGoogleOAuthTokenRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		response, err := handler.userService.GetGoogleOAuthCodeResponse(request.Token, request.ClientType)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler UserHandler) HandleOAuth2GoogleToken(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.GetGoogleOAuthTokenRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		response, err := handler.userService.GetGoogleOAuthTokenResponse(request.Token, request.ClientType)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler UserHandler) HandleOAuth2GoogleRegister(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.GetGoogleOAuthTokenRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		response, err := handler.userService.GetGoogleOAuthRegisterResponse(request.Token, request.ClientType)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler UserHandler) HandlePreRegister(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.PreregistrationRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		guid, err := handler.userService.SendOtp(request.Email)
		util.CheckErr(err)
		response := &contract.PreregistrationResponse{Guid: guid}
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler UserHandler) HandlePreRegisterVerification(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.ValidateAuthorizationCodeRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		err := handler.userService.ValidateOtp(request.Guid, request.Email, request.Otp)
		util.CheckErr(err)
		responseStatus := &entity.ResponseStatus{IsSucccess: true}
		ctx.JSON(http.StatusOK, responseStatus)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler UserHandler) HandleRegister(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.RegistrationRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {
		err := handler.userService.Register(request.Guid, request.Email, request.FirstName, request.LastName, request.MobilePhone, request.Password, request.CityCode)
		util.CheckErr(err)
		responseStatus := &entity.ResponseStatus{IsSucccess: true}
		ctx.JSON(http.StatusOK, responseStatus)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}
