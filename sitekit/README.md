DM Sitekit
==========
Sitekit is a toolkit which is used to build sites in a 'template way'.

### Core features

- easy templating
- multi sites support
- load multi sites from yaml config file or api, or both.
- powerful template override based on content conditions
- nice url(from niceurl package) and extendable
- extend site router with templating
- customize template functions, filters, macro

### Template engine
The template engine used is pongo2 https://github.com/flosch/pongo2.

### Demosite
See [Demosite](../demosite) for example use.

### Template Functions
Examples:

Fetch content by id:

    {%set content = dm.fetch_byid(8)%}

Fetch children:

    {%set children = dm.children( content, "article" )%}


#### Filters

Output variable:

    {{var|dmshow}}

Format time to local:

    {{timestamp|dm_format_time}}


### Template Override
