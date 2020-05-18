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
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

type RefreshTokenManager interface {
	Store(id string, Expiry int64) error
	Get(id string) interface{}
	Delete(id string) error
}

var refreshKey = "testtesttest11111"
var accessKey = "testtest22222"
var tokenManager RefreshTokenManager

type UserClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	Name   string `json:"user_name"`
}

func newRefreshToken(userID int) (string, error) {
	guid := util.GenerateGUID()
	expiry := time.Now().Add(time.Minute * 60 * 5).Unix()
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"guid":    guid,
		"exp":     expiry} //todo: make it configurable
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	token, err := jwt.SignedString([]byte(refreshKey)) //todo: make it configurable
	if err != nil {
		return "", err
	}

	//store meta info
	err = tokenManager.Store(guid, expiry)
	if err != nil {
		log.Error(err.Error(), "")
		return "", errors.New("Error when storing refresh token info.")
	}
	//store it in db
	return token, nil
}

func newAccessToken(refreshToken string, r *http.Request) (string, error) {
	//check refresh token
	refreshClaims := struct {
		jwt.StandardClaims
		UserID int    `json:"user_id"`
		GUID   string `json:"guid"`
	}{}

	token, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid refresh token!")
	}

	guid := refreshClaims.GUID

	userID := refreshClaims.UserID
	existingToken := tokenManager.Get(guid)
	if existingToken == nil {
		log.Warning("Someone is trying to use revoked token. guid: "+guid+" ip: "+util.GetIP(r)+". user in the refresh token: "+strconv.Itoa(userID), "")
		return "", errors.New("Invalid refresh token!")
	}

	user, err := handler.Querier().GetUser(userID)
	if user == nil || err != nil {
		if err != nil {
			log.Error(err.Error(), "")
		}
		return "", errors.New("User not found.")
	}

	//Generate new access token
	accessClaims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": user.GetName(),
		"exp":       time.Now().Add(time.Minute * 5).Unix()} //todo: make it configurable

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := jwt.SignedString([]byte(accessKey))
	if err != nil {
		log.Error(err, "")
		return "", err
	}

	return accessToken, nil
}

//Grant  refresh toke and access token
func AuthAuthenticate(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
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

	userID := user.GetCID()
	refreshToken, err := newRefreshToken(userID)
	if err != nil {
		log.Error("Error in generating token on "+strconv.Itoa(userID)+": "+err.Error(), "")
		HandleError(errors.New("Error when generating refresh token."), w)
		return
	}

	accessToken, err := newAccessToken(refreshToken, r)
	if err != nil {
		log.Error("Error in generating token on "+strconv.Itoa(userID)+": "+err.Error(), "")
		HandleError(errors.New("Error when generating refresh token."), w)
		return
	}

	resultObj := map[string]interface{}{"refresh_token": refreshToken, "access_token": accessToken}
	result, _ := json.Marshal(resultObj)

	w.Write(result)
}

func AuthRevoke(w http.ResponseWriter, r *http.Request) {

}

func Auth(w http.ResponseWriter, r *http.Request) (bool, error) {
	authStr := r.Header.Get("Authorization")
	if authStr == "" {
		return false, errors.New("Empty Authentication")
	}
	authSlice := strings.Split(authStr, " ")
	if authSlice[0] != "Bearer" {
		return false, errors.New("Only bearer is supported.")
	}

	token := authSlice[1]
	accessClaims := UserClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &accessClaims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Wrong signing method: %v", token.Header["alg"])
		}
		return []byte(accessKey), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if jwtToken.Valid {
		return true, nil
	} else {
		return false, nil
	}

}

func AuthVerify(w http.ResponseWriter, r *http.Request) {
	verified, err := Auth(w, r)
	if err != nil {
		HandleError(err, w)
		return
	}
	if verified {
		w.Write([]byte("1"))
	} else {
		w.Write([]byte("0"))
	}
}

func AuthRenewRefreshToken(w http.ResponseWriter, r *http.Request) {
}

func AuthRenewAccessToken(w http.ResponseWriter, r *http.Request) {
}

func RegisterTokenManager(manager RefreshTokenManager) {
	tokenManager = manager
}

func init() {
	RegisterRoute("/auth/verify", AuthVerify)
	RegisterRoute("/auth/auth", AuthAuthenticate)
	RegisterRoute("/auth/revoke", AuthRevoke)
	RegisterRoute("/auth/token/refresh", AuthRenewRefreshToken)
	RegisterRoute("/auth/token/rewnew_access", AuthRenewAccessToken)
}
