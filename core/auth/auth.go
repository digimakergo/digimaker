//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
)

type RefreshTokenManager interface {
	Store(id string, Expiry int64, claims map[string]interface{}) error
	Get(id string) interface{}
	Delete(id string) error
}

var tokenManager RefreshTokenManager

type RefreshClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	GUID   string `json:"guid"`
}

type UserClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	Name   string `json:"user_name"`
}

func NewRefreshToken(userID int) (string, error) {
	guid := util.GenerateGUID()
	refreshExpiry := util.GetConfigSectionI("auth")["refresh_token_expiry"].(int)
	expiry := time.Now().Add(time.Minute * time.Duration(refreshExpiry)).Unix()
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"guid":    guid,
		"exp":     expiry} //todo: make it configurable
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	//todo: better way to read configuration.
	refreshKey := util.GetConfigSectionI("auth")["refresh_token_secret_key"].(string)
	token, err := jwt.SignedString([]byte(refreshKey))
	if err != nil {
		return "", err
	}

	//store guid in db
	err = tokenManager.Store(guid, expiry, refreshClaims)
	if err != nil {
		log.Error(err.Error(), "")
		return "", errors.New("Error when storing refresh token info.")
	}
	return token, nil
}

func NewAccessToken(refreshToken string, r *http.Request) (string, error) {
	//check refresh token
	refreshClaims, err := VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}
	if refreshClaims.UserID == 0 {
		return "", errors.New("Invalid refresh token")
	}

	//generate new access token
	guid := refreshClaims.GUID
	userID := refreshClaims.UserID
	existingToken := tokenManager.Get(guid)
	if existingToken == nil {
		log.Warning("Someone is trying to use revoked token. guid: "+guid+" ip: "+util.GetIP(r)+". user in the refresh token: "+strconv.Itoa(userID), "")
		return "", errors.New("Invalid refresh token")
	}

	user, err := query.GetUser(userID)
	if user == nil || err != nil {
		if err != nil {
			log.Error(err.Error(), "")
		}
		return "", errors.New("User not found")
	}

	//Generate new access token
	accessExpiry := util.GetConfigSectionI("auth")["access_token_expiry"].(int)
	accessClaims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": user.GetName(),
		"exp":       time.Now().Add(time.Minute * time.Duration(accessExpiry)).Unix()}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessKey := util.GetConfigSectionI("auth")["access_token_secret_key"].(string)
	accessToken, err := jwt.SignedString([]byte(accessKey))
	if err != nil {
		log.Error(err, "")
		return "", err
	}

	return accessToken, nil
}

//if failed there will be always err
func VerifyRefreshToken(token string) (RefreshClaims, error) {
	claims := RefreshClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Wrong signing method: %v", token.Header["alg"])
		}
		refreshKey := util.GetConfigSectionI("auth")["refresh_token_secret_key"].(string)
		return []byte(refreshKey), nil
	})
	if err != nil {
		return RefreshClaims{}, err
	}
	if jwtToken.Valid {
		entity := tokenManager.Get(claims.GUID)
		if entity == nil {
			return claims, TokenErrorRevoked
		}
		return claims, nil
	} else {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch ve.Errors {
			case jwt.ValidationErrorExpired:
				return RefreshClaims{}, TokenErrorExpired
			default:
				return RefreshClaims{}, err
			}
		}
		return RefreshClaims{}, err
	}
}

var (
	TokenErrorExpired = errors.New("Expired token")
	TokenErrorRevoked = errors.New("Expired revoked")
)

//Verify access token, return nil, TokenErrorExpired or other err
//@todo: maybe store refresh's guid in access token also to check if it's there? It will have access token revoked in refresh token is revoked.
func VerifyToken(token string) (error, UserClaims) {
	accessClaims := UserClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &accessClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Wrong signing method: %v", token.Header["alg"])
		}
		accessKey := util.GetConfigSectionI("auth")["access_token_secret_key"].(string)
		return []byte(accessKey), nil
	})

	if jwtToken.Valid {
		return nil, accessClaims
	} else {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch ve.Errors {
			case jwt.ValidationErrorExpired:
				return TokenErrorExpired, UserClaims{}
			default:
				return err, UserClaims{}
			}
		} else {
			return err, UserClaims{}
		}
	}
}

func GetTokenManager() RefreshTokenManager {
	return tokenManager
}

func RegisterTokenManager(manager RefreshTokenManager) {
	tokenManager = manager
}
