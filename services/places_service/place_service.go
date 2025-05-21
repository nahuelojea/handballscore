package places_service

import (
	"context"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/places_repository"
	"go.mongodb.org/mongo-driver/bson"
)

type PlaceService struct {
	Repository places_repository.PlaceRepository
}

func NewPlaceService() *PlaceService {
	return &PlaceService{
		Repository: *places_repository.NewPlaceRepository(),
	}
}

func (s *PlaceService) CreatePlace(ctx context.Context, place *models.Place) (string, bool, error) {
	return s.Repository.CreatePlace(ctx, place)
}

func (s *PlaceService) GetPlace(ctx context.Context, id string) (models.Place, bool, error) {
	return s.Repository.GetPlace(ctx, id)
}

func (s *PlaceService) GetPlaces(ctx context.Context, filter bson.M, page, pageSize int) ([]models.Place, int64, error) {
	return s.Repository.GetPlaces(ctx, filter, page, pageSize)
}

func (s *PlaceService) UpdatePlace(ctx context.Context, id string, place models.Place) (bool, error) {
	return s.Repository.UpdatePlace(ctx, id, place)
}

func (s *PlaceService) DeletePlace(ctx context.Context, id string) (bool, error) {
	return s.Repository.DeletePlace(ctx, id)
}
