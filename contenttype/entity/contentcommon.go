package entity

type ContentCommon struct {
	CID       int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Published int    `boil:"published" json:"published" toml:"published" yaml:"published"`
	Modified  int    `boil:"modified" json:"modified" toml:"modified" yaml:"modified"`
	RemoteID  string `boil:"remote_id" json:"remote_id" toml:"remote_id" yaml:"remote_id"`
}

func (c ContentCommon) Values() map[string]interface{} {
	result := make(map[string]interface{})
	result["id"] = c.CID
	result["published"] = c.Published
	result["modified"] = c.Modified
	result["remote_id"] = c.RemoteID
	return result
}

//TODO: add more common methods related to content here.
