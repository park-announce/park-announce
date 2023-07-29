package handler

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/service"
	"github.com/park-announce/pa-api/types"
	"github.com/park-announce/pa-api/util"
)

func (handler SocketHandler) HandleSocketConnection(ctx *gin.Context, hub *service.SocketHub) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(fmt.Sprintf("error in recover : %v, stack : %s", err, string(debug.Stack())))
			util.HandleErr(ctx, err)
		}
	}()

	userData, isExist := ctx.Get("User")
	var user entity.User
	if !isExist {
		panic(types.NewBusinessException("system exception", "exp.systemexception"))
	}

	user = userData.(entity.User)

	handler.socketService.CreateSocketConnection(ctx, user, hub)
}
