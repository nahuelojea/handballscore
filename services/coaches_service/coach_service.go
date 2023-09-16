package coaches_service

import (
	"bytes"
	"context"
	"errors"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/coaches_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
	"github.com/nahuelojea/handballscore/storage"
)

const AvatarUrl = "avatars/coaches/"

type GetCoachsOptions struct {
	Name          string
	Surname       string
	Dni           string
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

func GetCoachs(filterOptions GetCoachsOptions) ([]models.Coach, int64, error) {
	filters := coaches_repository.GetCoachsOptions{
		Name:          filterOptions.Name,
		Surname:       filterOptions.Surname,
		Dni:           filterOptions.Dni,
		TeamId:        filterOptions.TeamId,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}

	return coaches_repository.GetCoachsFilteredAndPaginated(filters)
}

func UpdateCoach(coach models.Coach, ID string) (bool, error) {
	if len(coach.Dni) > 0 {
		_, exist, _ := GetCoachByDni(coach.AssociationId, coach.Dni)
		if exist {
			return false, errors.New("There is already a registered coach with this dni")
		}
	}
	return coaches_repository.UpdateCoach(coach, ID)
}

func DisableCoach(ID string) (bool, error) {
	return coaches_repository.DisableCoach(ID)
}

func GetCoachByDni(associationId, dni string) (models.Coach, bool, string) {
	return coaches_repository.GetCoachByDni(associationId, dni)
}

func GetAvatar(id string, ctx context.Context) (*bytes.Buffer, string, error) {
	coach, _, err := GetCoach(id)
	if err != nil {
		return nil, "", errors.New("Error to get coach: " + err.Error())
	}

	var filename = coach.Avatar
	if len(filename) < 1 {
		return nil, "", errors.New("The coach has no avatar")
	}

	file, err := storage.GetFile(ctx, filename)
	if err != nil {
		return nil, "", errors.New("Error to download file in S3 " + err.Error())
	}
	return file, filename, nil
}

func uploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var coach models.Coach

	filename = AvatarUrl + id + ".jpg"
	coach.Avatar = filename

	err := storage.UploadImage(ctx, contentType, body, filename)
	if err != nil {
		return errors.New("Error to upload image: " + err.Error())
	}

	status, err := coaches_repository.UpdateCoach(coach, id)
	if err != nil || !status {
		return errors.New("Error to update coach " + err.Error())
	}
	return nil
}
