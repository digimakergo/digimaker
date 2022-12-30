package fieldtypes

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/digimakergo/digimaker/core/fieldtype"
)

type RelationHandler struct {
	fieldtype.FieldDef
}

//max 30 length
func (handler RelationHandler) LoadInput(ctx context.Context, input interface{}, mode string) (interface{}, error) {
	str := strings.TrimSpace(fmt.Sprint(input))
	var i int
	if str != "" {
		var err error
		i, err = strconv.Atoi(str)
		if err != nil {
			return nil, fieldtype.NewValidationError("Input is not a int. ref value:" + str)
		}
	}
	//todo: verify with parameters
	return i, nil
}

func (handler RelationHandler) DBField() string {
	return "INT NOT NULL DEFAULT 0"
}

func init() {
	fieldtype.Register(
		fieldtype.Definition{Name: "relation",
			DataType: "int",
			NewHandler: func(def fieldtype.FieldDef) fieldtype.Handler {
				return RelationHandler{FieldDef: def}
			}})
}
