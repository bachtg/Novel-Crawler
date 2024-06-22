package source_adapter

import (
	"novel_crawler/constant"
	"novel_crawler/internal/model"
)

type SourceAdapter interface {
	Connect() SourceAdapter
	GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)

	GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error)

	GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error)

	GetAllGenres() ([]*model.Genre, error)

	GetDomain() string
}

type SourceAdapterManager struct {
	CurrentSource   *SourceAdapter
	SourceMapping   map[string]*SourceAdapter
	PriorityMapping map[string]int
}

func (sourceAdapterManager *SourceAdapterManager) AddNewSource(sources ...*SourceAdapter) error {

	if sourceAdapterManager.CurrentSource == nil {
		sourceAdapterManager.SourceMapping = make(map[string]*SourceAdapter)
		sourceAdapterManager.PriorityMapping = make(map[string]int)
		if len(sources) > 0 {
			sourceAdapterManager.CurrentSource = sources[0]
		}
	}

	for _, source := range sources {
		sourceDomain := (*source).GetDomain()
		sourceAdapterManager.SourceMapping[sourceDomain] = source
		sourceAdapterManager.PriorityMapping[sourceDomain] = len(sourceAdapterManager.SourceMapping) - 1
	}

	if sourceAdapterManager.CurrentSource != nil {
		return nil
	}

	return &model.Err{
		Code:    constant.NoSourceFound,
		Message: "No source found",
	}
}

func (sourceAdapterManager *SourceAdapterManager) RemoveSource(sourceDomain string) error {
	if sourceAdapterManager.SourceMapping != nil {
		delete(sourceAdapterManager.SourceMapping, sourceDomain)
		delete(sourceAdapterManager.PriorityMapping, sourceDomain)

		index := 0
		for key := range sourceAdapterManager.PriorityMapping {
			sourceAdapterManager.PriorityMapping[key] = index
			index += 1
		}
	}
	return nil
}

func (sourceAdapterManager *SourceAdapterManager) GetAllSources() ([]*model.Source, error) {
	numSource := len(sourceAdapterManager.SourceMapping)

	if numSource == 0 {
		return nil, &model.Err{
			Code:    constant.NoSourceFound,
			Message: "No source found",
		}
	}

	sources := make([]*model.Source, numSource)

	for key, value := range sourceAdapterManager.PriorityMapping {
		sources[value] = &model.Source{
			Id: key,
		}
	}

	return sources, nil
}
