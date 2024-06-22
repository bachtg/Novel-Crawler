package business

import (
	"fmt"
	"plugin"
	"sync"

	"golang.org/x/sync/errgroup"

	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/internal/repository/exporter"
	"novel_crawler/internal/repository/source_adapter"
)

type Service struct {
	SourceAdapterManager *source_adapter.SourceAdapterManager
	ExporterManager      *exporter.ExporterManager
}

func NewService(sourceAdapterManager *source_adapter.SourceAdapterManager, exporterManager *exporter.ExporterManager) *Service {
	return &Service{sourceAdapterManager, exporterManager}
}

func (service *Service) GetAllSources() ([]*model.Source, error) {
	return service.SourceAdapterManager.GetAllSources()
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
	//var errRes error
	for key, value := range service.SourceAdapterManager.SourceMapping {
		go func(key string, value *source_adapter.SourceAdapter) {
			defer wg.Done()
			adapter := value
			resp, err := (*adapter).GetDetailChapter(&model.GetDetailChapterRequest{
				NovelId:      request.NovelId,
				ChapterId:    request.ChapterId,
				SourceDomain: key,
			})

			//errRes = err

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
	var response *model.GetDetailChapterResponse
	check := 0
	for result := range resultChan {
		if result.CurrentChapter.Title != "" {
			sources = append(sources, result.CurrentSource)
			if result.CurrentSource == request.SourceDomain {
				response = result
				check = 1
			}
			if result.CurrentSource == source.GetDomain() {
				response = result
				check = 1
			}
			if check == 0 {
				response = result
			}
		}
	}

	novel, _ := source.GetDetailNovel(&model.GetDetailNovelRequest{
		NovelId: request.NovelId,
	})

	if response == nil || response.Novel == nil || novel.Novel == nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Not found novel",
		}
	}
	response.Sources = sources
	response.Novel.CoverImage = novel.Novel.CoverImage
	return response, nil
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

func (service *Service) RemoveSourceAdapter(sourceDomain string) error {
	return service.SourceAdapterManager.RemoveSource(sourceDomain)
}

func (service *Service) RegisterNewSourceAdapter(sourceAdapterId string) error {
	path := fmt.Sprintf("./plugin/source_adapter_plugin/%s/%s.so", sourceAdapterId, sourceAdapterId)
	plg, err := plugin.Open(path)
	if err != nil {
		return err
	}
	symSourceAdapter, err := plg.Lookup("SourceAdapter")
	if err != nil {
		return &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}
	sourceAdapter, ok := symSourceAdapter.(source_adapter.SourceAdapter)
	if !ok {
		return &model.Err{
			Code:    constant.InternalError,
			Message: "Cannot add new source",
		}
	}
	sourceAdapter.Connect()
	err = service.SourceAdapterManager.AddNewSource(&sourceAdapter)
	if err != nil {
		return &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}
	return nil
}

func (service *Service) RegisterNewExporter(typeId string) error {
	path := fmt.Sprintf("./plugin/exporter_plugin/%s/%s.so", typeId, typeId)
	plg, err := plugin.Open(path)
	if err != nil {
		return err
	}
	symExporter, err := plg.Lookup("Exporter")
	if err != nil {
		return &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}
	newExporter, ok := symExporter.(exporter.Exporter)
	if !ok {
		return &model.Err{
			Code:    constant.InternalError,
			Message: "Cannot add new exporter",
		}
	}
	newExporter.New()
	err = service.ExporterManager.AddNewExporter(&newExporter)
	if err != nil {
		return &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}
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

	exp := service.ExporterManager.ExporterMapping[request.Type]

	bytesData, err := (*exp).Generate("<p>" + getDetailChapterResponse.CurrentChapter.Content + "</p>\n")
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

	for key := range service.ExporterManager.ExporterMapping {
		result = append(result, key)
	}
	return result
}

func (service *Service) DeleteType(extension string) error {
	err := service.ExporterManager.RemoveExporter(extension)
	return err
}

func GetNovelsGoRoutine(f func(*model.GetNovelsRequest) (*model.GetNovelsResponse, error), numPage int, request *model.GetNovelsRequest) []*model.Novel {
	var novels []*model.Novel
	numGoroutines := numPage / 2
	size := numPage

	chunkSize := size / numGoroutines
	var mu sync.Mutex
	g := errgroup.Group{}
	for i := 0; i < numGoroutines; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if i == numGoroutines-1 {
			end = size
		}
		g.Go(func() error {
			func(start, end int) {
				var partialNovels []*model.Novel
				mu.Lock()
				for j := start; j < end; j++ {
					requestTemp := request
					requestTemp.Page = fmt.Sprintf("%d", j+1)
					resp, err := f(requestTemp)

					if err != nil {
						return
					}

					partialNovels = append(partialNovels, resp.Novels...)

				}
				novels = append(novels, partialNovels...)
				mu.Unlock()
			}(start, end)
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return novels
}
