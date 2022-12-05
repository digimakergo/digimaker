GraphQL Query
====

`url:/api/graphql`

`header:{"api_key":"ddddxxxx7383423424sjfshfgfysifsik"}`

`method:POST`

Query Format
----
### Format: 
```graphql
{
  query: { <content type>: {filter:<filter>, sort:<sort>, limit:<limit>, offset:<offset>)
  {<fields>}}
  }
}
```

filter: support array(means `or`) or object(means `and`)




### Example
```json
query{
	article(filter:{author:1,title:"Test"})
	{
	 title,
	 summary,
	 body
	}
}
```

```json
query{
	article(filter:[{id:11}, {id:12},{title:"Test"}])
	{
	title,
	summary,
		body,
		fullbody
	},
	folder{
		title
	}
}
```
