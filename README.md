Digimaker Content Management Framework
----------------
Written in Go language, Digimaker is a simple, high performance and flexible Content Management(but not limited to) Framework to develop your web site and web application.

The framework emphasizes below design principles:
- minimal design, avoiding over-engineering.
- balance of simplicity, performance and flexibility
- easy to maintain(eg. good logging, debug info, error message)

Typical application cases are like website, document management system, cloud based application, also some generic software like workflow systems, even small CRM systems or similar.

### Simple
- [website]simple templating syntax
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

### Flexible
- rest api to query/change contents
- extendable permission&user system
- powerful content model(content type&field type) to extending&operating content fitting your need
- clear callback & debug mechanisms.
- online debug so administrator can monitor request's processing data, time spent, etc.


Documentation
--------
See [doc](core/doc) for ideas detail and evolving.

See [Progress](core/doc/9.Progress.md) for progresses.

License
--------
MIT license. 

Support & Services
--------
Almost all activites can be done via community. 

For specific support/service you can contact Digimaker Go AS. Also there are paid plugins providing specific features(eg. maintenance tools). 

Check more on [www.digimaker.com](http://www.digimaker.com).
