package e

var MsgFlags = map[CODE]string{
	SUCCESS:                   "success",
	INVALID_PARAMS:            "invalid parameters",
	ERROR:                     "unknown error happened",
	ERROR_ADD_DUPLICATED_USER: "add duplicated user",
	ERROR_USER_NOT_FOUND:      "user not found",
}

func GetMsg(code CODE) string {
	if msg, exists := MsgFlags[code]; exists {
		return msg
	}
	return MsgFlags[ERROR]
}
