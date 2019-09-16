Digimaker Content Management Framework
----------------
Digimaker is a simple, high performance and flexible Content Management Framework suitable to develop your web site and web application. Used scenarios are like websites, internal document management system, or generic software like workflow system, even CRM systems, or cloud based application.


### Simple

- Django-syntax like templating
- clear template structure
- easily support multi side
- Simple go language to extend api


### High performance
Written in Go language, with performance-prioritied practise, Digimaker CMF provides best performance among main stream languages. Benchmark ref:xxxx. Query data you need most.


### Powerful&Flexible
- rest api to query/change contents
- most common features are set in configuration, no coding needed.
- reuse built-in modules like login, displaying content
- extendable permission&user system
- powerful content model so extending&operating content is like operating database tables.
- clear & beautiful callback & debug mechanisms.


Doc
--------
See [doc](dm/doc) for ideas detail and evolving.

See [Progress](dm/doc/9.Progress.md) for progresses.


Progress
---------
### Phrase 1
1) Implement core api, including content type, Datatype, version, language.
 - Limit content types to: folder, article
 - Limit datatypes to: text, plaintext, datetime

2) Implement basic rest api for publishing, fetching

3) Implement A demo site for frontend

### Phrase 2
1) Extend 1) 2) above
2) Implement basic Admin UI.

In this stage it may be used in a small project. And then we involve project together with product.
