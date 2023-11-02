package models

import (
	"gorm.io/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func ExistTagByName(name string) (bool, error) {
	var tag Tag
	err := gormDB.Select("id").Where("name = ? AND deleted_on = ?", name, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}

	return false, nil
}

func AddTag(name string, state int, createdBy string) error {
	tag := Tag{
		Name:      name,
		State:     state,
		CreatedBy: createdBy,
	}

	if err := gormDB.Create(&tag).Error; err != nil {
		return err
	}

	return nil
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
