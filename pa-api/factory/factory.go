package factory

import (
	"github.com/park-announce/pa-api/entity"
)

var factoryList = make(map[string]IFactory)

func InitFactoryList() {
	factoryList["User"] = UserFactory{}
}

type IFactory interface {
	GetInstance() entity.IEntity
}

type UserFactory struct {
}

func (userFactory UserFactory) GetInstance() entity.IEntity {
	return &entity.User{}
}

func GetEntityInstance(name string) entity.IEntity {
	return factoryList[name].GetInstance()
}
