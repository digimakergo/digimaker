package content

import (
	base "dm/models"
	orm "dm/models/orm"
)

type Article struct {
	*base.Content
	*orm.Article
}

func (article Article) Publish() {
	article.Content.Publish() //call parent
}
