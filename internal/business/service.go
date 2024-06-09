package business

import (
	"fmt"
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/internal/repository"
)

type Service struct {
	*repository.SourceAdapterManager
}

func NewService(sourceAdapterManager *repository.SourceAdapterManager) *Service {
	return &Service{sourceAdapterManager}
}

func (service *Service) GetNovels(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	source := *service.SourceAdapterManager.CurrentSource

	if request.Page == "" {
		request.Page = "1"
	}
	if request.Keyword != "" {
		return source.GetNovelsByKeyword(request)
	}
	if request.GenreId != "" {
		return source.GetNovelsByGenre(request)
	}
	if request.CategoryId != "" {
		return source.GetNovelsByCategory(request)
	}
	if request.AuthorId != "" {
		return source.GetNovelsByAuthor(request)
	}
	return nil, &model.Err{
		Code:    constant.InvalidRequest,
		Message: "Invalid Request",
	}
}

func (service *Service) GetAllGenres() ([]*model.Genre, error) {
	source := *service.SourceAdapterManager.CurrentSource
	return source.GetAllGenres()
}

func (service *Service) GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error) {
	if request.Page == "" {
		request.Page = "1"
	}
	source := *service.SourceAdapterManager.CurrentSource
	return source.GetDetailNovel(request)
}

func (service *Service) GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error) {
	fmt.Println(request)

	if request.SourceDomain != "" {
		source, exists := service.SourceAdapterManager.SourceMapping[request.SourceDomain]
		if exists {
			return (*source).GetDetailChapter(request)
		}
	}
	source := *service.SourceAdapterManager.CurrentSource
	return source.GetDetailChapter(request)
}

//func (service *Service) Download(request *model.DownloadChapterRequest) (*model.DownloadChapterResponse, error) {
//
//	getDetailChapterResponse, err := service.SourceAdapter.GetDetailChapter(&model.GetDetailChapterRequest{
//		ChapterId: request.ChapterId,
//		NovelId:   request.NovelId,
//	})
//
//	bytesData, err := Generate("<p>" + getDetailChapterResponse.CurrentChapter.Content + "</p>\n")
//	if err != nil {
//		return nil, &model.Err{
//			Code:    constant.InternalError,
//			Message: err.Error(),
//		}
//	}
//
//	filename := "[" + getDetailChapterResponse.Novel.Title + "] " + getDetailChapterResponse.CurrentChapter.Title
//
//	return &model.DownloadChapterResponse{
//		Filename:  filename,
//		BytesData: bytesData,
//	}, nil
//}
