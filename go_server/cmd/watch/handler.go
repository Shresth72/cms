package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/Shresth72/server/internals/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (app *application) getAllVideos(w http.ResponseWriter, r *http.Request) {
  metadata, err := app.models.MetadataModel.FindFiles()
  if err != nil {
    app.serverError(w, r, err, "error finding files from db")
    return
  }

  err = utils.WriteJSON(w, r, http.StatusOK, utils.Envelope{
    "data": metadata,
  })
  if err != nil {
    app.serverError(w, r, err)
    return
  }
}

func (app *application) getVideo(w http.ResponseWriter, r *http.Request) {
  videoKey := r.URL.Query().Get("key")
  if videoKey == "" {
    app.badRequestError(w, r, errors.New("key not present in query params")) 
  }
  
  signedUrl, signedHeaders, err := app.generateSignedUrl(videoKey)
  if err != nil {
    app.serverError(w, r, err) 
    return 
  }

  for key, values := range signedHeaders {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	err = utils.WriteJSON(w, r, http.StatusOK, utils.Envelope{
		"signed_url": signedUrl,
	})

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) generateSignedUrl(key string) (string, http.Header, error) {
  s3Svc, bucket, key, err := utils.InitS3Service(key)

	s3Params := s3.GetObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
	}

  resp, _ := s3Svc.GetObjectRequest(&s3Params)
  url, signedHeaders, err := resp.PresignRequest(60 * time.Minute)
  if err != nil {
    return "", nil, err
  }

  return url, signedHeaders, nil
}
