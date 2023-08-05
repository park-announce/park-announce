package repository

import (
	"log"

	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/factory"
)

type UserRepository struct {
	BaseRepository
}

func NewUserRepository(dbClient *DBClient) UserRepository {
	return UserRepository{BaseRepository: BaseRepository{dbClient: dbClient}}
}

func (repository *UserRepository) QueryX(instanceType string, query string, args ...interface{}) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query, args...)
	if err != nil {
		log.Println("error :", err)
		return nil, err
	}

	var instance entity.IEntity
	if result != nil && len(result) > 0 {
		instance = factory.GetEntityInstance(instanceType)
		convertMapToStruct(result[0], instance)
	}

	return instance, err
}
