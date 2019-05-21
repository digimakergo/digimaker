//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

type RolePolicy struct {
	ID       int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	RoleID   string `boil:"name" json:"name" toml:"name" yaml:"name"`
	PolicyID string `boil:"description" json:"description" toml:"description" yaml:"description"`
	Subtree  string `boil:"subtree" json:"subtree" toml:"subtree" yaml:"subtree"`
	Scope    string `boil:"scope" json:"scope" toml:"scope" yaml:"scope"`
}
