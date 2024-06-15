package business

import (
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/internal/repository"
	"sync"
	"fmt"
)

type Service struct {
	*repository.SourceAdapterManager
	*repository.ExporterManager
}

func NewService(sourceAdapterManager *repository.SourceAdapterManager, exporterManager *repository.ExporterManager) *Service {
	return &Service{sourceAdapterManager, exporterManager}
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
	sourceNum := len(service.SourceAdapterManager.SourceMapping)
	var wg sync.WaitGroup
	wg.Add(sourceNum)
	resultChan := make(chan *model.GetDetailChapterResponse, sourceNum)
	var errRes error
	for key, value := range service.SourceAdapterManager.SourceMapping {
		go func(key string, value *repository.SourceAdapter) {
			defer wg.Done()
			adapter := value
			resp, err := (*adapter).GetDetailChapter(&model.GetDetailChapterRequest{
				NovelId: request.NovelId,
				ChapterId: request.ChapterId,
				SourceDomain: key,
			})

			errRes = err

			if err == nil {
				resp.CurrentSource = key
				resultChan <- resp
			}
		}(key, value)
	}
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	source := *service.SourceAdapterManager.CurrentSource
	var sources []string
	var respone *model.GetDetailChapterResponse

	for result := range resultChan {
		if result.CurrentChapter.Title != "" {
            sources = append(sources, result.CurrentSource)
			if result.CurrentSource == request.SourceDomain {
				respone = result
			}
			if result.CurrentSource == source.GetDomain() {
				respone = result
			}
        }
	}
	if(errRes != nil) {
		return nil, errRes
	}
	novel, _ := source.GetDetailNovel(&model.GetDetailNovelRequest{
		NovelId: request.NovelId,
	})
	respone.Sources = sources
	respone.Novel.CoverImage = novel.Novel.CoverImage
	return respone, errRes
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
		ChapterId:    request.ChapterId,
		NovelId:      request.NovelId,
		SourceDomain: request.Domain,
	})

	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	if getDetailChapterResponse.CurrentChapter == nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Not found",
		}
	}

	if getDetailChapterResponse.CurrentChapter.Content == "" {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Not found",
		}
	}

	var exporter repository.Exporter

	if request.Type == "PDF" {
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

func (service *Service) GetAllTypes() []string {
	var result []string
	
	for key, _ := range service.ExporterManager.ExporterMapping{ 
		fmt.Println(key)
		result = append(result, key)
	}
	return result
}
