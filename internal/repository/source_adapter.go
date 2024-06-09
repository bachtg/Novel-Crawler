package repository

import (
	"novel_crawler/constant"
	"novel_crawler/internal/model"
)

type SourceAdapter interface {
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
	CurrentSource *SourceAdapter
	SourceMapping map[string]*SourceAdapter
}

func (sourceAdapterManager *SourceAdapterManager) AddNewSource(sources ...*SourceAdapter) error {

	if sourceAdapterManager.SourceMapping == nil {
		sourceAdapterManager.SourceMapping = make(map[string]*SourceAdapter)
	}

	if len(sources) != 0 && sourceAdapterManager.CurrentSource == nil {
		sourceAdapterManager.CurrentSource = sources[0]
	}

	for _, source := range sources {
		sourceDomain := (*source).GetDomain()
		sourceAdapterManager.SourceMapping[sourceDomain] = source
	}

	if sourceAdapterManager.CurrentSource != nil {
		return nil
	}

	return &model.Err{
		Code:    constant.NoSourceFound,
		Message: "No source found",
	}
}
