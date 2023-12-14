package e

var MsgFlags = map[CODE]string{
	SUCCESS:        "success",
	INVALID_PARAMS: "invalid parameters",
	ERROR:          "unknown error happened",
}

func GetMsg(code CODE) string {
	if msg, exists := MsgFlags[code]; exists {
		return msg
	}
	return MsgFlags[ERROR]
}
