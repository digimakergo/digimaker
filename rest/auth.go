//Author xc, Created on 2019-08-11 16:49
//{COPYRIGHTS}
package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xc/digimaker/core/handler"
	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

var key = "test123456"

type UserClaims struct {
	jwt.StandardClaims
	UserID int    `json:"user_id"`
	Name   string `json:"user_name"`
}

func newRefreshToken(userID int) (string, error) {
	refreshClaims := jwt.MapClaims{
		"user_id":      userID,
		"security_key": "222",
		"exp":          time.Now().Add(time.Minute * 60 * 5).Unix()} //todo: make it configurable
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	token, err := claims.SignedString([]byte(key)) //todo: make it configurable
	if err != nil {
		return "", err
	}
	//store it in db
	return token, nil
}

func newAccessToken(refreshToken string, r *http.Request) (string, error) {
	//check refresh token
	refreshClaims := struct {
		jwt.StandardClaims
		UserID      int    `json:"user_id"`
		SecurityKey string `json:"security_key"`
	}{}

	fmt.Println(refreshToken)
	token, err := jwt.ParseWithClaims(refreshToken, &refreshClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Invalid refresh token!")
	}

	securityKey := refreshClaims.SecurityKey

	userID := refreshClaims.UserID
	if securityKey != "222" {
		log.Warning("Someone is trying to use revoked token. ip: "+util.GetIP(r)+". user in the refresh token: "+strconv.Itoa(userID), "")
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
	atClaims := jwt.MapClaims{
		"user_id":   userID,
		"user_name": user.GetName(),
		"exp":       time.Now().Add(time.Minute * 5).Unix()} //todo: make it configurable

	atClaims1 := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessTokenKey := "fdsfdsfsfsdfsdf21"
	accessToken, err := atClaims1.SignedString([]byte(accessTokenKey))
	if err != nil {
		log.Error(err, "")
		return "", err
	}

	return accessToken, nil
}

func AuthGrant(w http.ResponseWriter, r *http.Request) {
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

func AuthRenewRefreshToken(w http.ResponseWriter, r *http.Request) {
}

func AuthRenewAccessToken(w http.ResponseWriter, r *http.Request) {
}

func init() {
	RegisterRoute("/auth/grant", AuthGrant)
	RegisterRoute("/auth/revoke", AuthRevoke)
	RegisterRoute("/auth/refresh", AuthRenewRefreshToken)
	RegisterRoute("/auth/rewnew_access", AuthRenewAccessToken)
}
