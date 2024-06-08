package repository

import (
	"github.com/gocolly/colly/v2"
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/util"
	"strconv"
	"strings"
)

type SourceAdapter interface {
	GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)
	GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error)

	GetDetailNovel(novelId string, page string) (*model.Novel, int, error)

	GetDetailChapter(novelId string, chapterId string) (*model.DetailChapterResponse, error)

	GetAllGenres() ([]*model.Genre, error)
}

type TruyenFullAdapter struct {
	collector *colly.Collector
}

func NewTruyenFullAdapter() SourceAdapter {
	collector := colly.NewCollector(
		colly.AllowedDomains("truyenfull.vn"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.Async(true),
		colly.MaxDepth(1),
	)
	collector.AllowURLRevisit = true
	return &TruyenFullAdapter{collector: collector}
}

func (truyenFullAdapter *TruyenFullAdapter) GetAllGenres() ([]*model.Genre, error) {
	var genres []*model.Genre
	url := "https://truyenfull.vn/"

	truyenFullAdapter.collector.OnHTML(".dropdown-menu a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		name := e.Text
		id := util.GetId(url)
		if strings.Contains(url, "the-loai") {
			genres = append(genres, &model.Genre{
				Id:   id,
				Name: name,
			})
		}
	})
	err := truyenFullAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Cannot visit url: " + url,
		}
	}
	truyenFullAdapter.collector.Wait()
	return genres, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetNovels(url string) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage = 1
	)

	truyenFullAdapter.collector.OnHTML(".row[itemscope]", func(e *colly.HTMLElement) {
		id := util.GetId(e.ChildAttr(".truyen-title > a", "href"))
		coverImage := e.ChildAttr("div[data-image]", "data-image")
		title := e.ChildAttr(".truyen-title > a", "title")
		latestChapter := e.ChildText(".text-info")

		var authors []*model.Author
		e.ForEach(".author", func(_ int, child *colly.HTMLElement) {
			authors = append(authors, &model.Author{Name: child.Text})
		})

		novels = append(novels, &model.Novel{
			Id:            id,
			Title:         title,
			Author:        authors,
			CoverImage:    coverImage,
			LatestChapter: latestChapter,
		})
	})

	// get number of page
	truyenFullAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			numPage = max(numPage, util.GetNumPage(child.Attr("href"), "trang-", "page="))
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := truyenFullAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	truyenFullAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://truyenfull.vn/the-loai/" + request.GenreId + "/trang-" + request.Page
	getNovelsResponse, err := truyenFullAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://truyenfull.vn/danh-sach/" + request.CategoryId + "/trang-" + request.Page
	getNovelsResponse, err := truyenFullAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetDetailNovel(novelId string, page string) (*model.Novel, int, error) {
	var (
		novel    *model.Novel
		authors  []*model.Author
		genres   []*model.Genre
		chapters []*model.Chapter
		numPage  = 1
		url      = "https://truyenfull.vn/" + novelId + "/trang-" + page
	)

	truyenFullAdapter.collector.OnHTML(".col-truyen-main", func(e *colly.HTMLElement) {
		coverImage := e.ChildAttr(".book > img", "src")
		title := e.ChildText(".title")
		description, _ := e.DOM.Find("div.desc-text.desc-text-full[itemprop='description']").Html()
		rate, _ := strconv.ParseFloat(e.ChildAttr(".rate-holder", "data-score"), 32)
		// just raw text
		//description := e.ChildText("div.desc-text.desc-text-full[itemprop='description']")

		e.ForEach("a[itemprop='author']", func(_ int, child *colly.HTMLElement) {
			authors = append(authors, &model.Author{
				Id:   util.GetId(child.Attr("href")),
				Name: child.Text,
			})
		})

		e.ForEach("a[itemprop='genre']", func(_ int, child *colly.HTMLElement) {
			genres = append(genres, &model.Genre{
				Id:   util.GetId(child.Attr("href")),
				Name: child.Text,
			})
		})

		e.ForEach(".list-chapter", func(_ int, el *colly.HTMLElement) {
			el.ForEach("a", func(_ int, child *colly.HTMLElement) {
				id := util.GetId(child.Attr("href"))
				chapters = append(chapters, &model.Chapter{
					Id:    id,
					Title: child.Text,
				})
			})

		})

		status := e.ChildText(".text-success")

		novel = &model.Novel{
			Title:       title,
			Rate:        float32(rate),
			Author:      authors,
			CoverImage:  coverImage,
			Description: description,
			Genre:       genres,
			Status:      status,
			Chapters:    chapters,
		}
	})

	// get number of page
	truyenFullAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			numPage = max(numPage, util.GetNumPage(child.Attr("href"), "trang-"))
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := truyenFullAdapter.collector.Visit(url)
	if err != nil {
		return nil, 0, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	truyenFullAdapter.collector.Wait()

	return novel, numPage, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://truyenfull.vn/tac-gia/" + request.AuthorId + "/trang-" + request.Page
	getNovelsResponse, err := truyenFullAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://truyenfull.vn/tim-kiem/?tukhoa=" + request.Keyword + "&page=" + request.Page
	getNovelsResponse, err := truyenFullAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (truyenFullAdapter *TruyenFullAdapter) GetDetailChapter(novelId string, chapterId string) (*model.DetailChapterResponse, error) {
	var (
		detailChapterResponse = &model.DetailChapterResponse{}
		url                   = "https://truyenfull.vn/" + novelId + "/" + chapterId
	)

	truyenFullAdapter.collector.OnHTML(".truyen-title", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Text
		detailChapterResponse.Novel = &model.Novel{
			Id:    id,
			Title: title,
		}
	})

	truyenFullAdapter.collector.OnHTML(".chapter-title", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Text
		detailChapterResponse.CurrentChapter = &model.Chapter{
			Id:    id,
			Title: title,
		}
	})

	truyenFullAdapter.collector.OnHTML(".chapter-c", func(e *colly.HTMLElement) {
		detailChapterResponse.CurrentChapter.Content, _ = e.DOM.Html()
	})

	truyenFullAdapter.collector.OnHTML(".btn-group", func(e *colly.HTMLElement) {
		previousChapterId := util.GetId(e.ChildAttr("#prev_chap", "href"))
		detailChapterResponse.PreviousChapter = &model.Chapter{
			Id: previousChapterId,
		}
		nextChapterId := util.GetId(e.ChildAttr("#next_chap", "href"))
		detailChapterResponse.NextChapter = &model.Chapter{
			Id: nextChapterId,
		}
	})

	err := truyenFullAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	truyenFullAdapter.collector.Wait()

	return detailChapterResponse, nil
}
