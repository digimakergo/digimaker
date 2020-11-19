//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
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

func newRefreshToken(userID int) (string, error) {
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

func newAccessToken(refreshToken string, r *http.Request) (string, error) {
	//check refresh token
	refreshClaims, err := verifyRefreshToken(refreshToken)
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

	user, err := handler.Querier().GetUser(userID)
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

//AuthAuthenticate generate refresh toke and access token based on username and password
func AuthAuthenticate(w http.ResponseWriter, r *http.Request) {
	//Check matches
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")
	if username == "" || password == "" {
		HandleError(errors.New("Please input username or password"), w)
		return
	}
	err, user := handler.CanLogin(username, password)
	if err != nil {
		log.Error("Error when verifying username and password of "+username+": "+err.Error(), "", r.Context())
		HandleError(errors.New("Verifying failed"), w, 400)
		return
	}

	//Generate refresh token and access token
	userID := user.GetCID()
	refreshToken, err := newRefreshToken(userID)
	if err != nil {
		log.Error("Error in generating token on "+strconv.Itoa(userID)+": "+err.Error(), "")
		HandleError(errors.New("Error when generating refresh token"), w)
		return
	}

	accessToken, err := newAccessToken(refreshToken, r)
	if err != nil {
		log.Error("Error in generating token on "+strconv.Itoa(userID)+": "+err.Error(), "")
		HandleError(errors.New("Error when generating refresh token"), w)
		return
	}

	resultObj := map[string]interface{}{"refresh_token": refreshToken, "access_token": accessToken}
	result, _ := json.Marshal(resultObj)

	w.Write(result)
}

func AuthRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	//Verify refresh token and delete.
	token := r.FormValue("token")
	claims, err := verifyRefreshToken(token)
	if err != nil {
		HandleError(err, w)
		return
	}
	if claims.UserID == 0 {
		HandleError(errors.New("No valid token"), w)
		return
	}

	guid := claims.GUID
	err = tokenManager.Delete(guid)
	if err != nil {
		log.Error("Deleting token error: "+err.Error(), "", r.Context())
		return
	}

	w.Write([]byte("1"))
}

func getToken(r *http.Request) (string, error) {
	authStr := r.Header.Get("Authorization")
	if authStr == "" {
		return "", errors.New("Empty Authentication header")
	}
	authSlice := strings.Split(authStr, " ")
	if len(authSlice) != 2 {
		return "", errors.New("Wrong format of bearer")
	}
	if authSlice[0] != "Bearer" {
		return "", errors.New("Only bearer is supported")
	}
	return authSlice[1], nil
}

//if failed there will be always err
func verifyRefreshToken(token string) (RefreshClaims, error) {
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
func VerifyToken(r *http.Request) (error, UserClaims) {
	token, err := getToken(r)
	if err != nil {
		return err, UserClaims{}
	}
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

func AuthVerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	err, _ := VerifyToken(r)
	if err != nil {
		HandleError(err, w, StatusUnauthed)
		return
	}
	w.Write([]byte("1"))
}

//Renew refresh token
func AuthRenewRefreshToken(w http.ResponseWriter, r *http.Request) {
	//verify refresh token
	token := r.FormValue("token")
	if token == "" {
		HandleError(errors.New("Token parameter is needed"), w, StatusUnauthed)
		return
	}

	refreshClaims, err := verifyRefreshToken(token)
	if err != nil {
		log.Error(err.Error(), "", r.Context())
		if err == TokenErrorExpired || err == TokenErrorRevoked {
			HandleError(err, w, StatusExpired)
			return
		}
		HandleError(errors.New("Invalid token"), w, StatusUnauthed)
		return
	}

	//generate new refresh token.
	userID := refreshClaims.UserID
	guid := refreshClaims.GUID
	newToken, err := newRefreshToken(userID)
	if err != nil {
		HandleError(err, w)
		return
	}
	err = tokenManager.Delete(guid)
	if err != nil {
		log.Error("Error when deleting token: "+err.Error(), "", r.Context())
		HandleError(err, w)
	}

	w.Write([]byte(newToken))
}

//Renew access token
func AuthRenewAccessToken(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	if token == "" {
		HandleError(errors.New("Token parameter is needed"), w)
		return
	}
	refreshClaims, err := verifyRefreshToken(token)
	if err != nil || refreshClaims.GUID == "" {
		if err != nil {
			log.Error(err.Error(), "", r.Context())
		}
		HandleError(errors.New("Not valid refresh token"), w, StatusUnauthed)
		return
	}

	accessToken, err := newAccessToken(token, r)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.Write([]byte(accessToken))
}

func RegisterTokenManager(manager RefreshTokenManager) {
	tokenManager = manager
}

func init() {
	RegisterRoute("/auth/auth", AuthAuthenticate)
	RegisterRoute("/auth/token/revoke", AuthRevokeRefreshToken)
	RegisterRoute("/auth/token/refresh/renew", AuthRenewRefreshToken)
	RegisterRoute("/auth/token/access/renew", AuthRenewAccessToken)
	RegisterRoute("/auth/token/access/verify", AuthVerifyAccessToken)
}
