package e

var MsgFlags = map[CODE]string{
	SUCCESS:        "success",
	ERROR:          "unknown error happened",
	INVALID_PARAMS: "invalid parameters",
}

func GetMsg(code CODE) string {
	if msg, exists := MsgFlags[code]; exists {
		return msg
	}
	return MsgFlags[ERROR]
}
