//Author xc, Created on 2019-04-23 22:20
//{COPYRIGHTS}

package fieldtype

type RelationField struct {
	Priority    int
	Description string
	Data        string
}

//restrcuture data to make it as real map instead of embeding data into.
func (r *RelationField) Restructure() {
	// fmt.Println(setting.RelationSettings["value_pattern"])
}
