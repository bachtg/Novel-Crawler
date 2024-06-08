package business

import (
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/internal/repository"
)

type Service struct {
	repository.SourceAdapter
}

func NewService(sourceAdapter repository.SourceAdapter) *Service {
	return &Service{sourceAdapter}
}

func (service *Service) GetNovels(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	if request.Page == "" {
		request.Page = "1"
	}
	if request.Keyword != "" {
		return service.SourceAdapter.GetNovelsByKeyword(request)
	}
	if request.GenreId != "" {
		return service.SourceAdapter.GetNovelsByGenre(request)
	}
	if request.CategoryId != "" {
		return service.SourceAdapter.GetNovelsByCategory(request)
	}
	if request.AuthorId != "" {
		return service.SourceAdapter.GetNovelsByAuthor(request)
	}
	return nil, &model.Err{
		Code:    constant.InvalidRequest,
		Message: "Invalid Request",
	}
}

func (service *Service) GetAllGenres() ([]*model.Genre, error) {
	return service.SourceAdapter.GetAllGenres()
}

func (service *Service) GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error) {
	if request.Page == "" {
		request.Page = "1"
	}
	return service.SourceAdapter.GetDetailNovel(request)
}

func (service *Service) GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error) {
	return service.SourceAdapter.GetDetailChapter(request)
}
