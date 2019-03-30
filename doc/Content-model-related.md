Content field:
-----------
Content name field should be 'title' as identifier(also fieldname in db).


ContentType generating
----------
Code generating is a good way, but we need to validate the performance and auto deployment. In some cases it will need a lightweight way to do everything from web(for instance, design form online which reused field types, the form design is done/modified by editor sometimes). So.

```
type Article struct{
  Title TextField
  Body RichTextField
}
```
can be a generic(so you just generate json for the definition, not .go)
```
type ContentType struct{
  Fields map[string]Fielder
}
```



Queries
---------
Typical content select query condition include:
  - common attribute: eg. published, modified, parent_location
  - field: eg. by user/job_title = 'engineer',
  - types: 'frontpage, frontpage_sub', with common field: have_children = 0, {"/have_children": 0, "types": "frontpage, frontpage_sub"}. This is more like a union.

Sorting should be by type first, it's not needed to mix types eg. by published.


In terms of query, it's important to have right content model. There are 2 types of model for query(take folder, frontpage, frontpage_sub as example ):

**Model1**:
```
Table attribute_data:
identifier, type, value
----------------------
'title', 'folder', 'Home'
'title', 'frontpage', 'Front page'
'left_menu', 'frontpage', '223'
'title', 'frontpage_sub', 'Oslo club front page'
'frontpage_sub', 'club_logo', '2255'
```

This helps to query multiple type because you can

```
SELECT * FROM attribute_data WHERE type IN ('folder', 'fontpage', 'fontpage_sub') AND ...
```

**Model2**
```
Table folder:
title, summary
--------------
'Home', ''

Table frontpage:
left_menu
---------
223


Table frontpage_sub:
club_logo
---------
2255

```
Model2 following the rational database and normal data principle, but it will have big problem when it has many type query at the same time(you have to use many unions and we should minimize sort and limit after union).

A complex site(which can have 100 sub sites even) can have types like:
folder
frontpage
frontpage_club
infobox_container
infobox
infobox_club
campaign

Idea: it would be good to have a sub_type concept, which is a special type of common type.
1) folder can have folder_type: image, organization, building - they will have different icon, template rewrite rule(rewrite rule will support attribute in general.).

2) Can 2 content type be in one datatable, eg. frontpage_club, frontpage, it's useful when there are not too many columns together(with a type frontpage_type: 'club')? (it can be both good for selecting, but bad for name conversion since one content type doesn't mean direct table match(use 'frontpage/club in this case?'))

Language
---------
Should language have it's own location?
**Have it's location, and with a relation(can call translation) between language articles**
- Language article is not different from normal article, so you can do copy, move, set permission .
- Query needs language

**Not in location, but part of an article**
- Always need to provide language when query(otherwise default language), visit(the system can use site language as a solution).
- You need special development to view an article's translation(if it's a normal article with different id you don't need special treatment).
- When set permission, need to select language. It makes the development alway have to think about language because it's part of the key of an article(in this case: key( article_id, language )).
- Good thing is that you can show main language automatically if it doesn't have a translation.

Conclusion: language can be an attribute to the object(article), but it shouldn't need language to identify an article, meaning it can be in location table. Translation is a relation between locations.

It's not needed to create a translation site where its content is from main site. If it's needed, we can provide a feature to "copy sub content with translation" which do the copy and build translation relation.

Trash
-------
We can reuse p field in location instead of create a table, since we alway have p when we do query. In this context, trash is a logical partition - similar to archive. We can use 'trash' in p. When me move a folder to trash, we mark all the subtree's p as trash, it will be quick.

Idea: can we have similar/same feature between trash and other partitions? Restore is actually same as "moving to c partition". In that case, you can actually see the structure in the trash - which might be helpful. eg.

```
- News
  - News1
    - News1 Test1
    - News1 Test2(removed)
  - News2
    - News 22(removed)
      - News 22 Test1(removed)
      - News 22 Test2(removed)
- Trash
 - Root<where parent is not in trash>
    - News1 Test2(shown in list)
 - Neews22
    - News 22 Test1(shown in list)
    - News 22 Test2(shown in list)
```


Databases
----------
### dm_location
p: short name for partition to do query easiler. value 'c' means current partition, default will be c when creating
