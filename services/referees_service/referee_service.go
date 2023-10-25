package referees_service

import (
	"context"
	"errors"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/referees_repository"
	"github.com/nahuelojea/handballscore/storage"
)

const AvatarUrl = "avatars/referees/"

type GetRefereesOptions struct {
	Name                    string
	Surname                 string
	Dni                     string
	Gender                  string
	OnlyEnabled             bool
	TeamId                  string
	AssociationId           string
	ExcludeExpiredInsurance bool
	YearLimitFrom           int
	YearLimitTo             int
	Page                    int
	PageSize                int
	SortField               string
	SortOrder               int
}

func CreateReferee(association_id string, referee models.Referee) (string, bool, error) {
	_, exist, _ := referees_repository.GetRefereeByDni(association_id, referee.Dni)
	if exist {
		return "", false, errors.New("There is already a registered referee with this dni")
	}
	return referees_repository.CreateReferee(association_id, referee)
}

func GetReferee(ID string) (models.Referee, error) {
	return referees_repository.GetReferee(ID)
}

func GetReferees(filterOptions GetRefereesOptions) ([]models.Referee, int64, error) {
	filters := referees_repository.GetRefereesOptions{
		Name:          filterOptions.Name,
		Surname:       filterOptions.Surname,
		Dni:           filterOptions.Dni,
		Gender:        filterOptions.Gender,
		OnlyEnabled:   filterOptions.OnlyEnabled,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortOrder:     filterOptions.SortOrder,
	}
	return referees_repository.GetReferees(filters)
}

func UpdateReferee(referee models.Referee, ID string) (bool, error) {
	if len(referee.Dni) > 0 {
		result, exist, _ := referees_repository.GetRefereeByDni(referee.AssociationId, referee.Dni)
		if exist && referee.Id != result.Id {
			return false, errors.New("There is already a registered referee with this dni")
		}
	}
	return referees_repository.UpdateReferee(referee, ID)
}

func DeleteReferee(ID string) (bool, error) {
	return referees_repository.DeleteReferee(ID)
}

func GetRefereeByDni(association_id, dni string) (models.Referee, bool, string) {
	return referees_repository.GetRefereeByDni(association_id, dni)
}

func UploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var referee models.Referee

	filename = AvatarUrl + id + ".jpg"

	err := storage.UploadImage(ctx, contentType, body, filename)
	if err != nil {
		return errors.New("Error to upload image: " + err.Error())
	}

	referee.SetAvatarURL(filename)
	status, err := referees_repository.UpdateAvatar(referee, id)
	if err != nil || !status {
		return errors.New("Error to update referee " + err.Error())
	}
	return nil
}
