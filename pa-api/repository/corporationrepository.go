package repository

import (
	"log"

	"github.com/park-announce/pa-api/entity"
	"github.com/park-announce/pa-api/factory"
)

type CorporationRepository struct {
	BaseRepository
}

func NewCorporationRepository(dbClient *DBClient) CorporationRepository {
	return CorporationRepository{BaseRepository: BaseRepository{dbClient: dbClient}}
}

func (repository *CorporationRepository) QueryX(instanceType string, query string, args ...interface{}) (interface{}, error) {

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

func (repository *CorporationRepository) UpdateX(query string, args ...interface{}) error {

	_, err := Update(repository.GetConnection(), query, args...)
	if err != nil {
		log.Println("error :", err)
		return err
	}

	return err
}

func (repository *CorporationRepository) InsertX(query string, args ...interface{}) error {

	_, err := Insert(repository.GetConnection(), query, args...)
	if err != nil {
		log.Println("error :", err)
		return err
	}

	return nil
}
