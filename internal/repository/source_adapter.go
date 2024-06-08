package repository

import "novel_crawler/internal/model"

type SourceAdapter interface {
	GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)

	GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error)

	GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error)

	GetAllGenres() ([]*model.Genre, error)
}
