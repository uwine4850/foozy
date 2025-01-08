## namelib
Сommon names for data access.

Further, the designation UC will indicate that this name will be used in the UserContext. More specifically:
* `UC` — data retrieval only.
* `UC<data_type>` — data installation. They can also be retrieved in the same way.
* `CK` — used in cookies.
* `RN` — retrieving data from the templating engine.
___
### AuthNames
The name for the embedded auth package.

*AUTH_TABLE* — the name of the authentication table in the database.\
`CK` *COOKIE_AUTH* — cookie authentication name.\
`CK` *COOKIE_AUTH_DATE* — access to the cookie authentication time.
___
### RouterNames
`UC` *URL_PATTERN* — current url pattern.\
`CK` *COOKIE_CSRF_TOKEN* — CSRF token in the cookie.\
`UC<string>` *MDDL_ERROR* — middleware error.\
`UC<bool>` *SKIP_NEXT_PAGE* — skipping page rendering.\
`RN` *REDIRECT_ERROR* — error with redirection to another page. Set in URl, once set the data can be retrieved in the templating engine.\
`UC<string>` *SERVER_ERROR* — error with redirection to another page.\
`UC<string>` *SERVER_FORBIDDEN_ERROR* — error with redirection to another page. Access to UserContext.
___
### ObjectNames
`UC` *OBJECT_CONTEXT* — access to the context of the object set in the view.\
`UC` *OBJECT_CONTEXT_FORM* — access to the form context in the view. Used only in *FormView*.\
`UC` *OBJECT_DB* — access to the database that is opened in the view.