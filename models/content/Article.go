package content

import (
	base "models"
	orm "models/orm"
)

type Article struct {
	*base.Content
	*orm.Article
}

func (article Article) Publish() {
	article.Content.Publish() //call parent
}
