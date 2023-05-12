package app_error

// ID
const (
	BAD_REQUEST_ID     = "ID10000"
	UNAUTHENTICATED_ID = "ID10001"
	USER_ID_EXIST      = "ID1234"
	USER_EMAIL_EXIST   = "ID1235"
)

// Message
const (
	BAD_REQUEST_MESSSAGE     = "Bad request"
	UNAUTHENTICATED_MESSAGE  = "Un authorize"
	USER_ID_EXIST_MESSAGE    = "User id was exist"
	USER_EMAIL_EXIST_MESSAGE = "User email was exist"
)

func GetErrorDict() map[string]string {
	innserMap := map[string]string{
		BAD_REQUEST_ID:     BAD_REQUEST_MESSSAGE,
		UNAUTHENTICATED_ID: UNAUTHENTICATED_MESSAGE,
		USER_ID_EXIST:      USER_ID_EXIST_MESSAGE,
		USER_EMAIL_EXIST:   USER_EMAIL_EXIST_MESSAGE,
	}
	return innserMap
}

func GetErrorMessage(key string) string {
	errorDict := GetErrorDict()
	if v, found := errorDict[key]; found {
		return v
	}
	return ""
}
