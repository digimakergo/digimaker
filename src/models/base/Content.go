package base

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

import (
	"models/orm"
)

type Content struct {
	*orm.Location
	fields map[string]Field //can we remove the fields and article.title directly?
}
