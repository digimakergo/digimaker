package dm/content

type Content struct{
  id int,
  content_id int,
  parent_id int,
  type string,
  name string,
  fields map[string]Field
}
