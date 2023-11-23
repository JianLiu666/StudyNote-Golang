package e

var MsgFlags = map[CODE]string{
	SUCCESS:              "success",
	ERROR:                "unknown error happened",
	INVALID_PARAMS:       "invalid parameters",
	ERROR_MARSHAL:        "system error",
	ERROR_UNMARSHAL:      "system error",
	ERROR_REDIS_COMMAND:  "system error",
	ERROR_DATA_NOT_FOUND: "data not found",
}

func GetMsg(code CODE) string {
	if msg, exists := MsgFlags[code]; exists {
		return msg
	}
	return MsgFlags[ERROR]
}
