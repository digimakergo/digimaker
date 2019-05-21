//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

type Role struct {
	ID          int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name        string `boil:"name" json:"name" toml:"name" yaml:"name"`
	Description string `boil:"description" json:"description" toml:"description" yaml:"description"`
}
