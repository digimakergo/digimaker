Basic way to do administration
===============

### Model change(content type change)
Unlike some other cms, this framework needs technical people to do model change, so everything which is for advanced editor's part should be implemented in editor part, eg. options for select datatype in a content type.

This part can be hidden from editor's view completely, or create a special role for maintenancer. In that case administrator is the editor administrator, who doesn't have access to model change either.

Content type model change will
- trigger generating orm enities.
- create a sql for administrator to execute(or execute online), but will not execute it automatically.
- need to regenerate orm and maybe others after execution.
- never delete column or table automatically.
