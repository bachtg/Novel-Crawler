package repository

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"novel_crawler/constant"
	"novel_crawler/internal/model"
	"novel_crawler/util"
	"regexp"
	"strconv"
	"strings"
)

type TangThuVienAdapter struct {
	collector *colly.Collector
}

func NewTangThuVienAdapter() SourceAdapter {
	collector := colly.NewCollector(
		colly.AllowedDomains("truyen.tangthuvien.vn"),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"),
		colly.Async(true),
		colly.MaxDepth(1),
	)
	collector.AllowURLRevisit = true
	return &TangThuVienAdapter{collector: collector}
}

var totalGenre map[string]string

func checkExist(url string, listGenre []*model.Genre) (bool, string) {
	id := strings.Split(url, "the-loai/")[1]
	for _, genre := range listGenre {
		if genre.Id == id {
			return true, ""
		}
	}
	return false, id
}

func (tangThuVienAdapter *TangThuVienAdapter) GetDomain() string {
	return "truyen.tangthuvien.vn"
}

// Complete
func (tangThuVienAdapter *TangThuVienAdapter) GetAllGenres() ([]*model.Genre, error) {
	var genres []*model.Genre
	tangThuVienAdapter.collector.OnHTML("a", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		name := e.ChildText("span.info i")
		if strings.Contains(url, "the-loai") {
			exist, id := checkExist(url, genres)
			if !exist {
				genres = append(genres, &model.Genre{
					Id:   id,
					Name: name,
				})
			}
		}
	})
	err := tangThuVienAdapter.collector.Visit("https://truyen.tangthuvien.vn/")
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: "Cannot visit url: " + "https://truyen.tangthuvien.vn/",
		}
	}
	tangThuVienAdapter.collector.Wait()
	return genres, nil
}

// Complete
func (tangThuVienAdapter *TangThuVienAdapter) GetNovelsByGenre(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage int
	)
	totalGenre = make(map[string]string)
	// Set map
	geners, _ := tangThuVienAdapter.GetAllGenres()
	for key, val := range geners {
		totalGenre[val.Id] = strconv.Itoa(key + 1)
	}
	tangThuVienAdapter.collector.OnHTML(".book-img-text li", func(e *colly.HTMLElement) {
		title := e.ChildText(".book-mid-info h4 a")
		titleHref := e.ChildAttr(".book-mid-info h4 a", "href")

		subs := strings.Split(titleHref, "/")
		subTitle := subs[len(subs)-1]
		image := e.ChildAttr("img", "src")
		chapterNumberStr := e.ChildText(".KIBoOgno")
		authorName := e.ChildText(".book-mid-info .author .name")
		authorHref := e.ChildAttrs(".book-mid-info .author .name", "href")[0]
		authorId := strings.Split(authorHref, "author=")[1]
		var authors []*model.Author

		authors = append(authors, &model.Author{
			Id:   authorId,
			Name: authorName,
		})
		novels = append(novels, &model.Novel{
			Id:         subTitle,
			Title:      title,
			CoverImage: image,
			Author:     authors,
			LatestChapter: &model.Chapter{
				Id:    chapterNumberStr,
				Title: "",
			},
		})
	})

	tangThuVienAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			num, _ := strconv.Atoi(child.Text)
			numPage = max(numPage, num)
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := tangThuVienAdapter.collector.Visit("https://truyen.tangthuvien.vn/tong-hop?ctg=" + totalGenre[request.GenreId])
	if err != nil {
		return &model.GetNovelsResponse{
				Novels:  nil,
				NumPage: 0,
			}, &model.Err{
				Code:    constant.InternalError,
				Message: err.Error(),
			}
	}

	tangThuVienAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}

func (tangThuVienAdapter *TangThuVienAdapter) GetNovelsByCategory(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage int
	)

	tangThuVienAdapter.collector.OnHTML(".book-img-text li", func(e *colly.HTMLElement) {
		title := e.ChildText(".book-mid-info h4 a")
		titleHref := e.ChildAttr(".book-mid-info h4 a", "href")

		subs := strings.Split(titleHref, "/")
		subTitle := subs[len(subs)-1]
		image := e.ChildAttr("img", "src")

		authorName := e.ChildText(".book-mid-info .author .name")
		authorHref := e.ChildAttrs(".book-mid-info .author .name", "href")[0]
		authorId := strings.Split(authorHref, "author=")[1]
		var authors []*model.Author

		authors = append(authors, &model.Author{
			Id:   authorId,
			Name: authorName,
		})
		novels = append(novels, &model.Novel{
			Id:         subTitle,
			Title:      title,
			CoverImage: image,
			Author:     authors,
		})
	})

	// get number of page
	tangThuVienAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			num, _ := strconv.Atoi(child.Text)
			numPage = max(numPage, num)
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})
	var url string
	if request.CategoryId == "truyen-hot" {
		url = "https://truyen.tangthuvien.vn/tong-hop?rank=nm&page=" + request.Page
	} else {
		url = "https://truyen.tangthuvien.vn/tong-hop?fns=ht&page=" + request.Page
	}

	err := tangThuVienAdapter.collector.Visit(url)
	if err != nil {
		return &model.GetNovelsResponse{
				Novels:  nil,
				NumPage: 0,
			}, &model.Err{
				Code:    constant.InternalError,
				Message: err.Error(),
			}
	}

	tangThuVienAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}

