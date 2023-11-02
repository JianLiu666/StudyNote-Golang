package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "請求參數錯誤",
	ERROR_EXIST_TAG:                 "已存在該標籤名稱",
	ERROR_EXIST_TAG_FAIL:            "獲取已存在標籤失敗",
	ERROR_NOT_EXIST_TAG:             "該標籤不存在",
	ERROR_GET_TAGS_FAIL:             "獲取所有標籤失敗",
	ERROR_COUNT_TAG_FAIL:            "統計標籤失敗",
	ERROR_ADD_TAG_FAIL:              "新增標籤失敗",
	ERROR_EDIT_TAG_FAIL:             "修改標籤失敗",
	ERROR_DELETE_TAG_FAIL:           "刪除標籤失敗",
	ERROR_EXPORT_TAG_FAIL:           "導出標籤失敗",
	ERROR_IMPORT_TAG_FAIL:           "導入標籤失敗",
	ERROR_NOT_EXIST_ARTICLE:         "該文章不存在",
	ERROR_ADD_ARTICLE_FAIL:          "新增文章失敗",
	ERROR_DELETE_ARTICLE_FAIL:       "刪除文章失敗",
	ERROR_CHECK_EXIST_ARTICLE_FAIL:  "檢查文章是否存在失敗",
	ERROR_EDIT_ARTICLE_FAIL:         "修改文章失敗",
	ERROR_COUNT_ARTICLE_FAIL:        "統計文章失敗",
	ERROR_GET_ARTICLES_FAIL:         "獲取多個文章失敗",
	ERROR_GET_ARTICLE_FAIL:          "獲取單一文章失敗",
	ERROR_GEN_ARTICLE_POSTER_FAIL:   "生成文章海報失敗",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鑑權失敗",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超時",
	ERROR_AUTH_TOKEN:                "Token生成失敗",
	ERROR_AUTH:                      "Token錯誤",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存圖片失敗",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "檢查圖片失敗",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校驗圖片錯誤，圖片格式或大小有問題",
}

func GetMsg(code int) string {
	if msg, ok := MsgFlags[code]; ok {
		return msg
	}
	return MsgFlags[ERROR]
}
