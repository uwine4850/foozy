## Package builtin_mddl

__Auth__
```
Auth(loginUrl string, db *database.Database) middlewares.MddlFunc
```
Auth is used to determine when to change the encoding of the AUTH cookie.
When keys are changed, the change date is set. If the date does not match, then you need to change the coding.
It is important to note that only previous keys are stored; accordingly, it is impossible to update the coding
if two or more key iterations have passed because the old keys are no longer known.