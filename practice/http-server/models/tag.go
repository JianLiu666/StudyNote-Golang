package models

import (
	"time"

	"gorm.io/gorm"
)

type Tag struct {
	Model

	Name       string `json:"name"`
	CreatedBy  string `json:"created_by"`
	ModifiedBy string `json:"modified_by"`
	State      int    `json:"state"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) (err error) {
	timeNow := time.Now().Unix()
	t.CreatedOn = timeNow
	t.ModifiedOn = timeNow
	return nil
}

func (t *Tag) BeforeUpdate(tx *gorm.DB) (err error) {
	// TODO: fix this
	t.ModifiedOn = time.Now().Unix()
	return nil
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

func ExistTagByID(id int) (bool, error) {
	var tag Tag
	err := gormDB.Select("id").Where("id = ? AND deleted_on = ?", id, 0).First(&tag).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if tag.ID > 0 {
		return true, nil
	}
	return false, nil
}

func GetTags(pageNum int, pageSize int, maps any) ([]Tag, error) {
	var tags []Tag
	var err error

	if pageSize > 0 && pageNum > 0 {
		err = gormDB.Where(maps).Offset(pageNum).Limit(pageSize).Find(&tags).Error
	} else {
		err = gormDB.Where(maps).Find(&tags).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tags, err
}

func GetTagTotal(maps any) (int, error) {
	var count int64

	if err := gormDB.Model(&Tag{}).Where(maps).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
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

func EditTag(id int, data any) error {
	if err := gormDB.Model(&Tag{}).Where("id = ? AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func DeleteTag(id int) error {
	if err := gormDB.Where("id = ?", id).Delete(&Tag{}).Error; err != nil {
		return err
	}
	return nil
}

func CleanAllTag() (bool, error) {
	if err := gormDB.Unscoped().Where("deleted_on != ?", 0).Delete(&Tag{}).Error; err != nil {
		return false, err
	}
	return true, nil
}
