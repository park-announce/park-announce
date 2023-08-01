package handler

import (
	"github.com/park-announce/pa-api/service"
)

type SocketHandler struct {
	socketService service.SocketService
}

type UserHandler struct {
	userService service.UserService
}

type CorporationHandler struct {
	corporationService service.CorporationService
}

func NewSocketHandler(socketService service.SocketService) SocketHandler {
	return SocketHandler{socketService: socketService}
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}

func NewCorporationHandler(corporationService service.CorporationService) CorporationHandler {
	return CorporationHandler{corporationService: corporationService}
}
