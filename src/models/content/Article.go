package content

import (
	orm "models/orm"
)

type Article struct {
	*Content
	*orm.Article
}
