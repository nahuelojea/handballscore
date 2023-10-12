package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/nahuelojea/handballscore/config/awsgo"
	"github.com/nahuelojea/handballscore/dto"
)

const maxImageSize = 600 * 1024

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, contentType, body, filename string) error {
	bucket := aws.String(ctx.Value(dto.Key("bucketName")).(string))

	fmt.Println("ContentType: > " + contentType)

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}

	fmt.Println("MediaType: > " + mediaType)

	if strings.HasPrefix(mediaType, "multipart/") {

		body, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return err
		}

		fmt.Println("Llegue 1")

		imageSize := len(body)
		if imageSize > maxImageSize {
			return errors.New(fmt.Sprintf("Image size exceeds the maximum allowed size of %d bytes", maxImageSize))
		}

		fmt.Println("Llegue 2")

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			return err
		}

		fmt.Println("Llegue 3")

		if err != io.EOF {
			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					return err
				}

				fmt.Println("Llegue 4")

				sess, err := session.NewSession(&aws.Config{
					Region: aws.String(awsgo.DefaultRegion)})

				if err != nil {
					return err
				}

				fmt.Println("Llegue 5")

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				fmt.Println("Llegue 6")

				if err != nil {
					return err
				}
			}
		}
	} else {
		return errors.New("You must send an image with the 'Content-Type' of type 'multipart/' in the Header")
	}
	return nil
}

func GetFile(ctx context.Context, filename string) (*bytes.Buffer, error) {
	svc := s3.NewFromConfig(awsgo.Cfg)

	bucket := ctx.Value(dto.Key("bucketName")).(string)
	obj, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()

	file, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(file)

	return buffer, nil
}
