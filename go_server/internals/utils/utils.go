package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	maxChunkSize = int64(5 << 20) // 5Mb
	maxReqBytes  = 1 << 20        // 1Mb
	maxRetries   = 3
)

type Envelope map[string]any

func InitS3Service(filename string) (*s3.S3, string, string, error) {
	bucket, region := LoadConfig()
	if filename == "" {
		return nil, "", "", errors.New("empty filename not allowed")
	}

	// Todo: Fix credentials
	sess := session.Must(session.NewSession())
	s3Svc := s3.New(sess, &aws.Config{
		Region: aws.String(region),
	})

	return s3Svc, bucket, filename, nil
}

func GetEnvInt(key string, fallback int) int {
	vatStr := os.Getenv(key)
	if vatStr == "" {
		return fallback
	}

	val, err := strconv.Atoi(vatStr)
	if err != nil {
		return fallback
	}

	return val
}

func WriteJSON(w http.ResponseWriter, _ *http.Request, status int, data Envelope) error {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	err := encoder.Encode(data)

	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxReqBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		var maxBytesError *http.MaxBytesError

		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func LoadConfig() (string, string) {
	bucket := os.Getenv("BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	if bucket == "" || region == "" {
		panic("BUCKET_NAME and AWS_REGION env var must be set")
	}

	return bucket, region
}

