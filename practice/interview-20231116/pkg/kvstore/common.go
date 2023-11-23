package kvstore

func genListKey() string {
	return "list"
}

func genPageKey(pageKey string) string {
	return "page/" + pageKey
}
