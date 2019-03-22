package dm/content

/**
This is a parent struct which consits of location and the content itself(eg. article).
*/

type Content struct{
  id int,
  content_id int,
  parent_id int,
  content_type string,
  name string,
  fields map[string]Field //can we remove the fields and article.title directly?
}
