package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MatchHeaderView struct {
	Id                     primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Date                   time.Time          `bson:"date" json:"date"`
	TeamHomeId             primitive.ObjectID `bson:"team_home_id" json:"team_home_id"`
	TeamHomeVariant        string             `bson:"team_home_variant" json:"team_home_variant"`
	TeamHomeName           string             `bson:"team_home_name" json:"team_home_name"`
	TeamHomeAvatar         string             `bson:"team_home_avatar" json:"team_home_avatar"`
	TeamHomeInitials       string             `bson:"team_home_initials" json:"team_home_initials"`
	TeamHomeMainColor      string             `bson:"team_home_main_color" json:"team_home_main_color"`
	TeamHomeSecondaryColor string             `bson:"team_home_secondary_color" json:"team_home_secondary_color"`
	TeamAwayId             primitive.ObjectID `bson:"team_away_id" json:"team_away_id"`
	TeamAwayVariant        string             `bson:"team_away_variant" json:"team_away_variant"`
	TeamAwayName           string             `bson:"team_away_name" json:"team_away_name"`
	TeamAwayAvatar         string             `bson:"team_away_avatar" json:"team_away_avatar"`
	TeamAwayInitials       string             `bson:"team_away_initials" json:"team_away_initials"`
	TeamAwayMainColor      string             `bson:"team_away_main_color" json:"team_away_main_color"`
	TeamAwaySecondaryColor string             `bson:"team_away_secondary_color" json:"team_away_secondary_color"`
	Referees               []string           `bson:"referees" json:"referees"`
	Place                  string             `bson:"place" json:"place"`
	Status                 string             `bson:"status" json:"status"`
	StreamingUrl           string             `bson:"streaming_url" json:"streaming_url"`
	GoalsHome              int                `bson:"goals_home" json:"goals_home"`
	GoalsAway              int                `bson:"goals_away" json:"goals_away"`
	TournamentCategoryId   string             `bson:"tournament_category_id" json:"tournament_category_id"`
	TournamentCategoryName string             `bson:"tournament_category_name" json:"tournament_category_name"`
	LeaguePhaseWeekId      string             `bson:"league_phase_week_id" json:"league_phase_week_id"`
	PlayoffRoundKeyId      string             `bson:"playoff_round_key_id" json:"playoff_round_key_id"`
	PlayoffRound           string             `bson:"playoff_round" json:"playoff_round"`
	Category_Id            string             `bson:"category_id" json:"category_id"`
	AssociationId          string             `bson:"association_id" json:"association_id"`
	Status_Data            `bson:"status_data" json:"status_data"`
}

func (matchView *MatchHeaderView) SetCreatedDate() {
	matchView.CreatedDate = time.Now()
}

func (matchView *MatchHeaderView) SetModifiedDate() {
	matchView.ModifiedDate = time.Now()
}

func (matchView *MatchHeaderView) SetAssociationId(associationId string) {
	matchView.AssociationId = associationId
}

func (matchView *MatchHeaderView) SetId(id primitive.ObjectID) {
	matchView.Id = id
}
