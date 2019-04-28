package handler

import (
	"dm/contenttype"
	"dm/contenttype/entity"
	"dm/db"
	"dm/query"
	"dm/util"
	"strconv"

	"github.com/pkg/errors"
)

type ContentQuery struct{}

//Fetch content by location id.
//If no location found. it will return nil and error message.
func (cq ContentQuery) FetchByID(locationID int) (contenttype.ContentTyper, error) {
	//get type first by location.
	dbhandler := db.DBHanlder()
	location := entity.Location{}
	err := dbhandler.GetEnity("dm_location", query.Cond("id", locationID), &location)
	if err != nil {
		return nil, errors.Wrap(err, "[contentquery.fetchbyid]Can not fetch location by locationID "+strconv.Itoa(locationID))
	}
	if location.ID == 0 {
		return nil, errors.New("[contentquery.fetchbyid]Location is empty.")
	}

	//fetch by content id.
	contentID := location.ContentID
	contentType := location.ContentType
	result, err := cq.FetchByContentID(contentType, contentID)
	return result, err
}

//Fetch a content by content id.
func (cq ContentQuery) FetchByContentID(contentType string, contentID int) (contenttype.ContentTyper, error) {
	return cq.Fetch(contentType, query.Cond("content.id", contentID))
}

//Fetch one content
func (cq ContentQuery) Fetch(contentType string, condition query.Condition) (contenttype.ContentTyper, error) {
	//todo: use limit in this case so it doesn't fetch more into memory.
	content := entity.NewInstance(contentType)
	err := cq.Fill(contentType, condition, content)
	if err != nil {
		return nil, err
	}
	return content, err
}

//Fetch a list of content, return eg. *[]Article
func (cq ContentQuery) List(contentType string, condition query.Condition) (interface{}, error) {
	contentList := entity.NewList(contentType)
	err := cq.Fill(contentType, condition, contentList)
	if err != nil {
		return nil, err
	}
	return contentList, err
}

//Fill all data into content which is a pointer
func (cq ContentQuery) Fill(contentType string, condition query.Condition, content interface{}) error {
	dbhandler := db.DBHanlder()
	err := dbhandler.GetByFields(contentType, condition, content)
	if err != nil {
		message := "[List]Content Query error"
		util.Error(message, err.Error())
		return errors.Wrap(err, message)
	}
	return nil
}

//todo: use method instead of global variable
var querier ContentQuery = ContentQuery{}

func Querier() ContentQuery {
	return querier
}
