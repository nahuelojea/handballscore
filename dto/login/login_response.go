package dto

type LoginResponse struct {
	Token         string `json:"token"`
	RefreshToken  string `json:"refresh_token"`
	Avatar        string `json:"avatar"`
	AssociationId string `json:"association_id"`
}
