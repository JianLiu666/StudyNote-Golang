package article

import "interview20231109/question-1/model"

func ArticleWithPagination(page, size int) []model.Article {
	result := []model.Article{}

	skips := (page - 1) * size

	instance.db.Table("article").Select("*").Limit(size).Offset(skips).Find(&result)

	return result
}
