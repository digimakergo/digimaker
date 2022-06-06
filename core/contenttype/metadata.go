package contenttype

import (
	"time"

	"github.com/digimakergo/digimaker/core/definition"
)

type Metadata struct {
	Contenttype string    `boil:"_contenttype" json:"contenttype" toml:"contenttype" yaml:"contenttype"`
	Name        string    `boil:"_name" json:"name" toml:"name" yaml:"name"`
	Version     int       `boil:"_version" json:"version" toml:"version" yaml:"version"`
	Published   time.Time `boil:"_published" json:"published" toml:"published" yaml:"published"`
	Modified    time.Time `boil:"_modified" json:"modified" toml:"modified" yaml:"modified"`
	CUID        string    `boil:"_cuid" json:"cuid" toml:"cuid" yaml:"cuid"`
	Author      int       `boil:"_author" json:"author" toml:"author" yaml:"author"`
	AuthorName  string    `boil:"_author_name" json:"author_name" toml:"author_name" yaml:"author_name"`
	//Relations is used for binding all relationlist. See FinishBind to assign to different relationlist
	Relations ContentRelationList `boil:"_relations" json:"-" toml:"relations" yaml:"relations"`
}

//IdentifierList return list of all field names
func (c Metadata) IdentifierList() []string {
	return []string{"version", "name", "published", "modified", "author", "author_name", "cuid"}
}

//Values return values for insert/update DB. todo: rename to ToDBValues()
func (c Metadata) ToDBValues() map[string]interface{} {
	result := map[string]interface{}{
		"_version":   c.Version,
		"_name":      c.Name,
		"_published": c.Published,
		"_modified":  c.Modified,
		"_author":    c.Author,
		"_cuid":      c.CUID,
	}

	//for non-location content, delete undefined
	def, _ := definition.GetDefinition(c.Contenttype)
	if !def.HasLocation {
		for identifier, _ := range result {
			if !def.HasDataField(identifier) {
				delete(result, identifier)
			}
		}
	}

	return result
}

func (c Metadata) GetName() string {
	return c.Name
}

func (c Metadata) ContentType() string {
	return c.Contenttype
}

func (c Metadata) Definition(language ...string) definition.ContentType {
	def, _ := definition.GetDefinition(c.Contenttype, language...)
	return def
}

func (c *Metadata) GetRelations() ContentRelationList {
	return c.Relations
}
