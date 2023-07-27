package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

var (
	AccessExpiresAt  = jwt.NewNumericDate(time.Now().Add(time.Minute * 15))
	RefreshExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))
	IssuedAt         = jwt.NewNumericDate(time.Now())
	NotBefore        = jwt.NewNumericDate(time.Now())
	Secret           = []byte(os.Getenv("secret"))
)

type Claims struct {
	Username string
	ID       int
	jwt.RegisteredClaims
}

type TokenDetail struct {
	Username  string
	UserId    int
	TokenID   string
	Token     *string
	ExpiresAt time.Time
	Subject   string
}

type Token struct {
	AccessToken  *TokenDetail
	RefreshToken *TokenDetail
}

type LoginResponse struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateLoginResponse(id int, username string) (*LoginResponse, *Token, error) {
	tokenID := uuid.NewV4().String()
	accessTokenDetail, err := GenerateToken(id, username, tokenID, "access_token")
	if err != nil {
		return nil, nil, err
	}
	refreshTokenDetail, err := GenerateToken(id, username, tokenID, "refresh_token")
	if err != nil {
		return nil, nil, err
	}
	return &LoginResponse{
			Username:     username,
			AccessToken:  *accessTokenDetail.Token,
			RefreshToken: *refreshTokenDetail.Token,
		}, &Token{
			AccessToken:  accessTokenDetail,
			RefreshToken: refreshTokenDetail,
		}, nil

}

func GenerateToken(id int, username, tokenID, subject string) (*TokenDetail, error) {
	claims := &Claims{}
	if subject == "access_token" {
		claims = &Claims{
			Username: username,
			ID:       id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: AccessExpiresAt,
				IssuedAt:  IssuedAt,
				NotBefore: NotBefore,
				Subject:   subject,
				ID:        tokenID,
			},
		}
	}
	if subject == "refresh_token" {
		claims = &Claims{
			Username: username,
			ID:       id,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: RefreshExpiresAt,
				IssuedAt:  IssuedAt,
				NotBefore: NotBefore,
				Subject:   subject,
				ID:        tokenID,
			},
		}
	}
	if subject == "" {
		return nil, errors.New("no subject")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(Secret)
	if err != nil {
		return nil, err
	}
	return &TokenDetail{
		Username:  claims.Username,
		UserId:    claims.ID,
		TokenID:   claims.RegisteredClaims.ID,
		Token:     &signedToken,
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Time,
		Subject:   claims.RegisteredClaims.Subject,
	}, nil
}
