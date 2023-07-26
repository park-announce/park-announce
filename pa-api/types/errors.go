package types

import (
	"encoding/json"

	"github.com/go-errors/errors"
)

var MaxStackDepth = 50

type GameOverException struct {
	*BusinessException
}

type BusinessException struct {
	*BaseException
}

type BaseException struct {
	message string
	code    string
	base    *errors.Error
}

type ExceptionMessage struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Stack   string `json:"stack"`
}

func NewGameOverException(message string, code string) *GameOverException {

	a := &GameOverException{BusinessException: &BusinessException{BaseException: &BaseException{}}}
	x := errors.New(a)

	e := &GameOverException{BusinessException: &BusinessException{BaseException: &BaseException{}}}
	e.message = message
	e.code = code
	e.base = x

	return e
}

func NewBusinessException(message string, code string) *BusinessException {
	a := &BusinessException{BaseException: &BaseException{}}
	x := errors.New(a)

	e := &BusinessException{BaseException: &BaseException{}}
	e.message = message
	e.code = code
	e.base = x

	return e
}

func (e *GameOverException) Error() string {
	return e.GetExceptionMessage()
}

func (e *BusinessException) Error() string {
	return e.GetExceptionMessage()
}

func (e *BaseException) Error() string {
	return e.GetExceptionMessage()
}

func (e *BaseException) GetExceptionMessage() string {
	exp := &ExceptionMessage{Code: e.code, Message: e.message, Stack: string(e.base.Stack())}
	expBytes, _ := json.Marshal(exp)
	return string(expBytes)
}
