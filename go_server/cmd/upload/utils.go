package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/oauth2"
  "github.com/golang-jwt/jwt/v5"
)

type envelope map[string]any

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

// ResponseWriter Utils
func (app *application) writeJSON(w http.ResponseWriter, _ *http.Request, status int, data envelope) error {
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
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

func loadConfig() (string, string) {
	bucket := os.Getenv("BUCKET_NAME")
	region := os.Getenv("AWS_REGION")

	if bucket == "" || region == "" {
		panic("BUCKET_NAME and AWS_REGION env var must be set")
	}

	return bucket, region
}

// Multipart Utils
func (app *application) uploadPart(sv *s3.S3, params s3.UploadPartInput) (*s3.CompletedPart, error) {
	tryNum := 1
	partNumber := params.PartNumber

	for tryNum <= maxRetries {
		uploadResult, err := sv.UploadPart(&params)
		if err != nil {
			if tryNum == maxRetries {
				if _, ok := err.(awserr.Error); ok {
					return nil, err
				}
				return nil, err
			}
			app.logMsgFmt("Retrying to upload part %d", partNumber)
			tryNum++
		} else {
			app.logMsgFmt("Uploaded part %d", partNumber)
			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: partNumber,
			}, nil
		}
	}

	return nil, nil
}

func (app *application) initS3Service(filename string) (*s3.S3, string, string, error) {
	bucket, region := loadConfig()
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

// OAuth2 Utils
func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	userInfoEndpoint := "https://www.googleapis.com/oauth2/v2/userinfo"

	resp, err := http.Get(fmt.Sprintf("%s?access_token=%s", userInfoEndpoint, accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func SignJWT(userInfo map[string]interface{}) (string, string, error) {
	claims := jwt.MapClaims{
		"sub":   userInfo["id"],
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
		"name":  userInfo["name"],
		"email": userInfo["email"],
		"iss":   "oauth2-golang",
	}
  id := fmt.Sprintf("%v", userInfo["id"])

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")

	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", "", err
	}

	return signedToken, id, nil
}

func (app *application) saveRefreshToken(userId string, refreshToken string) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh_token_%s", userId)
	err := app.redis.Set(ctx, key, refreshToken, 0).Err()
	return err
}

func (app *application) refreshToken(refreshToken string) (*oauth2.Token, error) {
	tokenSource := app.oauth.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}

	return newToken, nil
}

type UserInfoClaims struct {
  Sub string `json:"sub"`
  Exp string `json:"exp"`
  Name string `json:"name"`
  Email string `json:"email"`
  Iss string `json:"iss"`
}

func (app *application) VerifyJWT(tokenString string) (*UserInfoClaims, error) {
  token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
      if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
      }
      hmacSecret := []byte(os.Getenv("SECRET_KEY"))
      return hmacSecret, nil
    })

    if err != nil {
      return nil, err
    }

   if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &UserInfoClaims{
			Sub:   claims["sub"].(string),
			Exp:   claims["exp"].(string),
			Name:  claims["name"].(string),
			Email: claims["email"].(string),
			Iss:   claims["iss"].(string),
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}
