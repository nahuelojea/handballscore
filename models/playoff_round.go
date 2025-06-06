package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ThirtyTwoFinals = "thirty_two_finals"
	SixteenFinals   = "sixteen_finals"
	EightFinals     = "eight_finals"
	QuarterFinals   = "quarter_finals"
	SemiFinal       = "semi_final"
	Final           = "final"
	ThirdPlace      = "third_place"
	FifthPlace      = "fifth_place"
	SeventhPlace    = "seventh_place"
	NinthPlace      = "ninth_place"
)

type PlayoffRound struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Round          string             `bson:"round" json:"round"`
	TeamsQuantity  int                `bson:"teams_quantity" json:"teams_quantity"`
	PlayoffPhaseId string             `bson:"playoff_phase_id" json:"playoff_phase_id"`
	Status_Data    `bson:"status_data" json:"status_data"`
	AssociationId  string `bson:"association_id" json:"association_id"`
}

func (playoffRound *PlayoffRound) SetAssociationId(associationId string) {
	playoffRound.AssociationId = associationId
}

func (playoffRound *PlayoffRound) SetCreatedDate() {
	playoffRound.CreatedDate = time.Now()
}

func (playoffRound *PlayoffRound) SetModifiedDate() {
	playoffRound.ModifiedDate = time.Now()
}

func (playoffRound *PlayoffRound) SetId(id primitive.ObjectID) {
	playoffRound.Id = id
}

func (playoffRound *PlayoffRound) PlayoffRoundNameTraduction() string {
	switch playoffRound.Round {
	case ThirtyTwoFinals:
		return "Treintaidosavos de Final"
	case SixteenFinals:
		return "Dieciseisavos de Final"
	case EightFinals:
		return "Octavos de Final"
	case QuarterFinals:
		return "Cuartos de Final"
	case SemiFinal:
		return "Semifinal"
	case Final:
		return "Final"
	case ThirdPlace:
		return "Tercer Puesto"
	case FifthPlace:
		return "Quinto Puesto"
	case SeventhPlace:
		return "Septimo Puesto"
	case NinthPlace:
		return "Noveno Puesto"
	default:
		return "Fase Desconocida"
	}
}

func GetRoundFromTeamsCount(teamsCount int) string {
	switch {
	case teamsCount <= 2:
		return Final
	case teamsCount <= 4:
		return SemiFinal
	case teamsCount <= 8:
		return QuarterFinals
	case teamsCount <= 16:
		return EightFinals
	case teamsCount <= 32:
		return SixteenFinals
	default:
		return ThirtyTwoFinals
	}
}

func GetNextRound(currentRound string) (string, error) {
	switch currentRound {
	case ThirtyTwoFinals:
		return SixteenFinals, nil
	case SixteenFinals:
		return EightFinals, nil
	case EightFinals:
		return QuarterFinals, nil
	case QuarterFinals:
		return SemiFinal, nil
	case SemiFinal:
		return Final, nil
	case Final:
		return "", errors.New("there is no next round after the Final")
	default:
		return "", errors.New("unknown round")
	}
}
