package models

const ImagesBaseURL = "https://handball-score.s3.sa-east-1.amazonaws.com/"

type Avatar interface {
	SetAvatarURL()
}
