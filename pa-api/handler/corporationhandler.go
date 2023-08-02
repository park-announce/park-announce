package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/contract"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"
)

func (handler CorporationHandler) HandleCorporationLocationUpdate(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.CorporationLocationUpdateRequest{}
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&request); err == nil {

		userData, isExist := ctx.Get("User")
		var user entity.User
		if !isExist {
			panic(types.NewBusinessException("system exception", "exp.systemexception"))
		}

		user = userData.(entity.User)

		err := handler.corporationService.UpdateCorporationLocation(user, id, request.CorporationId, request.Count)
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

func (handler CorporationHandler) HandleCorporationToken(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.CorporationOAuthTokenRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {
		response, err := handler.corporationService.GetCorporationToken(request.Password, request.Email)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}

func (handler CorporationHandler) HandleCorporationUserInsert(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.CorporationUserInsertRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		userData, isExist := ctx.Get("User")
		var user entity.User
		if !isExist {
			panic(types.NewBusinessException("system exception", "exp.systemexception"))
		}

		user = userData.(entity.User)

		err := handler.corporationService.InsertCorporationUser(user, request.Email, request.CorporationId)
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
