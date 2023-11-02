package models

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      string `json:"state"`
}

func GetTags(pageNum int, pageSize int, maps any) (tags []Tag) {
	gormDB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags)
	return
}

func GetTagTotal(maps any) (count int) {
	var res int64
	gormDB.Model(&Tag{}).Where(maps).Count(&res)
	return int(res)
}
