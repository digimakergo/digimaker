//Author xc, Created on 2020-10-04 17:19
//{COPYRIGHTS}

//Package handlers implements build-in action callbacks.
package handlers

import (
	"fmt"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/fieldtype"
	"github.com/digimakergo/digimaker/core/handler"
	"github.com/digimakergo/digimaker/core/util"
)

type UserHandler struct {
}

func (uh UserHandler) ValidateCreate(inputs handler.InputMap, parentID int) (bool, handler.ValidationResult) {
	login := fmt.Sprint(inputs["login"])
	email := fmt.Sprint(inputs["email"])

	result := handler.ValidationResult{Fields: map[string]string{}}
	querier := handler.Querier()

	existing, _ := querier.Fetch("user", db.Cond("login", login))
	if existing != nil {
		result.Fields["login"] = "Username is used already"
	}

	if util.GetConfigSectionI("user")["unique_email"].(bool) {
		existing, _ = querier.Fetch("user", db.Cond("email", email))
		if existing != nil {
			result.Fields["email"] = "Email is used already"
		}
	}

	return result.Passed(), result
}

func (uh UserHandler) ValidateUpdate(inputs handler.InputMap, content contenttype.ContentTyper) (bool, handler.ValidationResult) {
	login := fmt.Sprint(inputs["login"])
	email := fmt.Sprint(inputs["email"])

	result := handler.ValidationResult{Fields: map[string]string{}}
	querier := handler.Querier()

	loginField := content.Value("login").(*fieldtype.Text)
	if loginField.String.String != login {
		existing, _ := querier.Fetch("user", db.Cond("login", login))
		if existing != nil && existing.GetCID() != content.GetCID() { //NB. uppcase change is allowed
			result.Fields["login"] = "Username is used already"
		}
	}

	emailField := content.Value("email").(*fieldtype.Text)
	if emailField.String.String != email {
		if util.GetConfigSectionI("user")["unique_email"].(bool) {
			existing, _ := querier.Fetch("user", db.Cond("email", email))
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
