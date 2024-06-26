package main

import (
	"bytes"
	"net/http"

	"github.com/Shresth72/server/internals/data"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type InitMultipartRequest struct {
	Filename string `json:"filename"`
}

type UploadChunkRequest struct {
	Filename   string `json:"filename"`
	ChunkIndex int    `json:"chunk_index"`
	UploadId   string `json:"upload_id"`
}

type CompleteMultipartRequest struct {
	Filename    string   `json:"filename"`
	UploadId    string   `json:"upload_id"`
	Title       string   `json:"title"`
	Description string   `json:"desc"`
	Etags       []string `json:"e_tag"`
}

func (app *application) initMultipartUpload(w http.ResponseWriter, r *http.Request) {
	var body InitMultipartRequest
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	s3Svc, bucket, key, err := app.initS3Service(body.Filename)
	if err != nil {
		app.serverError(w, r, err)
	}

	s3Params := s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String("video/mp4"),
	}

	resp, err := s3Svc.CreateMultipartUpload(&s3Params)
	if err != nil {
		app.serverError(w, r, err, "error creating multipart upload request")
		return
	}
	app.logSuccess(r, "Created multipart upload request")

	err = app.writeJSON(w, r, http.StatusOK, envelope{
		"upload_id": resp.UploadId,
	})

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) uploadChunk(w http.ResponseWriter, r *http.Request) {
	var body UploadChunkRequest
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	fileBytes, ok := r.Context().Value("fileBytes").([]byte)
	if !ok || len(fileBytes) == 0 {
		app.errorResponse(w, r, http.StatusBadRequest, "missing file bytes")
		return
	}

	s3Svc, bucket, key, err := app.initS3Service(body.Filename)
	if err != nil {
		app.serverError(w, r, err)
	}

	partParams := s3.UploadPartInput{
		Bucket:     aws.String(bucket),
		Key:        aws.String(key),
		UploadId:   &body.UploadId,
		PartNumber: aws.Int64(int64(body.ChunkIndex) + 1),
		Body:       bytes.NewReader(fileBytes),
	}

	completedPart, err := app.uploadPart(s3Svc, partParams)
	if err != nil {
		app.serverError(w, r, err, "error uploading chunk, aborting...")

		abortParams := s3.AbortMultipartUploadInput{
			Bucket:   aws.String(bucket),
			Key:      aws.String(key),
			UploadId: &body.UploadId,
		}
		_, err := s3Svc.AbortMultipartUpload(&abortParams)
		if err != nil {
			app.serverError(w, r, err, "error aborting multipart upload")
		}
		return
	}

	app.logSuccess(r, "Created multipart upload request")

	w.Header().Set("Etag", *completedPart.ETag)
	err = app.writeJSON(w, r, http.StatusOK, envelope{
		// "e_tag":       completedPart.ETag,
		"part_number": completedPart.PartNumber,
	})

	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) completeMultipart(w http.ResponseWriter, r *http.Request) {
	var body CompleteMultipartRequest
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	s3Svc, bucket, key, err := app.initS3Service(body.Filename)
	if err != nil {
		app.serverError(w, r, err)
	}

	var completedParts []*s3.CompletedPart
	for i, etag := range body.Etags {
		completedParts = append(completedParts, &s3.CompletedPart{
			ETag:       aws.String(etag),
			PartNumber: aws.Int64(int64(i + 1)),
		})
	}

	completeParams := s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		UploadId: &body.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}

	completedResponse, err := s3Svc.CompleteMultipartUpload(&completeParams)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logMsgFmt("Completed multipart upload to S3: %s\n", completedResponse.String())

	err = app.writeJSON(w, r, http.StatusCreated, envelope{
		"msg": "successfully completed multipart upload",
	})
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Add Details to DB
  userInfo, ok := r.Context().Value("userInfo").(*UserInfoClaims)
	if !ok {
		app.errorResponse(w, r, http.StatusBadRequest, "user information not found, unauthorized access or some error")
		return
	}

  details := data.Metadata{
    Key: *completedResponse.Key,
    Title: body.Title,
    Description: body.Description,
    Author: userInfo.Email,
    Url: *completedResponse.Location,
  }

  err = app.models.MetadataModel.CreateFile(&details)
  if err != nil {
    app.serverError(w, r, err, "Could not save details to db, try again")
    return 
  }

	// Push for Encoding using Kafka

  err = app.writeJSON(w, r, http.StatusOK, envelope{
		"data": details,
	})
	if err != nil {
		app.serverError(w, r, err)
	}
}

// Delete with Incomplete Deletion for fast backup
func (app *application) uploadToDb(w http.ResponseWriter, r *http.Request) {

}
