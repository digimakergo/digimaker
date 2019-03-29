Permisson system
----------------


###Change log
Permission change should be logged. It's not necessary to use EndDate for permission so all the history is kept, but it will be useful to track user's permission. Like below.

id, log_type, user, target_id, detail 
22, role_assigment,122, 121<here is user group id>, 'user group Dev department is assigned to role 21321'


It is used for non-technial administrator, so "what to log" should consider that. Technical log is there already.