// Complete
func (tangThuVienAdapter *TangThuVienAdapter) GetDetailNovel(request *model.GetDetailNovelRequest) (*model.GetDetailNovelResponse, error) {
	fmt.Println("GetDetailNovel")
	res := &model.Novel{}
	var story_id string
	tangThuVienAdapter.collector.OnHTML(".book-detail-wrap", func(e *colly.HTMLElement) {
		name := e.ChildText(".book-info h1")
		image := e.ChildAttr(".book-img img", "src")
		var authors []*model.Author
		authorName := e.ChildText(".tag a:first-child")
		authorHref := e.ChildAttrs(".tag a:first-child", "href")[0]
		authorId := strings.Split(authorHref, "author=")[1]
		authors = append(authors, &model.Author{
			Name: authorName,
			Id:   authorId,
		})
		temp := e.ChildText(".nav-wrap ul li:nth-child(2)")
		intro := e.ChildText(".intro")
		re := regexp.MustCompile(`\d+`)
		matches := re.FindStringSubmatch(temp)
		chapterNumber, err := strconv.Atoi(matches[0])
		if err != nil {
			fmt.Println(err)
		}
		var genres []*model.Genre
		genreName := e.ChildText(".tag .red")
		genreHref := e.ChildAttr(".tag .red", "href")
		genreID := strings.Split(genreHref, "the-loai/")[1]

		genres = append(genres, &model.Genre{
			Id:   genreID,
			Name: genreName,
		})

		rateStr := e.ChildText("#myrate")
		rateFloat, _ := strconv.ParseFloat(strings.TrimSpace(rateStr), 32)

		story_id = e.ChildAttr("input[name=story_id]", "value")

		res = &model.Novel{
			Id:          request.NovelId,
			Title:       name,
			CoverImage:  image,
			Rate:        float32(rateFloat),
			Chapters:    nil,
			Author:      authors,
			Genre:       genres,
			Description: intro,
			LatestChapter: &model.Chapter{
				Id:      "",
				Title:   "Chương " + strconv.Itoa(chapterNumber),
				Content: "",
			},
			Status: "success",
		}
	})

	err := tangThuVienAdapter.collector.Visit("https://truyen.tangthuvien.vn/doc-truyen/" + request.NovelId)
	if err != nil {
		return &model.GetDetailNovelResponse{
				Novel:   nil,
				NumPage: 0,
			}, &model.Err{
				Code:    constant.InternalError,
				Message: err.Error(),
			}
	}

	tangThuVienAdapter.collector.Wait()

	chapters := tangThuVienAdapter.GetListChapters(story_id, request.Page)
	res.Chapters = chapters

	return &model.GetDetailNovelResponse{
			Novel:   res,
			NumPage: 1,
		},
		nil
}

func (tangThuVienAdapter *TangThuVienAdapter) GetNovelsByAuthor(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage int
	)

	tangThuVienAdapter.collector.OnHTML(".book-img-text li", func(e *colly.HTMLElement) {
		title := e.ChildText(".book-mid-info h4 a")
		titleHref := e.ChildAttr(".book-mid-info h4 a", "href")

		subs := strings.Split(titleHref, "/")
		subTitle := subs[len(subs)-1]
		image := e.ChildAttr("img", "src")

		authorName := e.ChildText(".book-mid-info .author .name")
		authorHref := e.ChildAttrs(".book-mid-info .author .name", "href")[0]
		authorId := strings.Split(authorHref, "author=")[1]
		var authors []*model.Author

		authors = append(authors, &model.Author{
			Id:   authorId,
			Name: authorName,
		})
		novels = append(novels, &model.Novel{
			Id:         subTitle,
			Title:      title,
			CoverImage: image,
			Author:     authors,
		})
	})

	// get number of page
	tangThuVienAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			num, _ := strconv.Atoi(child.Text)
			numPage = max(numPage, num)
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := tangThuVienAdapter.collector.Visit("https://truyen.tangthuvien.vn/tac-gia?author=" + request.AuthorId + "&page=" + request.Page)
	if err != nil {

		return &model.GetNovelsResponse{
				Novels:  nil,
				NumPage: 0,
			}, &model.Err{
				Code:    constant.InternalError,
				Message: err.Error(),
			}
	}

	tangThuVienAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}

