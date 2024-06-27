package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

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
