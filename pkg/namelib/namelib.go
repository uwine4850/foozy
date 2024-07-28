package namelib

type AuthNames struct {
	AUTH_TABLE       string
	COOKIE_AUTH      string
	COOKIE_AUTH_DATE string
}

var AUTH = AuthNames{
	AUTH_TABLE:       "auth",
	COOKIE_AUTH:      "AUTH",
	COOKIE_AUTH_DATE: "AUTH_DATE",
}

// AUTH_TABLE string = "auth"
// COOKIE_AUTH string = "AUTH"
// COOKIE_AUTH_DATE string = "AUTH_DATE"

type RouterNames struct {
	URL_PATTERN            string
	COOKIE_CSRF_TOKEN      string
	MDDL_ERROR             string
	SKIP_NEXT_PAGE         string
	REDIRECT_ERROR         string
	SERVER_ERROR           string
	SERVER_FORBIDDEN_ERROR string
}

var ROUTER = RouterNames{
	URL_PATTERN:            "URL_PATTERN",
	COOKIE_CSRF_TOKEN:      "CSRF_TOKEN",
	MDDL_ERROR:             "MDDL_ERROR",
	SKIP_NEXT_PAGE:         "SKIP_NEXT_PAGE",
	REDIRECT_ERROR:         "REDIRECT_ERROR",
	SERVER_ERROR:           "SERVER_ERROR",
	SERVER_FORBIDDEN_ERROR: "SERVER_ERROR",
}

// // Router
// const URL_PATTERN string = "URL_PATTERN"
// const COOKIE_CSRF_TOKEN string = "CSRF_TOKEN"
// const MDDL_ERROR string = "MDDL_ERROR"
// const SKIP_NEXT_PAGE string = "SKIP_NEXT_PAGE"
// const REDIRECT_ERROR string = "REDIRECT_ERROR"
// const SERVER_ERROR string = "SERVER_ERROR"
// const SERVER_FORBIDDEN_ERROR string = "SERVER_ERROR"

type ObjectNames struct {
	OBJECT_CONTEXT      string
	OBJECT_CONTEXT_FORM string
	OBJECT_DB           string
}

var OBJECT = ObjectNames{
	OBJECT_CONTEXT:      "OBJECT_CONTEXT",
	OBJECT_CONTEXT_FORM: "OBJECT_CONTEXT_FORM",
	OBJECT_DB:           "OBJECT_DB",
}

// Object
// const OBJECT_CONTEXT string = "OBJECT_CONTEXT"
// const OBJECT_CONTEXT_FORM string = "OBJECT_CONTEXT_FORM"
// const OBJECT_DB = "OBJECT_DB"
