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
	url := "https://truyenfull.vn/"
	return service.SourceAdapter.GetAllGenres(url)
}

func (service *Service) GetNovelsByGenre(genreId string, page string) ([]*model.Novel, int, error) {
	if page == "" {
		page = "1"
	}
	url := "https://truyenfull.vn/the-loai/" + genreId + "/trang-" + page
	return service.SourceAdapter.GetNovelsByGenre(url)
}

func (service *Service) GetNovelsByCategory(categoryId string, page string) ([]*model.Novel, int, error) {
	if page == "" {
		page = "1"
	}
	url := "https://truyenfull.vn/danh-sach/" + categoryId + "/trang-" + page
	return service.SourceAdapter.GetNovelsByCategory(url)
}

func (service *Service) GetDetailNovel(novelId string, page string) (*model.Novel, int, error) {
	if page == "" {
		page = "1"
	}
	url := "https://truyenfull.vn/" + novelId + "/trang-" + page
	return service.SourceAdapter.GetDetailNovel(url)
}
