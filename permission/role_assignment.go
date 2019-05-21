//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

type RoleAssignment struct {
	ID         int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	AssignTo   string `boil:"assign_to" json:"assign_to" toml:"assign_to" yaml:"assign_to"`
	AssignType string `boil:"assign_type" json:"assign_type" toml:"assign_type" yaml:"assign_type"`
}
