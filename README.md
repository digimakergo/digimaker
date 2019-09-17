Digimaker Content Management Framework
----------------
Written in Go language, Digimaker is a simple, high performance and flexible Content Management(but not limited to) Framework to develop your web site and web application.

The framework emphasizes below design principles:
- minimal core+plugin design, avoiding over-engineering.
- providing features with good balance of simplicity and flexibility
- important to make everything easy to maintain(eg. good logging, debug info, error message)

Typical application cases are like website, document management system, cloud based application, also some generic software like workflow systems, even CRM systems.

### Simple
- [website]Django/Twig-syntax like templating
- [website]clear template structure
- [website]easily support multi side
- Go style api
- Manipulate content via rest api.


### High performance
Thanks to performance of Go language, with performance-prioritied design, Digimaker CMF provides very good performance without using cache server.
- Straightforward content model to database, query data you need directly.
- Support cluster
- Support database partition from model so querying > 10 millions of data can have same performance as querying under 1 million.
Benchmark reference:xxxx.

### Powerful&Flexible
- rest api to query/change contents
- most common features are set in configuration, no coding needed.
- reuse built-in modules like login, displaying content
- extendable permission&user system
- powerful content model(content type&field type) to extending&operating content fitting your need
- clear & beautiful callback & debug mechanisms.


Documentation
--------
See [doc](dm/doc) for ideas detail and evolving.

See [Progress](dm/doc/9.Progress.md) for progresses.

License
--------
Digimaker is honored to use MIT license(confirmed?). There are Paid Digimaker Plugins providing additional valuable features(eg. maintenance tools), check more on [http://www.digimaker.com](www.digimaker.com).
