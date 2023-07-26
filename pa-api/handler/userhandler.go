package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/contract"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"
)

func (handler UserHandler) HandleOAuth2Google(ctx *gin.Context) {

	defer func() {
		if err := recover(); err != nil {
			util.HandleErr(ctx, err)
		}
	}()

	request := contract.GetGoogleOAuthTokenRequest{}

	if err := ctx.ShouldBindJSON(&request); err == nil {

		response, err := handler.userService.GetGoogleOAuthTokenResponse(request.Token)
		util.CheckErr(err)
		ctx.JSON(http.StatusOK, response)
	} else {
		exp := &types.ExceptionMessage{}
		_ = json.Unmarshal([]byte(fmt.Sprint(err)), exp)
		responseSatus := util.PrepareResponseStatusWithMessage(false, exp.Message, exp.Code, exp.Stack)
		ctx.JSON(http.StatusBadRequest, responseSatus)
	}
}
