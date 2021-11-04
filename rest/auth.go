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

	"github.com/digimakergo/digimaker/core/auth"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/log"
	"github.com/golang-jwt/jwt"

	"github.com/digimakergo/digimaker/core/util"
)

type AuthInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	RememberMe bool   `json:"rememberme"`
}

//AuthAuthenticate generate refresh toke and access token based on username and password
func AuthAuthenticate(w http.ResponseWriter, r *http.Request) {
	//Check matches
	input := AuthInput{}
	decorder := json.NewDecoder(r.Body)
	err := decorder.Decode(&input)
	if err != nil {
		HandleError(err, w)
		return
	}

	username := input.Username
	password := input.Password
	if username == "" || password == "" {
		HandleError(errors.New("Please input username or password"), w)
		return
	}
	err, user := handler.CanLogin(r.Context(), username, password)
	if err != nil {
		log.Error("Error when verifying username and password of "+username+": "+err.Error(), "", r.Context())
		HandleError(errors.New("Verifying failed"), w, 400)
		return
	}

	//Generate refresh token and access token
	userID := user.GetCID()
	rememberMe := false
	if util.GetConfigSectionI("auth")["enable_rememberme"].(bool) {
		rememberMe = input.RememberMe
	}
	refreshToken, err := auth.NewRefreshToken(r.Context(), userID, rememberMe)
	if err != nil {
		log.Error("Error in generating token on "+strconv.Itoa(userID)+": "+err.Error(), "")
		HandleError(errors.New("Error when generating refresh token"), w)
		return
	}

	accessToken, err := auth.NewAccessToken(refreshToken, r)
	if err != nil {
		log.Error("Error in generating token on "+strconv.Itoa(userID)+": "+err.Error(), "")
		HandleError(errors.New("Error when generating refresh token"), w)
		return
	}

	resultObj := map[string]interface{}{"refresh_token": refreshToken, "access_token": accessToken}
	WriteResponse(resultObj, w)
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
	err = auth.GetTokenManager().Delete(r.Context(), guid)
	if err != nil {
		log.Error("Deleting token error: "+err.Error(), "", r.Context())
		return
	}

	WriteResponse(true, w)
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
func verifyRefreshToken(token string) (auth.RefreshClaims, error) {
	claims := auth.RefreshClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Wrong signing method: %v", token.Header["alg"])
		}
		refreshKey := util.GetConfigSectionI("auth")["refresh_token_secret_key"].(string)
		return []byte(refreshKey), nil
	})
	if err != nil {
		return auth.RefreshClaims{}, err
	}
	if jwtToken.Valid {
		entity := auth.GetTokenManager().Get(claims.GUID)
		if entity == nil {
			return claims, auth.TokenErrorRevoked
		}
		return claims, nil
	} else {
		if ve, ok := err.(*jwt.ValidationError); ok {
			switch ve.Errors {
			case jwt.ValidationErrorExpired:
				return auth.RefreshClaims{}, auth.TokenErrorExpired
			default:
				return auth.RefreshClaims{}, err
			}
		}
		return auth.RefreshClaims{}, err
	}
}

func VerifyAccessToken(r *http.Request) (error, auth.UserClaims) {
	token, err := getToken(r)
	if err != nil {
		return err, auth.UserClaims{}
	}

	err, claims := auth.VerifyToken(token)
	if err != nil {
		return err, auth.UserClaims{}
	}
	return nil, claims
}

func AuthVerifyAccessToken(w http.ResponseWriter, r *http.Request) {
	err, _ := VerifyAccessToken(r)
	if err != nil {
		HandleError(err, w, StatusUnauthed)
		return
	}
	WriteResponse(true, w)
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

	accessToken, err := auth.NewAccessToken(token, r)
	if err != nil {
		HandleError(err, w)
		return
	}

	WriteResponse(accessToken, w)
}

func init() {
	RegisterRoute("/auth/auth", AuthAuthenticate, "POST")
	RegisterRoute("/auth/token/revoke", AuthRevokeRefreshToken)
	RegisterRoute("/auth/token/access/renew", AuthRenewAccessToken)
	RegisterRoute("/auth/token/access/verify", AuthVerifyAccessToken)
}
