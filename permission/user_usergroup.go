//Author xc, Created on 2019-05-22 22:10
//{COPYRIGHTS}

package permission

type RoleAssignment struct {
	ID          int    `boil:"id" json:"id" toml:"id" yaml:"id"`
	UserID      int    `boil:"user_id" json:"user_id" toml:"user_id" yaml:"user_id"`
	UsergroupID string `boil:"usergroup_id" json:"usergroup_id" toml:"usergroup_id" yaml:"usergroup_id"`
}
