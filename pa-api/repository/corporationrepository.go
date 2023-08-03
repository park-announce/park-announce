package repository

import (
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
		return nil, err
	}

	var instance entity.IEntity
	if result != nil && len(result) > 0 {
		instance = factory.GetEntityInstance(instanceType)
		convertMapToStruct(result[0], instance)
	}

	return instance, err
}

/*
func (repository *CorporationRepository) CheckCorporationUserRole(instanceType string, query string, roleId string) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query, roleId)
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

func (repository *CorporationRepository) CheckCorporationUser(instanceType string, query string, id string, corporationId string) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query, id, corporationId)
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

func (repository *CorporationRepository) CheckCorporationUserWithMail(instanceType string, query string, mail string) (interface{}, error) {

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

func (repository *CorporationRepository) ValidateCorporationUser(instanceType string, query string, email string, status int32) (interface{}, error) {

	result, err := Query(repository.GetConnection(), query, email, status)
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

*/

func (repository *CorporationRepository) UpdateX(query string, args ...interface{}) error {

	_, err := Update(repository.GetConnection(), query, args...)
	if err != nil {
		return err
	}

	return err
}

func (repository *CorporationRepository) InsertX(query string, args ...interface{}) error {

	_, err := Insert(repository.GetConnection(), query, args...)
	if err != nil {
		return err
	}

	return nil
}
