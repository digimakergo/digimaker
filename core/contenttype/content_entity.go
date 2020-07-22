package contenttype

type ContentEntity struct {
	ID          int                 `boil:"id" json:"id" toml:"id" yaml:"id"`
	ContentType string              `boil:"content_type" json:"content_type" toml:"content_type" yaml:"content_type"`
	CID         int                 `boil:"cid" json:"cid" toml:"cid" yaml:"cid"`
	Relations   ContentRelationList `boil:"relations" json:"relations" toml:"relations" yaml:"relations"`
}

func (c *ContentEntity) GetRelations() *ContentRelationList {
	return &c.Relations
}

func (c *ContentEntity) GetCID() int {
	return c.ID
}

//TODO: add more common methods related to content here.
