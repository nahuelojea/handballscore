package players_service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/players_repository"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
	"github.com/nahuelojea/handballscore/services/categories_service"
	"github.com/nahuelojea/handballscore/storage"
	"github.com/xuri/excelize/v2"
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

func ImportFromExcel(fileContent string) (bool, []error) {
	decodedFile, err := base64.StdEncoding.DecodeString(fileContent)
	if err != nil {
		return false, []error{errors.New("Invalid file format")}
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(decodedFile))
	if err != nil {
		return false, []error{errors.New("Cannot read the excel file")}
	}

	rows, err := excelFile.GetRows("Jugadores")
	if err != nil {
		return false, []error{errors.New("Cannot get rows from the excel file")}
	}

	var errorsList []error

	for i, row := range rows[1:] {
		dateOfBirth, err := time.Parse("02/01/2006", row[2])
		if err != nil {
			errorsList = append(errorsList, fmt.Errorf("Invalid date of birth format in row %d: %v", i+2, row))
			continue
		}

		expirationInsurance, err := time.Parse("02/01/2006", "01/03/2025")
		if err != nil {
			errorsList = append(errorsList, fmt.Errorf("Invalid expiration insurance date format in row %d: %v", i+2, row))
			continue
		}

		fmt.Println("Expiration insurance: ", expirationInsurance)

		player := models.Player{
			Personal_Data: models.Personal_Data{
				Name:        row[0],
				Surname:     row[1],
				DateOfBirth: dateOfBirth,
				Dni:         row[3],
				Gender:      row[4],
				PhoneNumber: row[5],
				Disabled:    false,
			},
			AffiliateNumber:     row[6],
			TeamId:              row[9],
			ExpirationInsurance: expirationInsurance,
			AssociationId:       row[10],
		}
		_, _, err = CreatePlayer(player.AssociationId, player)
		if err != nil {
			errorsList = append(errorsList, fmt.Errorf("Failed to create player for row %d: %v", i+2, row))
			continue
		}
	}

	if len(errorsList) > 0 {
		return false, errorsList
	}
	return true, nil
}
