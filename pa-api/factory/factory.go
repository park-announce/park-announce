package factory

import (
	"github.com/park-announce/pa-api/entity"
)

var factoryList = make(map[string]IFactory)

func InitFactoryList() {
	factoryList["User"] = UserFactory{}
	factoryList["CorporationUser"] = CorporationUserFactory{}
	factoryList["CorporationUserRole"] = CorporationUserRoleFactory{}

}

type IFactory interface {
	GetInstance() entity.IEntity
}

type UserFactory struct {
}

func (userFactory UserFactory) GetInstance() entity.IEntity {
	return &entity.User{}
}

type CorporationUserFactory struct {
}

func (corporationUserFactory CorporationUserFactory) GetInstance() entity.IEntity {
	return &entity.CorporationUser{}
}

type CorporationUserRoleFactory struct {
}

func (corporationUserRoleFactory CorporationUserRoleFactory) GetInstance() entity.IEntity {
	return &entity.CorporationUserRole{}
}

func GetEntityInstance(name string) entity.IEntity {
	return factoryList[name].GetInstance()
}
