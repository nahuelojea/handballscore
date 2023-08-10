package dto

import (
	jwt "github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claim struct {
	Email         string             `json:"email"`
	Role          string             `json:"role"`
	AssociationId string             `json:"association_id"`
	Id            primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	jwt.RegisteredClaims
}
