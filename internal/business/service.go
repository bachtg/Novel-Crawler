package business

import (
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/internal/repository"
	"fmt"
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
	if request.SourceDomain != "" {
		source, exists := service.SourceAdapterManager.SourceMapping[request.SourceDomain]
		if exists {
			return (*source).GetDetailChapter(request)
		}
	}
	source := *service.SourceAdapterManager.CurrentSource
	return source.GetDetailChapter(request)
}

func (service *Service) UpdateSourcePriority(sources []string) error {
	updateSuccess := true
	newPriorityMapping := make(map[string]int)
	for index, source := range sources {
		if _, exist := service.SourceAdapterManager.SourceMapping[source]; !exist {
			updateSuccess = false
			break
		}
		newPriorityMapping[source] = index
	}
	if updateSuccess {
		service.SourceAdapterManager.PriorityMapping = newPriorityMapping
		for domain, priority := range service.SourceAdapterManager.PriorityMapping {
			if priority == 0 {
				service.SourceAdapterManager.CurrentSource = service.SourceAdapterManager.SourceMapping[domain]
				break
			}
		}
		return nil
	}
	return &model.Err{
		Code:    constant.InvalidRequest,
		Message: "Invalid source",
	}
}

func (service *Service) RegisterSourceAdapter(domain string) error {
	return nil
}

func (service *Service) Download(request *model.DownloadChapterRequest) (*model.DownloadChapterResponse, error) {

	getDetailChapterResponse, err := service.GetDetailChapter(&model.GetDetailChapterRequest{
		ChapterId: request.ChapterId,
		NovelId:   request.NovelId,
	})
	if err != nil {
		return nil, &model.Err{
            Code:    constant.InternalError,
            Message: err.Error(),
        }
	}
	var exporter repository.Exporter
	fmt.Println(request.ChapterId)
	fmt.Println(request.NovelId)
	fmt.Println(request.Type)

	if (request.Type == "pdf") {
		exporter = repository.NewPDFExporter()
	} else {
		exporter = repository.NewEpubExporter()
	}

	bytesData, err := exporter.Generate("<p>" + getDetailChapterResponse.CurrentChapter.Content + "</p>\n")
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	filename := "[" + getDetailChapterResponse.Novel.Title + "] " + getDetailChapterResponse.CurrentChapter.Title

	return &model.DownloadChapterResponse{
		Filename:  filename,
		BytesData: bytesData,
	}, nil
}
