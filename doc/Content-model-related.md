Content field:
-----------
Content name field should be 'title' as identifier(also fieldname in db).

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


Databases
----------
### dm_location
p: short name for partition to do query easiler. value 'c' means current partition, default will be c when creating
