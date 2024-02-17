package coaches_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
	"github.com/nahuelojea/handballscore/storage"
)

const AvatarUrl = "avatars/coaches/"

type GetCoachesOptions struct {
	Name          string
	Surname       string
	Dni           string
	Gender        string
	OnlyEnabled   bool
	TeamId        string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateCoach(association_id string, coach models.Coach) (string, bool, error) {
	_, exist, _ := teams_repository.GetTeam(coach.TeamId)
	if !exist {
		return "", false, errors.New("No team found with this id")
	}

	_, exist, _ = GetCoachByDni(association_id, coach.Dni)
	if exist {
		return "", false, errors.New("There is already a registered coach with this dni")
	}
	return coaches_repository.CreateCoach(association_id, coach)
}

func GetCoach(ID string) (models.Coach, bool, error) {
	return coaches_repository.GetCoach(ID)
}

func GetCoaches(filterOptions GetCoachesOptions) ([]models.Coach, int64, error) {
	filters := coaches_repository.GetCoachesOptions{
		Name:          filterOptions.Name,
		Surname:       filterOptions.Surname,
		Dni:           filterOptions.Dni,
		Gender:        filterOptions.Gender,
		OnlyEnabled:   filterOptions.OnlyEnabled,
		TeamId:        filterOptions.TeamId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}

	return coaches_repository.GetCoaches(filters)
}

func UpdateCoach(coach models.Coach, ID string) (bool, error) {
	if len(coach.Dni) > 0 {
		result, exist, _ := GetCoachByDni(coach.AssociationId, coach.Dni)
		if exist && coach.Id != result.Id {
			return false, errors.New("There is already a registered coach with this dni")
		}
	}
	return coaches_repository.UpdateCoach(coach, ID)
}

func DeleteCoach(ID string) (bool, error) {
	return coaches_repository.DeleteCoach(ID)
}

func GetCoachByDni(associationId, dni string) (models.Coach, bool, string) {
	return coaches_repository.GetCoachByDni(associationId, dni)
}

func UploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var coach models.Coach

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename = fmt.Sprintf("%s%d_%s.jpg", AvatarUrl, timestamp, id)

	err := storage.UploadImage(ctx, contentType, body, filename)
	if err != nil {
		return errors.New("Error to upload image: " + err.Error())
	}

	coach.SetAvatarURL(filename)
	status, err := coaches_repository.UpdateAvatar(coach, id)
	if err != nil || !status {
		return errors.New("Error to update coach " + err.Error())
	}
	return nil
}
