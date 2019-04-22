package entity

type ContentCommon struct {
	CID       int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Published int    `boil:"published" json:"published" toml:"published" yaml:"published"`
	Modified  int    `boil:"modified" json:"modified" toml:"modified" yaml:"modified"`
	RemoteID  string `boil:"remote_id" json:"remote_id" toml:"remote_id" yaml:"remote_id"`
}

func (c ContentCommon) IdentifierList() []string {
	return []string{"cid", "published", "modified", "remote_id"}
}

func (c ContentCommon) Values() map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = c.CID
	result["published"] = c.Published
	result["modified"] = c.Modified
	result["remote_id"] = c.RemoteID
	return result
}

func (c *ContentCommon) Value(identifier string) interface{} {
	var result interface{}
	switch identifier {
	case "cid":
		result = c.CID
	case "modified":
		result = c.Modified
	case "published":
		result = c.Published
	case "remote_id":
		result = c.RemoteID
	}
	return result
}

func (c *ContentCommon) SetValue(identifier string, value interface{}) error {
	switch identifier {
	case "id":
		c.CID = value.(int)
	case "published":
		c.Published = value.(int)
	case "modified":
		c.Modified = value.(int)
	case "remote_id":
		c.RemoteID = value.(string)
	}
	return nil
}

func GetCID(c *ContentCommon) int {
	return c.CID
}

//TODO: add more common methods related to content here.
