package main

import (
	"fmt"
	"novel_crawler/internal/repository/source_adapter"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"

	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/util"
)

type DTruyenAdapter struct {
	collector *colly.Collector
}

func (dtruyenAdapter *DTruyenAdapter) Connect() source_adapter.SourceAdapter {
	dtruyenAdapter.collector = colly.NewCollector(
		colly.AllowedDomains("dtruyen.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.Async(true),
		colly.MaxDepth(1),
	)
	dtruyenAdapter.collector.AllowURLRevisit = true
	return dtruyenAdapter
}

func (dtruyenAdapter *DTruyenAdapter) GetAllGenres() ([]*model.Genre, error) {
	var genres []*model.Genre
	url := "https://dtruyen.com"
	fmt.Println(url)
	dtruyenAdapter.collector.OnHTML(".categories a", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Attr("title")
		genres = append(genres, &model.Genre{
			Id:   id,
			Name: title,
		})
	})
	err := dtruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Cannot visit url: " + url,
		}
	}
	dtruyenAdapter.collector.Wait()
	return genres, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetNovels(url string) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage = 1
	)

	dtruyenAdapter.collector.OnHTML(".list-stories .story-list", func(e *colly.HTMLElement) {
		id := util.GetId(e.ChildAttr(".thumb", "href"))
		title := e.ChildAttr(".thumb", "title")
		var authors []*model.Author
		authors = append(authors, &model.Author{
			Name: e.ChildText(`[itemprop="author"]`),
		})
		coverImage := e.ChildAttr(".thumb > img", "data-layzr")
		lastChapter := e.ChildText(".last-chapter")
		latestChapter := &model.Chapter{
			Title: lastChapter,
		}
		novels = append(novels, &model.Novel{
			Id:            id,
			Title:         title,
			Author:        authors,
			CoverImage:    coverImage,
			LatestChapter: latestChapter,
		})
	})

	dtruyenAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			num, _ := strconv.Atoi(child.Text)
			numPage = max(numPage, num)
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := dtruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	dtruyenAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://dtruyen.com" + "/" + request.GenreId + "/" + request.Page
	getNovelsResponse, err := dtruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	var key string

	if request.CategoryId == "hoan-thanh" {
		key = "truyen-full"
	} else {
		key = "a"
	}
	url := "https://dtruyen.com" + "/" + key + "/" + request.Page
	getNovelsResponse, err := dtruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error) {
	var (
		novel    *model.Novel
		authors  []*model.Author
		genres   []*model.Genre
		chapters []*model.Chapter
		numPage  = 1
		url      = "https://dtruyen.com" + "/" + request.NovelId + "/" + request.Page
	)

	dtruyenAdapter.collector.OnHTML("#story-detail", func(e *colly.HTMLElement) {
		coverImage := e.ChildAttr("img", "src")
		title := e.ChildText(".title")
		description, _ := e.DOM.Find(".description").Html()
		rate, _ := strconv.ParseFloat(e.ChildAttr(".rate-holder", "data-score"), 32)

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

		status := e.ChildText(".infos p:nth-of-type(5)")

		novel = &model.Novel{
			Title:       title,
			Rate:        float32(rate),
			Author:      authors,
			CoverImage:  coverImage,
			Description: description,
			Genre:       genres,
			Status:      status,
		}
	})

	dtruyenAdapter.collector.OnHTML("#chapters .chapters", func(el *colly.HTMLElement) {
		el.ForEach("a", func(_ int, child *colly.HTMLElement) {
			id := util.GetId(child.Attr("href"))
			chapters = append(chapters, &model.Chapter{
				Id:    id,
				Title: child.Text,
			})
			fmt.Println(chapters[len(chapters)-1])
		})
		novel.Chapters = chapters
	})

	err := dtruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	dtruyenAdapter.collector.Wait()

	return &model.GetDetailNovelResponse{
		Novel:   novel,
		NumPage: numPage,
	}, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://dtruyen.com" + "/searching/" + request.Keyword + "/" + request.Page
	getNovelsResponse, err := dtruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error) {
	var (
		novel          = &model.Novel{}
		currentChapter = &model.Chapter{}
		prevChapter    = &model.Chapter{}
		nextChapter    = &model.Chapter{}
		url            = "https://dtruyen.com" + "/" + request.NovelId + "/" + request.ChapterId
	)
	fmt.Println("url")

	dtruyenAdapter.collector.OnHTML(".story-title a", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Text
		novel.Id = id
		novel.Title = title
	})

	dtruyenAdapter.collector.OnHTML(".chapter-title", func(e *colly.HTMLElement) {
		//id := util.GetId(e.Attr("href"))
		title := e.Text
		//currentChapter.Id = id
		currentChapter.Title = title
	})

	dtruyenAdapter.collector.OnHTML("#chapter-content", func(e *colly.HTMLElement) {
		currentChapter.Content, _ = e.DOM.Html()
	})

	dtruyenAdapter.collector.OnHTML(".chapter-button", func(e *colly.HTMLElement) {
		previousChapterId := util.GetId(e.ChildAttr(`a[title="Chương Trước"]`, "href"))
		prevChapter.Id = previousChapterId

		nextChapterId := util.GetId(e.ChildAttr(`[title="Chương Sau"]`, "href"))
		nextChapter.Id = nextChapterId
	})

	err := dtruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	dtruyenAdapter.collector.Wait()

	return &model.GetDetailChapterResponse{
		Novel:           novel,
		CurrentChapter:  currentChapter,
		PreviousChapter: prevChapter,
		NextChapter:     nextChapter,
	}, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := "https://dtruyen.com" + "/tac-gia/" + request.AuthorId + "/" + request.Page
	getNovelsResponse, err := dtruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (dtruyenAdapter *DTruyenAdapter) GetDomain() string {
	return "dtruyen.com"
}

var SourceAdapter DTruyenAdapter