func (tangThuVienAdapter *TangThuVienAdapter) GetListChapters(story_id string, page string) []*model.Chapter {
	var listChapters []*model.Chapter
	tangThuVienAdapter.collector.OnHTML(".cf", func(e *colly.HTMLElement) {

		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			chapterHref := child.Attr("href")
			chapterIDArr := strings.Split(chapterHref, "/")
			chapterID := chapterIDArr[len(chapterIDArr)-1]
			chapterName := child.Attr("title")
			listChapters = append(listChapters, &model.Chapter{
				Id:      chapterID,
				Title:   chapterName,
				Content: "",
			})
		})

	})
	pageUrl, _ := strconv.Atoi(page)
	err := tangThuVienAdapter.collector.Visit("https://truyen.tangthuvien.vn/doc-truyen/page/" + story_id + "?page=" + strconv.Itoa(pageUrl-1) + "&limit=60&web=1")
	if err != nil {
		return nil
	}

	tangThuVienAdapter.collector.Wait()
	return listChapters
}

func (tangThuVienAdapter *TangThuVienAdapter) GetDetailChapter(request *model.GetDetailChapterRequest) (*model.GetDetailChapterResponse, error) {
	var (
		detailChapterResponse = &model.GetDetailChapterResponse{}
		url                   = "https://truyen.tangthuvien.vn/doc-truyen/" + request.NovelId + "/" + request.ChapterId
	)

	tangThuVienAdapter.collector.OnHTML(".truyen-title a", func(e *colly.HTMLElement) {
		id := util.GetId(e.Attr("href"))
		title := e.Text
		detailChapterResponse.Novel = &model.Novel{
			Id:    id,
			Title: title,
		}
	})

	tangThuVienAdapter.collector.OnHTML(".chapter h2", func(e *colly.HTMLElement) {
		//id := util.GetId(e.Attr("href"))
		title := e.Text
		detailChapterResponse.CurrentChapter = &model.Chapter{
			Id:    "",
			Title: title,
		}
	})

	tangThuVienAdapter.collector.OnHTML(".chapter-c-content .box-chap", func(e *colly.HTMLElement) {
		detailChapterResponse.CurrentChapter.Content, _ = e.DOM.Html()
	})

	detailChapterResponse.NextChapter = &model.Chapter{
		Id: "",
	}
	detailChapterResponse.PreviousChapter = &model.Chapter{
		Id: "",
	}

	err := tangThuVienAdapter.collector.Visit(url)
	if err != nil {
		return nil, &model.Err{
			Code:    constant.InternalError,
			Message: err.Error(),
		}
	}

	tangThuVienAdapter.collector.Wait()

	return detailChapterResponse, nil
}

func (tangThuVienAdapter *TangThuVienAdapter) GetNovelsByKeyword(request *model.GetNovelsRequest) (*model.GetNovelsResponse, error) {
	var (
		novels  []*model.Novel
		numPage int
	)

	tangThuVienAdapter.collector.OnHTML(".book-img-text li", func(e *colly.HTMLElement) {
		title := e.ChildText(".book-mid-info h4 a")
		titleHref := e.ChildAttr(".book-mid-info h4 a", "href")

		subs := strings.Split(titleHref, "/")
		subTitle := subs[len(subs)-1]
		image := e.ChildAttr("img", "src")

		authorName := e.ChildText(".book-mid-info .author .name")
		authorHref := e.ChildAttrs(".book-mid-info .author .name", "href")[0]
		authorId := strings.Split(authorHref, "author=")[1]
		var authors []*model.Author

		authors = append(authors, &model.Author{
			Id:   authorId,
			Name: authorName,
		})
		novels = append(novels, &model.Novel{
			Id:         subTitle,
			Title:      title,
			CoverImage: image,
			Author:     authors,
		})
	})

	// get number of page
	tangThuVienAdapter.collector.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, child *colly.HTMLElement) {
			num, _ := strconv.Atoi(child.Text)
			numPage = max(numPage, num)
		})
		activePage, _ := strconv.Atoi(strings.Split(e.ChildText(".active"), " ")[0])
		numPage = max(numPage, activePage)
	})

	err := tangThuVienAdapter.collector.Visit("https://truyen.tangthuvien.vn/ket-qua-tim-kiem?term=" + request.Keyword + "&page=" + request.Page)
	if err != nil {
		return &model.GetNovelsResponse{
				Novels:  nil,
				NumPage: 0,
			}, &model.Err{
				Code:    constant.InternalError,
				Message: err.Error(),
			}
	}

	tangThuVienAdapter.collector.Wait()

	return &model.GetNovelsResponse{
		Novels:  novels,
		NumPage: numPage,
	}, nil
}
