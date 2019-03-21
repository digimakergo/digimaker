Application scenarios
==============
1. Website like CMS application, for enterprise mainly, which means for internet and intranet.
- hightlights: for subscriber visit(eg. ft or ftchinese) - dynamic content
               for mobile app which need login - dynamic content
               for distributed use(2 systems in 2 coutries). - maybe better than cdn


2. Application: for instance membership crm, trello, event snapchat for non-chating part 

3. Micro application: as a content engine it runs independent, somthing like solr(so non-database can be supported) 

Quality needs
==============
1) API should be minimal to use(write less code)

2) Database needs to be natual to application(not difficult database stucture that it's hard to extend in database level, like join).

3) Be clear what can be extended, what can not be


Future
===============

Performance
-------------
The system aims to have hundrends millsion/billion-intractive-level content, to achieve that, dm_location(which is the only table used by all other content) can be very big, so partition should be possible from start

Partition for dm_location, based on section first
