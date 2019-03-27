Application scenarios
==============
1. Website like CMS application, for enterprise mainly, which means for internet and intranet.
- hightlights: for subscriber visit(eg. ft or ftchinese) - dynamic content
               for mobile app which need login - dynamic content
               for distributed use(2 systems in 2 coutries). - maybe better than cdn


2. Application: for instance membership crm, trello, event snapchat for non-chating part

3. Micro application: as a content engine it runs independent, somthing like solr(so non-database can be supported)

Architecture
---------------
Some principles:
1. DO NOT try to do something which looks nice, but against the principle for long term. Some nice-try feature is actually blocking some good principles(but can try for short term). For instance,
  - try to change table structure online -  this is possible, but what if it failed, or timeout - we don't have rollback mechanism/backup online. Database change is better done by database console directly - part of the DBA role and there is a routine for backup. We can generate sql so the executor can preview what's inside and do related action.
  - try to do distributed system with too much data transferring by yourself. Distributed system need to be thought from whole projects perspective with more thoughts(with tools) for load balance, table design, database partition, etc. However, this framework can try to have content partition based on time for example, which can be used by database partition. Distributed system based on service is possible. But if we use load balance, all the services under load balance should be the same.
```
  - service 1
      - load balancer
         - service 1 implementation with same data(if data is too much, create partition data or move partition data with a new instance service3)
         - service 1 implementation with same data
         ...         
  - service 2
```
NOT THIS(this is why mysql master-slave or other sync is not a good solution because data1 and data2 below sync can be desaster, specally when traffic is high, you have no control of data sync). The principle for sync is that: service should not be available while it's data is syncing.
```
  - load balancer
     - service 1 implementation with data1
     - service 1 implementation with data1
     - service 1 implementation with data2 which sync with data1
     ...         
     - service 2
     ...
```


Content model
=================
1) [relate to performance]Should be able to move horizontal data to different partition or database.
And the system should be able to load whole bunch of content data in another view(maybe a tab on top or tab under the node). If it's moved to a different database, it's possible to instance a new system with the part data.

 This needs everything can be horizontally chunked - content, content relation, images should be inside 1 chunk,

2) [related to performance and migration]the images shouldn't be in a separated folder - they should be under the folder it belongs. But for ui, we can have a separate tab(library) for all images with structure. eg.

```
news
  - images(a image container type or a folder with folder type image)
  - files
  - <domestic news>
  - <tech news>
```

So the library can be like this:
```
library
  - images
     - news
     - <other virtual folder>
  - files
     -news
     - <other virtual folder>
```

NB: putting images&files into content structure can fundamentally make the parent into a partition, move to another installation, not global. It's impossible to separate it without doing this, especially if you even didn't create a separate folder in media for use images. It's like folder, if you have resource somewhere else, you have to copy serveral times. The media library is always a symlink in this pespective:).


3) [related to separation] Seldom the content images are used globally many times(maybe many times, but in similar location). The resource image(eg.logo) maybe used in many place. So
 - If you know it's not important to have update when image is updated, copy image to near images folder, instead of always using reference.
 - When migrate content partition, copy the shared content(typically share image which are outside of the partition)

4) [image]Image should have options to not version it - versioning take too much space.

5) [images]image&file need authentication&permission check.

6) [images]images can be done completely using a cdn image api(with permission check locally).


API
========
1. Rest api
------------
There should be a flexible&powerful query api that you can query once and get what you need.
eg.

Simple ones:
```
  /content/<id> you get content
  /content/list/<location_id> you get list of contents
```
complex one(like union): get name, id, created from article and files in recent one week
 ```
 /query/select/name,id,created/from/article,file/where/created>10
 ```

Should we support a query language(like Doctrine's DQL) or json like( { "select": "" } ). It all depends on application. Normally if it's not difficult query, url should be good enough.

2. Local api
------------
Should avoid to use sql directly(and we will support NOSQL), use query api(eg. where( "created", ">=32131321" )).


Quality needs
==============
1) API should be minimal to use(write less code)

2) Database needs to be natual to application(not difficult database stucture that it's hard to extend in database level, like join).

3) Be clear what can be extended, what can not be


Future
===============

Performance
-------------
The system aims to have hundrends million/billion-intractive-level content, to achieve that, dm_location(which is the only table used by all other content) can be very big, so partition should be possible from start

Partition for dm_location, based on section first
