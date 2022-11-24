GraphQL Query
====

`url:/api/graphql`

`header:{"api_key":"ddddxxxx7383423424sjfshfgfysifsik"}`

`method:POST`

Query Format
----
Format: 
```json
{
  "query": "{table(input){...ret}}",
  "operation":"content"
}
```

```
table:
    contenttype

input:
    filter:{field,gt(>),ge(>=),lt(<),le(<=),ne(!=)}
    sort:[\"id desc\"]
    limit:10
    offset:0
    
ret:
    output fields
```

Example:
```json
{
    "query": "{article(cid:464){id,name,title,published}}",
    "operation": "content"
}
```

```json
{
    "query": "{article(cid:[464,467]){id,name,title,published}}",
    "operation": "content"
}

```


```json
{
    "query": "{article(filter:{and:{body:\"ff\"}}){id,name,title,published,body}}",
    "operation": "content"
}
```


```json
{
    "query": "{article(filter:{le:{cid:471}}){id,name,title,published,body}}",
    "operation": "content"
}
```


```json
{
    "query": "{article(filter:{cid:[471,464]}){id,name,title,published,body},role{id,name}}",
    "operation": "content"
}
```

Some content type fields(defined by yourself)
---

| commonItem | type |
| :---:|:---:|
| id | string |
| version | string |
| published | string |
| modified | string |
| author | string |
| author_name | string |
| cuid | string |
| status | string |

---

| article | type |
| :---:|:---:|
| body | string |
| coverimage | string |
| editors | string |
| related_articles | string |
| related_articles | string |
| title | string |
| useful_resources | string |

---

| file | type |
| :---:|:---:|
| filetype | string |
| path | string |
| title | string |

---

| folder | type |
| :---:|:---:|
| display_type | string |
| summary | string |
| title | string |

---

| frontpage | type |
| :---:|:---:|
| mainarea | string |
| mainarea_blocks | string |
| sidearea | string |
| sidearea_blocks | string |
| slideshow | string |
| title | string |

---

| image | type |
| :---:|:---:|
| image | string |
| name | string |

---

| role | type |
| :---:|:---:|
| identifier | string |
| name | string |
| summary | string |
| under_folder | string |
