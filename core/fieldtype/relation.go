package fieldtype

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/digimakergo/digimaker/core/db"
)

type RelationParameters struct {
	Type      string                 `json:"type"`
	Value     string                 `json:"value"`
	Condition map[string]interface{} `json:"condition"`
}

type Relation struct {
	String
}

func (r Relation) Type() string {
	return "relation"
}

//Scan scan data from db.
func (r *Relation) Scan(src interface{}) error {
	switch src.(type) {
	case int64:
		r.String.String = strconv.Itoa(int(src.(int64)))
	default:
		return r.String.Scan(src)
	}
	return nil
}

func (r *Relation) LoadFromInput(input interface{}, params FieldParameters) error {
	err := r.String.LoadFromInput(input, params)
	if err != nil {
		return err
	}

	if r.String.String == "" {
		return nil
	}

	rParams, err := ConvertRelationParams(params)
	if err != nil {
		return err
	}

	value := r.String.String
	if rParams.Type == "" {
		return errors.New("Need type definition.")
	}

	valueColumn := "id"
	if rParams.Value != "" {
		valueColumn = rParams.Value
	}

	dbHandler := db.DBHanlder()
	condition := db.Cond(valueColumn, value)
	for cKey, cValue := range rParams.Condition {
		condition = condition.And(cKey, cValue)
	}

	//todo: support Type to be a content type, not table name - also for query.
	// to do so, we need to invoke content def from field - infiniate package look issue
	count, err := dbHandler.Count(rParams.Type, condition)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Not found in " + rParams.Type)
	}

	return nil
}

//Convert to parameters obj
func ConvertRelationParams(params FieldParameters) (RelationParameters, error) {
	paramsData, _ := json.Marshal(params)
	rParams := RelationParameters{}
	err := json.Unmarshal(paramsData, &rParams)
	if err != nil {
		return rParams, errors.New("Wrong definition of parameters:" + err.Error())
	}
	return rParams, nil
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "relation", Value: "fieldtype.Relation"},
		func() FieldTyper { return &Relation{} })
}
