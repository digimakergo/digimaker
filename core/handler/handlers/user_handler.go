//Author xc, Created on 2020-10-04 17:19
//{COPYRIGHTS}

//Package handlers implements build-in action callbacks.
package handlers

import (
	"context"
	"fmt"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/spf13/viper"
)

type UserHandler struct {
}

func (uh UserHandler) ValidateCreate(ctx context.Context, inputs handler.InputMap, parentID int) (bool, handler.ValidationResult) {
	login := fmt.Sprint(inputs["login"])
	email := fmt.Sprint(inputs["email"])

	result := handler.ValidationResult{Fields: map[string]string{}}

	existing, _ := query.Fetch(context.Background(), "user", db.Cond("login", login))
	if existing != nil {
		result.Fields["login"] = "Username is used already"
	}

	if viper.GetBool("user.unique_email") {
		existing, _ = query.Fetch(context.Background(), "user", db.Cond("email", email))
		if existing != nil {
			result.Fields["email"] = "Email is used already"
		}
	}

	return result.Passed(), result
}

func (uh UserHandler) ValidateUpdate(ctx context.Context, inputs handler.InputMap, content contenttype.ContentTyper) (bool, handler.ValidationResult) {
	login := fmt.Sprint(inputs["login"])
	email := fmt.Sprint(inputs["email"])

	result := handler.ValidationResult{Fields: map[string]string{}}

	loginField := content.Value("login").(string)
	if loginField != login {
		existing, _ := query.Fetch(context.Background(), "user", db.Cond("login", login))
		if existing != nil && existing.GetCID() != content.GetCID() { //NB. uppcase change is allowed
			result.Fields["login"] = "Username is used already"
		}
	}

	emailField := content.Value("email").(string)
	if emailField != email {
		if viper.GetBool("user.unique_email") {
			existing, _ := query.Fetch(context.Background(), "user", db.Cond("email", email))
			if existing != nil && existing.GetCID() != content.GetCID() { //NB. uppcase change is allowed
				result.Fields["email"] = "Email is used already"
			}
		}
	}

	return result.Passed(), result
}

func init() {
	handler.RegisterContentTypeHandler("user", UserHandler{})
}
