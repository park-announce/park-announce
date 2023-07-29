package repository

import (
	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/factory"
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(dbClient *DBClient) UserRepository {
	return UserRepository{BaseRepository: BaseRepository{dbClient: dbClient}}
}

func (repository *UserRepository) GetByMail(instanceType string, mail string, query string) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query, mail)
	if err != nil {
		return nil, err
	}

	var instance entity.IEntity
	if result != nil && len(result) > 0 {
		instance = factory.GetEntityInstance(instanceType)
		convertMapToStruct(result[0], instance)
	}

	return instance, err
}
