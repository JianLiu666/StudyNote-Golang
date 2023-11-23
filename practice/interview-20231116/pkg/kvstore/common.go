package kvstore

func genHashKey() string {
	return "list"
}

func genPageKey(pageKey string) string {
	return "page/" + pageKey
}
