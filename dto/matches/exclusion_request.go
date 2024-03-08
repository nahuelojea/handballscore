package dto

type ExclusionRequest struct {
	Add  bool   `json:"add"`
	Time string `json:"time"`
}
