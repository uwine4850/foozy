## Package builtin_mddl

__Auth__
```go
func Auth(excludePatterns []string, db *database.Database, onErr OnError) middlewares.MddlFunc
```
Auth is used to determine when to change the encoding of the AUTH cookie.
When keys are changed, the change date is set. If the date does not match, then you need to change the coding.
It is important to note that only previous keys are stored; accordingly, it is impossible to update the coding
if two or more key iterations have passed because the old keys are no longer known.
When this mddleware is enabled and the user is not logged in, a login error will be displayed. 
Therefore the __excludePatterns__ field is intended to allow the user to access some pages without having to log in.