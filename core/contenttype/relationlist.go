package contenttype

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
