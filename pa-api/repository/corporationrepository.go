package repository

type CorporationRepository struct {
	BaseRepository
}

func NewCorporationRepository(dbClient *DBClient) CorporationRepository {
	return CorporationRepository{BaseRepository: BaseRepository{dbClient: dbClient}}
}

func (repository *CorporationRepository) UpdateCorporationLocationAvailabilityCount(query string, count int32, id string, corporationId string) error {

	_, err := Update(repository.GetConnection(), query, count, id, corporationId)
	if err != nil {
		return err
	}

	return err
}
