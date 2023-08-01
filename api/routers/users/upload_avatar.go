package users

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/nahuelojea/handballscore/config/awsgo"
	"github.com/nahuelojea/handballscore/dto"
	"github.com/nahuelojea/handballscore/models"
	"github.com/nahuelojea/handballscore/repositories/users_repository"
)

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, request events.APIGatewayProxyRequest, claim dto.Claim) dto.RestResponse {

	var response dto.RestResponse
	response.Status = http.StatusBadRequest
	userId := claim.Id.Hex()

	var filename string
	var user models.User

	bucket := aws.String(ctx.Value(dto.Key("bucketName")).(string))
	filename = "users/" + userId + ".jpg"
	user.Avatar = filename

	mediaType, params, err := mime.ParseMediaType(request.Headers["Content-Type"])
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = err.Error()
		return response
	}

	if strings.HasPrefix(mediaType, "multipart/") {

		body, err := base64.StdEncoding.DecodeString(request.Body)
		if err != nil {
			response.Status = http.StatusInternalServerError
			response.Message = err.Error()
			return response
		}

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			response.Status = http.StatusInternalServerError
			response.Message = err.Error()
			return response
		}

		if err != io.EOF {
			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					response.Status = http.StatusInternalServerError
					response.Message = err.Error()
					return response
				}

				sess, err := session.NewSession(&aws.Config{
					Region: aws.String(awsgo.DefaultRegion)})

				if err != nil {
					response.Status = http.StatusInternalServerError
					response.Message = err.Error()
					return response
				}

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					response.Status = http.StatusInternalServerError
					response.Message = err.Error()
					return response
				}
			}
		}

		status, err := users_repository.UpdateUser(user, userId)
		if err != nil || !status {
			response.Status = http.StatusInternalServerError
			response.Message = "Error to update user " + err.Error()
			return response
		}
	} else {
		response.Message = "You must send an image with the 'Content-Type' of type 'multipart/' in the Header"
		response.Status = http.StatusBadRequest
		return response
	}

	response.Status = http.StatusOK
	response.Message = "Avatar uploaded"
	return response
}
