package players_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/players_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
	"github.com/nahuelojea/handballscore/services/categories_service"
	"github.com/nahuelojea/handballscore/storage"
)

const AvatarUrl = "avatars/players/"

type GetPlayersOptions struct {
	Name                    string
	Surname                 string
	Dni                     string
	Gender                  string
	OnlyEnabled             bool
	TeamId                  string
	AssociationId           string
	ExcludeExpiredInsurance bool
	CategoryId              string
	Page                    int
	PageSize                int
	SortField               string
	SortOrder               int
}

func CreatePlayer(association_id string, player models.Player) (string, bool, error) {
	_, exist, _ := teams_repository.GetTeam(player.TeamId)
	if !exist {
		return "", false, errors.New("No team found with this id")
	}

	_, exist, _ = players_repository.GetPlayerByDni(association_id, player.Dni)
	if exist {
		return "", false, errors.New("There is already a registered player with this dni")
	}
	return players_repository.CreatePlayer(association_id, player)
}

func GetPlayer(ID string) (models.Player, bool, error) {
	return players_repository.GetPlayer(ID)
}

func GetPlayers(filterOptions GetPlayersOptions) ([]models.Player, int64, int, error) {
	var gender = filterOptions.Gender

	var yearLimitFrom, yearLimitTo int
	var err error

	if len(filterOptions.CategoryId) > 1 {
		yearLimitFrom, yearLimitTo, gender, err = categories_service.GetLimitYearsByCategory(filterOptions.CategoryId)
		if err != nil {
			return nil, 0, 0, errors.New("Error to get category: " + err.Error())
		}
	}

	filters := players_repository.GetPlayersOptions{
		Name:                    filterOptions.Name,
		Surname:                 filterOptions.Surname,
		Dni:                     filterOptions.Dni,
		Gender:                  gender,
		OnlyEnabled:             filterOptions.OnlyEnabled,
		TeamId:                  filterOptions.TeamId,
		ExcludeExpiredInsurance: filterOptions.ExcludeExpiredInsurance,
		YearLimitFrom:           yearLimitFrom,
		YearLimitTo:             yearLimitTo,
		AssociationId:           filterOptions.AssociationId,
		Page:                    filterOptions.Page,
		PageSize:                filterOptions.PageSize,
		SortOrder:               filterOptions.SortOrder,
	}
	return players_repository.GetPlayers(filters)
}

func UpdatePlayer(player models.Player, ID string) (bool, error) {
	if len(player.Dni) > 0 {
		result, exist, _ := players_repository.GetPlayerByDni(player.AssociationId, player.Dni)
		if exist && player.Id != result.Id {
			return false, errors.New("There is already a registered player with this dni")
		}
	}
	return players_repository.UpdatePlayer(player, ID)
}

func DeletePlayer(ID string) (bool, error) {
	return players_repository.DeletePlayer(ID)
}

func GetPlayerByDni(association_id, dni string) (models.Player, bool, string) {
	return players_repository.GetPlayerByDni(association_id, dni)
}

func UploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var player models.Player

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename = fmt.Sprintf("%s%d_%s.jpg", AvatarUrl, timestamp, id)

	err := storage.UploadImage(ctx, contentType, body, filename)
	if err != nil {
		return errors.New("Error to upload image: " + err.Error())
	}

	player.SetAvatarURL(filename)
	status, err := players_repository.UpdateAvatar(player, id)
	if err != nil || !status {
		return errors.New("Error to update player " + err.Error())
	}

	return nil
}
