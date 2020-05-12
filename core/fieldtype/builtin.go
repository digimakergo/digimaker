package fieldtype

import (
	"database/sql/driver"
	"os"
	"strings"

	"github.com/xc/digimaker/core/log"
	"github.com/xc/digimaker/core/util"
)

/**
* Initial types END
 */
//Text
type Text struct {
	String
}

func (t Text) Type() string {
	return "text"
}

type RichText struct {
	String
}

func (rt RichText) Type() string {
	return "richtext"
}

type Number struct {
	Int
}

func (n Number) Type() string {
	return "number"
}

//Radio struct represent radio type
type Radio struct {
	Int
}

func (r Radio) Type() string {
	return "radio"
}

//Radio struct represent radio type
type Checkbox struct {
	Int
}

func (c Checkbox) Type() string {
	return "checkbox"
}

//Password struct represent password type
type Password struct {
	String
}

func (r Password) Type() string {
	return "password"
}

//Password struct represent password type
type Image struct {
	String
	changed bool
}

func (i Image) Type() string {
	return "image"
}

func (i *Image) LoadFromInput(input interface{}) error {
	original := i.String.String
	str := i.String
	err := str.LoadFromInput(input)
	if err == nil {
		i.String = str
		if i.String.String != "" && original != i.String.String {
			i.changed = true
		}
	}
	return err
}

func (i Image) Value() (driver.Value, error) {
	filepath := i.String.String
	if i.changed && filepath != "" {
		oldAbsPath := util.VarFolder() + "/" + filepath
		arr := strings.Split(filepath, "/")
		filename := arr[len(arr)-1]
		newPath := "uploaded/" + filename
		newAbsPath := util.VarFolder() + "/" + newPath
		//todo: create subfolder for it to save performance.
		//todo: create thumbnail
		err := os.Rename(oldAbsPath, newAbsPath)
		if err != nil {
			log.Error("Can not move "+filepath+". error: "+err.Error(), "")
		}
		//even if it failed, it will use new path(even if it failed - might because of file already been moved)
		filepath = newPath
	}
	return filepath, nil
}

type RelationList struct {
	JSON
}

func (r RelationList) Type() string {
	return "relationlist"
}

func init() {
	RegisterFieldType(
		FieldtypeDef{Type: "text", Value: "fieldtype.Text"},
		func() FieldTyper { return &Text{} })
	RegisterFieldType(
		FieldtypeDef{Type: "json", Value: "fieldtype.JSON"},
		func() FieldTyper { return &JSON{} })
	RegisterFieldType(
		FieldtypeDef{Type: "radio", Value: "fieldtype.Radio"},
		func() FieldTyper { return &Radio{} })
	RegisterFieldType(
		FieldtypeDef{Type: "checkbox", Value: "fieldtype.Checkbox"},
		func() FieldTyper { return &Checkbox{} })
	RegisterFieldType(
		FieldtypeDef{Type: "richtext", Value: "fieldtype.RichText"},
		func() FieldTyper { return &RichText{} })
	RegisterFieldType(
		FieldtypeDef{Type: "number", Value: "fieldtype.Number"},
		func() FieldTyper { return &Number{} })
	RegisterFieldType(
		FieldtypeDef{Type: "password", Value: "fieldtype.Password"},
		func() FieldTyper { return &Password{} })
	RegisterFieldType(
		FieldtypeDef{Type: "image", Value: "fieldtype.Image"},
		func() FieldTyper { return &Image{} })
	RegisterFieldType(
		FieldtypeDef{Type: "relationlist", Value: "fieldtype.RelationList", IsRelation: true},
		func() FieldTyper { return &RelationList{} })
}
