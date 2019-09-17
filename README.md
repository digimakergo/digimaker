Digimaker Content Management Framework
----------------
Written in Go language, Digimaker is a simple, high performance and flexible Content Management (but not limited to) Framework to develop your web site and web application. Typical application cases are like websites, internal document management system, cloud based application, also some generic software like workflow systems, even CRM systems.

The framework emphasizes below design principles:
- minimal core+plugin design
- avoiding over-engineering, features fitting scenarios. Get things done with less, clean, and beautiful code
- make things easy to maintain(eg. good logging, debug info, error message)

### Simple
- [websites]Django/Twig-syntax like templating
- [websites]clear template structure
- [websites]easily support multi side
- Go style api
- Manipulate content via rest api.


### High performance
Thanks to performance of Go language, with performance-prioritied design, Digimaker CMF provides very good performance.
- Support cluster
- Straightforward content model to database, query data you need directly.
- Support database partition from model so querying > 10 millions of data can be like querying under 1 million.
Benchmark reference:xxxx.

### Powerful&Flexible
- rest api to query/change contents
- most common features are set in configuration, no coding needed.
- reuse built-in modules like login, displaying content
- extendable permission&user system
- powerful content model(content type&field type) to extending&operating content fitting your need
- clear & beautiful callback & debug mechanisms.

License
--------
Digimaker is honored to use MIT license(confirmed?). There is Paid Digimaker Plugins which provide additional valuable features(eg. maintenance tools), check more [here](http://www.digimaker.com).

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
