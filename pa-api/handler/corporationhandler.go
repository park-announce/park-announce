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
		err := handler.corporationService.UpdateCorporationLocation(id, request.CorporationId, request.Count)
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
