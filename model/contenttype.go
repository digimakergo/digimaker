//Author xc, Created on 2019-03-28 20:00
//{COPYRIGHTS}

package model

type ContentType struct {
	DataID    int
	Published int
	Modified  int
	Title     string
	RemoteID  string
}

//The helpers that generated entity will use to reduce template logic.

//Query:
//Content.List("article", "id > 0")
//Content.List("article", "title = 'Welcome'")
//
//Style 1:
//Query("article", and( "title = 'Welcome'", "body !=''" ) ).Sort( "id desc" )

//Style2:
//Query("article", `{ "condition": "title='Welcome' and body !=''",
//									  "sort":"id desc" 	}` )
//

//Style 3: (way similar to this is imporssible in go I think)
//Article.List( Cond( ( Article.ID > 5) && (Article.Modified > 5) ).Sort( Article.ID, desc ) )

//Style 4. Json Style
//Content.List( "article", `{id: 123, name: "test" }` )
//Content.List( "article", `{id: 123, author: [12,23] }` )
//Content.List( "article, folder", `{modify: [">", "123123130"}]` )
//Note: in this style, we could support MongoDB json's syntax style.
//
//Style5:
//Content.GetByID().Subtree( []Cond{ CondLocation("12"), CondModifiedLT( 123123130 ) ] } ).SortBy( "id" )
//
