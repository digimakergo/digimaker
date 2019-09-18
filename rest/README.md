REST package
==============
Rest package handles all the rest request, eg. query content, add/edit/delete content. Use subroute to run it.


List of apis:
--------------------
Note: the real url can depend on routering. The default one is /api. eg. /api/content/get/2
#### Fetch content
content/get/{id}

content/children/{id}

content/list/{id}

#### Operate content



i18n library needed
---------
Besides standand transation, the i18n library can be used for customized message also. use key or message?
The library should be
- easy to use. eg. one function should be enough
- easy to edit and can be changed online. eg. json/toml
- better support count
- support override in project so many standard message(eg. validation, error) can be overridden based on special need(not only translation, but also english). Customzied template message(eg. error template) should be done in template override, not here.
- message should be written inline to improve readability.
- message supports non-english in code so you can actually write Chinese in source and translate to english.
Example:
 i18n.T( "File format is not supported" ) //message with global context
 i18n.T( "File format is not supported", "rest.error" ) // message+context
 note: the message can be treated as a key also. so you can actually customize message
 in ch.json:
{
 "context": "rest.error",
 "text": "File format is not supported",
 "translation": "不支持的文件类型."
}
also in eng.json in project{
    {
     "context": "rest.error",
     "text": "File format is not supported",
     "translation": "File is not supported, check (http://xxxf)[this] link?."
    }
}

so translation file can be a kind of template override.
text/locale/eng.json
text/locale/chi.json
text/customize/eng.json
text/customize/chi.json
