package business

import (
	"novel_crawler/internal/model"
	"novel_crawler/internal/repository"
)

type Service struct {
	repository.SourceAdapter
}

func NewService(sourceAdapter repository.SourceAdapter) *Service {
	return &Service{sourceAdapter}
}

func (service *Service) GetAllGenres() ([]*model.Genre, error) {
	return service.SourceAdapter.GetAllGenres()
}

func (service *Service) GetNovelsByGenre(genreId string, page string) (*model.GetNovelsResponse, error) {
	if page == "" {
		page = "1"
	}
	return service.SourceAdapter.GetNovelsByGenre(genreId, page)
}

func (service *Service) GetNovelsByCategory(categoryId string, page string) (*model.GetNovelsResponse, error) {
	if page == "" {
		page = "1"
	}
	return service.SourceAdapter.GetNovelsByCategory(categoryId, page)
}

func (service *Service) GetDetailNovel(novelId string, page string) (*model.Novel, int, error) {
	if page == "" {
		page = "1"
	}
	return service.SourceAdapter.GetDetailNovel(novelId, page)
}

func (service *Service) GetNovelsByAuthor(authorId string, page string) (*model.GetNovelsResponse, error) {
	if page == "" {
		page = "1"
	}
	return service.SourceAdapter.GetNovelsByAuthor(authorId, page)
}

func (service *Service) GetNovelsByKeyword(keyword string, page string) (*model.GetNovelsResponse, error) {
	if page == "" {
		page = "1"
	}
	return service.SourceAdapter.GetNovelByKeyword(keyword, page)
}

func (service *Service) GetDetailChapter(novelId string, chapterId string) (*model.DetailChapterResponse, error) {
	return service.SourceAdapter.GetDetailChapter(novelId, chapterId)
}
