package service

func (service *CorporationService) UpdateCorporationLocation(id string, corporationId string, count int32) error {
	return service.corporationRepository.UpdateCorporationLocationAvailabilityCount("update pa_corporation_locations set available_location_count = $1 where id = $2 and corporation_id = $3;", count, id, corporationId)
}
