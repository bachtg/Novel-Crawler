package repository

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"

	"novel_crawler/config"
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/util"
)

type NetTruyenAdapter struct {
	collector *colly.Collector
}

func NewNetTruyenAdapter() SourceAdapter {
	collector := colly.NewCollector(
		colly.AllowedDomains("dtruyen.com"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.Async(true),
		colly.MaxDepth(1),
	)
	collector.AllowURLRevisit = true
	return &NetTruyenAdapter{collector: collector}
}

func (netTruyenAdapter *NetTruyenAdapter) GetAllGenres() ([]*model.Genre, error) {
	var genres []*model.Genre
	url := config.Cfg.NetTruyenBaseUrl
	fmt.Println(url)
	netTruyenAdapter.collector.OnHTML(".categories a", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Attr("title")
		genres = append(genres, &model.Genre{
			Id:   id,
			Name: title,
		})
	})
	err := netTruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Cannot visit url: " + url,
		}
	}
	netTruyenAdapter.collector.Wait()
	return genres, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetNovels(url string) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage = 1
	)

	netTruyenAdapter.collector.OnHTML(".list-stories .story-list", func(e *colly.HTMLElement) {
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

	netTruyenAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			num, _ := strconv.Atoi(child.Text)
			numPage = max(numPage, num)
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := netTruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	netTruyenAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := config.Cfg.NetTruyenBaseUrl + "/" + request.GenreId + "/" + request.Page
	getNovelsResponse, err := netTruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	var key string

	if request.CategoryId == "hoan-thanh" {
		key = "truyen-full"
	} else {
		key = "a"
	}
	url := config.Cfg.NetTruyenBaseUrl + "/" + key + "/" + request.Page
	getNovelsResponse, err := netTruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error) {
	var (
		novel    *model.Novel
		authors  []*model.Author
		genres   []*model.Genre
		chapters []*model.Chapter
		numPage  = 1
		url      = config.Cfg.NetTruyenBaseUrl + "/" + request.NovelId + "/" + request.Page
	)

	netTruyenAdapter.collector.OnHTML("#story-detail", func(e *colly.HTMLElement) {
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

	netTruyenAdapter.collector.OnHTML("#chapters .chapters", func(el *colly.HTMLElement) {
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

	err := netTruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	netTruyenAdapter.collector.Wait()

	return &model.GetDetailNovelResponse{
		Novel:   novel,
		NumPage: numPage,
	}, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := config.Cfg.NetTruyenBaseUrl + "/searching/" + request.Keyword + "/" + request.Page
	getNovelsResponse, err := netTruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error) {
	var (
		novel          = &model.Novel{}
		currentChapter = &model.Chapter{}
		prevChapter    = &model.Chapter{}
		nextChapter    = &model.Chapter{}
		url            = config.Cfg.NetTruyenBaseUrl + "/" + request.NovelId + "/" + request.ChapterId
	)
	fmt.Println("url")

	netTruyenAdapter.collector.OnHTML(".story-title a", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Text
		novel.Id = id
		novel.Title = title
	})

	netTruyenAdapter.collector.OnHTML(".chapter-title", func(e *colly.HTMLElement) {
		//id := util.GetId(e.Attr("href"))
		title := e.Text
		//currentChapter.Id = id
		currentChapter.Title = title
	})

	netTruyenAdapter.collector.OnHTML("#chapter-content", func(e *colly.HTMLElement) {
		currentChapter.Content, _ = e.DOM.Html()
	})

	netTruyenAdapter.collector.OnHTML(".chapter-button", func(e *colly.HTMLElement) {
		previousChapterId := util.GetId(e.ChildAttr(`a[title="Chương Trước"]`, "href"))
		prevChapter.Id = previousChapterId

		nextChapterId := util.GetId(e.ChildAttr(`[title="Chương Sau"]`, "href"))
		nextChapter.Id = nextChapterId
	})

	err := netTruyenAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	netTruyenAdapter.collector.Wait()

	return &model.GetDetailChapterResponse{
		Novel:           novel,
		CurrentChapter:  currentChapter,
		PreviousChapter: prevChapter,
		NextChapter:     nextChapter,
	}, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	url := config.Cfg.NetTruyenBaseUrl + "/tac-gia/" + request.AuthorId + "/" + request.Page
	getNovelsResponse, err := netTruyenAdapter.GetNovels(url)
	if err != nil {
		return nil, err
	}
	return getNovelsResponse, nil
}

func (netTruyenAdapter *NetTruyenAdapter) GetDomain() string {
	return "dtruyen.com"
}
