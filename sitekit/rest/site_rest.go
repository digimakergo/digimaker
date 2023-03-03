package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/contenttype"
	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/query"
	"github.com/digimakergo/digimaker/core/util"
	"github.com/digimakergo/digimaker/rest"
	"github.com/digimakergo/digimaker/sitekit"
	"github.com/spf13/viper"
)

func GetView(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := util.CurrentUserID(ctx)
	if userID == 0 {
		//api key auth
		apiKey := viper.GetString("graphql.api_key") //todo: congfig api_key separate
		rApiKey := r.Header.Get("apiKey")
		if rApiKey != "" && apiKey == rApiKey {
			ctx = context.WithValue(r.Context(), util.CtxKeyUserID, 1)
		} else {
			rest.HandleError(errors.New("Invalid api key"), w)
			return
		}
	}

	idStr := r.URL.Query().Get("id")

	cType := r.URL.Query().Get("type")
	site := r.URL.Query().Get("site")
	viewmode := r.URL.Query().Get("viewmode")

	if cType == "" || site == "" || viewmode == "" {
		rest.HandleError(errors.New("Please input id, type, site and viewmode"), w)
		return
	}
	var contentList []contenttype.ContentTyper
	if idStr != "" {
		ids, _ := util.ArrayStrToInt(strings.Split(idStr, ","))
		contentList, _, _ = query.List(ctx, cType, db.Cond("c.id", ids))
	} else {
		parent, _ := strconv.Atoi(r.URL.Query().Get("parent"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit == 0 {
			limit = 10
		}
		if parent > 0 {
			contentList, _, _ = query.List(ctx, cType, db.Cond("l.parent_id", parent).Limit(0, limit))
		} else {
			rest.HandleError(errors.New("Invalid parameters"), w)
			return
		}
	}

	if contentList == nil {
		rest.HandleError(errors.New("No content found or no access"), w)
		return
	}
	result := map[int]string{}
	for _, content := range contentList {
		vars := map[string]interface{}{}
		vars["content"] = content
		vars["viewmode"] = viewmode //params.Mode
		output, err := sitekit.OutputString(vars, "content_view", sitekit.RequestInfo{Context: ctx, Site: site})
		if err != nil {
			rest.HandleError(err, w)
			return
		}
		result[content.GetID()] = output
	}
	rest.WriteResponse(result, w)
}

func init() {
	rest.RegisterRoute("/site/content/view", GetView)
}
