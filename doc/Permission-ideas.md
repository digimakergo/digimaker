Permisson system
----------------
### Support attribute level permission
Permission should support attribute level(at least for simple value) for condition as operation. This will make it very powerful when extending the permission system. Traditionally section, language, subtree are typical conditions, with attribute support we can
- easily support language condition(attribute language)
- set something like "users can access to articles under new in latest 30 days"
- set something like "users can edit article title but not body or upload new file but not change title". - this will be useful in editorial workflow system.

The operation in attribute might be less important for content editing, but very important for developing project, typically when a customized table needs to be controller to field level. So the permission should either provide extending way(eg. customziable json & ui callback) so it's possible to extend this.

### Change log
Permission change should be logged. It's not necessary to use EndDate for permission so all the history is kept, but it will be useful to track user's permission. Like below.
```
id, log_type, user, target_id, detail
22, role_assigment,122, 121<here is user group id>, 'user group Dev department is assigned to role 21321'
```

It is used for non-technial administrator, so "what to log" should consider that. Technical log is there already.
