package namelib

// The name for the embedded auth package.
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

// The name for the router package.
// Also used outside the router, but for interaction with it, e.g. MDDL_ERROR.
type RouterNames struct {
	URL_PATTERN            string
	COOKIE_CSRF_TOKEN      string
	SKIP_NEXT_PAGE         string
	REDIRECT_ERROR         string
	SERVER_ERROR           string
	SERVER_FORBIDDEN_ERROR string
}

var ROUTER = RouterNames{
	URL_PATTERN:            "URL_PATTERN",
	COOKIE_CSRF_TOKEN:      "CSRF_TOKEN",
	SKIP_NEXT_PAGE:         "SKIP_NEXT_PAGE",
	REDIRECT_ERROR:         "REDIRECT_ERROR",
	SERVER_ERROR:           "SERVER_ERROR",
	SERVER_FORBIDDEN_ERROR: "SERVER_ERROR",
}

// The name for the package object.
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

type TagNames struct {
	DB_MAPPER_NAME        string
	DB_MAPPER_EMPTY       string
	DB_MAPPER_DATE_F      string
	FORM_MAPPER_NAME      string
	FORM_MAPPER_EMPTY     string
	FORM_MAPPER_EXTENSION string
	FORM_MAPPER_NIL       string
	DTO                   string
}

var TAGS = TagNames{
	DB_MAPPER_NAME:        "db",
	DB_MAPPER_EMPTY:       "empty",
	DB_MAPPER_DATE_F:      "date-f",
	FORM_MAPPER_NAME:      "form",
	FORM_MAPPER_EMPTY:     "empty",
	FORM_MAPPER_EXTENSION: "ext",
	FORM_MAPPER_NIL:       "nil",
	DTO:                   "dto",
}
