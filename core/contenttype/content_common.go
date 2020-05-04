package contenttype

type ContentCommon struct {
	CID        int                 `boil:"cid" json:"cid" toml:"cid" yaml:"cid"`
	Version    int                 `boil:"version" json:"version" toml:"version" yaml:"version"`
	Published  int                 `boil:"published" json:"published" toml:"published" yaml:"published"`
	Modified   int                 `boil:"modified" json:"modified" toml:"modified" yaml:"modified"`
	CUID       string              `boil:"cuid" json:"cuid" toml:"cuid" yaml:"cuid"`
	Status     int                 `boil:"status" json:"status" toml:"status" yaml:"status"`
	Author     int                 `boil:"author" json:"author" toml:"author" yaml:"author"`
	AuthorName string              `boil:"author_name" json:"author_name" toml:"author_name" yaml:"author_name"`
	Relations  ContentRelationList `boil:"relations" json:"relations" toml:"relations" yaml:"relations"`
}

func (c ContentCommon) IdentifierList() []string {
	return []string{"cid", "version", "published", "modified", "author", "author_name", "cuid", "status"}
}

//Value for DB. todo: rename to DBValues()
func (c ContentCommon) Values() map[string]interface{} {
	result := map[string]interface{}{
		"id":        c.CID,
		"version":   c.Version,
		"published": c.Published,
		"modified":  c.Modified,
		"status":    c.Status,
		"author":    c.Author,
		"cuid":      c.CUID,
	}
	for identifier, relationValue := range c.Relations.Map {
		result[identifier] = relationValue
	}
	return result
}

func (c *ContentCommon) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {
	case "cid":
		result = c.CID
	case "version":
		result = c.Version
	case "modified":
		result = c.Modified
	case "published":
		result = c.Published
	case "author":
		result = c.Author
	case "author_name":
		result = c.AuthorName
	case "status":
		result = c.Status
	case "cuid":
		result = c.CUID
	case "relations":
		result = c.Relations
	}
	return result
}

func (c *ContentCommon) SetValue(identifier string, value interface{}) error {
	switch identifier {
	case "cid":
		c.CID = value.(int)
	case "version":
		c.Version = value.(int)
	case "published":
		c.Published = value.(int)
	case "modified":
		c.Modified = value.(int)
	case "status":
		c.Status = value.(int)
	case "author":
		c.Author = value.(int)
	case "cuid":
		c.CUID = value.(string)
	}
	return nil
}

func (c *ContentCommon) GetCID() int {
	return c.CID
}

func (c *ContentCommon) GetRelations() *ContentRelationList {
	return &c.Relations
}

//TODO: add more common methods related to content here.
