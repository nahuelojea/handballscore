package teams_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/teams_repository"
	"github.com/nahuelojea/handballscore/storage"
)

const AvatarUrl = "avatars/teams/"

type GetTeamsOptions struct {
	Name          string
	AssociationId string
	Page          int
	PageSize      int
	SortField     string
	SortOrder     int
}

func CreateTeam(association_id string, team models.Team) (string, bool, error) {
	return teams_repository.CreateTeam(association_id, team)
}

func GetTeam(ID string) (models.Team, bool, error) {
	return teams_repository.GetTeam(ID)
}

func GetTeams(filterOptions GetTeamsOptions) ([]models.Team, int64, error) {
	filters := teams_repository.GetTeamsOptions{
		Name:          filterOptions.Name,
		AssociationId: filterOptions.AssociationId,
		Page:          filterOptions.Page,
		PageSize:      filterOptions.PageSize,
		SortField:     filterOptions.SortField,
		SortOrder:     filterOptions.SortOrder,
	}
	return teams_repository.GetTeams(filters)
}

func UpdateTeam(team models.Team, ID string) (bool, error) {
	return teams_repository.UpdateTeam(team, ID)
}

func DeleteTeam(ID string) (bool, error) {
	return teams_repository.DeleteTeam(ID)
}

func UploadAvatar(ctx context.Context, contentType, body, id string) error {
	var filename string
	var team models.Team

	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	filename = fmt.Sprintf("%s%d_%s.jpg", AvatarUrl, timestamp, id)

	err := storage.UploadImage(ctx, contentType, body, filename)
	if err != nil {
		return errors.New("Error to upload image: " + err.Error())
	}

	team.SetAvatarURL(filename)
	status, err := teams_repository.UpdateAvatar(team, id)
	if err != nil || !status {
		return errors.New("Error to update team " + err.Error())
	}
	return nil
}
