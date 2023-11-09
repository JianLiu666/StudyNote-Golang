package repository

import "interview20231109/question-1/model"

type ArticleRepo interface {
	GetWithPagination(page, size int) []model.Article
}
