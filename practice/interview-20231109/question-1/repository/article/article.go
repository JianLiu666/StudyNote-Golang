package article

import (
	"interview20231109/question-1/model"
	"interview20231109/question-1/repository"

	"gorm.io/gorm"
)

type articleService struct {
	db *gorm.DB
}

func NewArticleService(db *gorm.DB) repository.ArticleRepo {
	return &articleService{
		db: db,
	}
}

func (a *articleService) GetWithPagination(page, size int) []model.Article {
	result := []model.Article{}

	skips := (page - 1) * size

	a.db.Table("article").Select("*").Limit(size).Offset(skips).Find(&result)

	return result
}
