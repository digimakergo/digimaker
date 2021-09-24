package contenttype

type ContentEntity struct {
	ID          int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	ContentType string `boil:"content_type" json:"content_type" toml:"content_type" yaml:"content_type"`
	CID         int    `boil:"cid" json:"cid" toml:"cid" yaml:"cid"`
	//todo: may put more here like Author, LocationID
	Author     int                 `boil:"author" json:"author" toml:"author" yaml:"author"`
	AuthorName string              `boil:"author_name" json:"author_name" toml:"author_name" yaml:"author_name"`
	Relations  ContentRelationList `boil:"relations" json:"relations" toml:"relations" yaml:"relations"`
}

func (c *ContentEntity) SetValue(identifier string, value interface{}) error {
	switch identifier {
	case "author":
		c.Author = value.(int)
	}
	return nil
}

func (c *ContentEntity) GetRelations() ContentRelationList {
	return c.Relations
}

func (c *ContentEntity) GetCID() int {
	return c.ID
}

func (c *ContentEntity) GetID() int {
	return c.ID
}

func (c *ContentEntity) GetAuthor() int {
	return c.Author
}

//TODO: add more common methods related to content here.
