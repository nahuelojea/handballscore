package dto

type LoginResponse struct {
	Token         string `json:"token"`
	Avatar        string `json:"avatar"`
	AssociationId string `json:"association_id"`
}
