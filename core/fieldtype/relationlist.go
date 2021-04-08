package fieldtype

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/db"
	"github.com/digimakergo/digimaker/core/definition"
	"github.com/digimakergo/digimaker/core/util"
)

type Relation struct {
	ID            int    `boil:"id" json:"-" toml:"id" yaml:"id"`
	ToContentID   int    `boil:"to_content_id" json:"-" toml:"to_content_id" yaml:"to_content_id"`
	ToType        string `boil:"to_type" json:"-" toml:"to_type" yaml:"to_type"`
	FromContentID int    `boil:"from_content_id" json:"from_content_id" toml:"from_content_id" yaml:"from_content_id"`
	FromType      string `boil:"from_type" json:"from_type" toml:"from_type" yaml:"from_type"`
	FromLocation  int    `boil:"from_location" json:"-" toml:"from_location" yaml:"from_location"`
	Priority      int    `boil:"priority" json:"priority" toml:"priority" yaml:"priority"`
	Identifier    string `boil:"identifier" json:"identifier" toml:"identifier" yaml:"identifier"`
	Description   string `boil:"description" json:"description" toml:"description" yaml:"description"`
	Data          string `boil:"data" json:"data" toml:"data" yaml:"data"`
	UID           string `boil:"uid" json:"-" toml:"uid" yaml:"uid"`
}

//A list of relation. Not used for bind so no Scan is not implemented
type RelationList []Relation

type RelationListHandler struct {
	definition.FieldDef
}

//LoadFromInput load data from input before validation
func (handler RelationListHandler) LoadInput(input interface{}, mode string) (interface{}, error) {
	s := fmt.Sprint(input)
	result := RelationList{}
	if s != "" {
		arr := strings.Split(s, ";")
		if len(arr) != 2 {
			return result, errors.New("wrong format")
		}
		arrInt, err := util.ArrayStrToInt(strings.Split(arr[0], ","))

		if err != nil {
			return result, errors.New("Not int(s)")
		}
		arrType := strings.Split(arr[1], ",")
		if len(arrType) != len(arrInt) {
			return result, errors.New("id and type are not same length")
		}

		for i, v := range arrInt {
			fromType := arrType[i]
			fromCid := v
			fromDef, err := definition.GetDefinition(fromType)
			if err != nil {
				return result, err
			}

			r := Relation{}

			r.FromContentID = fromCid
			r.FromType = fromType
			r.Priority = len(arrInt) - i
			r.Identifier = handler.Identifier

			relationDataFields := fromDef.RelationData
			if len(relationDataFields) > 0 {
				//get content
				contents := db.DatamapList{}
				db.BindEntity(context.Background(), &contents, fromDef.TableName, db.Cond("id", fromCid))
				if len(contents) == 0 {
					return result, errors.New("No content found on " + strconv.Itoa(fromCid))
				}

				//If there is one field, use it on data, otherwise use json map as data
				if len(relationDataFields) == 1 {
					r.Data = fmt.Sprint(contents[0][relationDataFields[0]])
				} else {
					datamap := map[string]interface{}{}
					for _, field := range relationDataFields {
						datamap[field] = contents[0][field]
					}
					data, _ := json.Marshal(datamap)
					r.Data = string(data)
				}
			}

			result = append(result, r)
		}
		//todo: update data after updating from_content
	}
	return result, nil
}

func (handler RelationListHandler) DBField() string {
	return ""
}

func init() {
	Register(
		Definition{
			Name:     "relationlist",
			DataType: "fieldtype.RelationList",
			Package:  "github.com/digimakergo/digimaker/core/fieldtype",
			NewHandler: func(fieldDef definition.FieldDef) Handler {
				return RelationListHandler{FieldDef: fieldDef}
			}})
}
