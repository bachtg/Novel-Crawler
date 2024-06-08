package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"novel_crawler/constant"
)

type Handler struct {
	*Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (handler *Handler) GetAllGenres(ctx *gin.Context) {
	genres, err := handler.Service.GetAllGenres()
	if err != nil {
		fmt.Println()
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"genres": genres,
		},
	})
}

func (handler *Handler) GetNovelsByGenre(ctx *gin.Context) {
	page := ctx.Query("page")
	genreId := ctx.Param("genre_id")

	getNovelsResponse, err := handler.Service.GetNovelsByGenre(genreId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":  getNovelsResponse.Novels,
			"numPage": getNovelsResponse.NumPage,
		},
	})
}

func (handler *Handler) Test(ctx *gin.Context) {
	novelId := ctx.Param("novel_id")
	chapterId := ctx.Param("chapter_id")
	ctx.JSON(http.StatusOK, gin.H{
		"novelId":   novelId,
		"chapterId": chapterId,
	})
}

func (handler *Handler) GetNovelByCategory(ctx *gin.Context) {
	page := ctx.Query("page")
	categoryId := ctx.Param("category_id")

	getNovelsResponse, err := handler.Service.GetNovelsByCategory(categoryId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":  getNovelsResponse.Novels,
			"numPage": getNovelsResponse.NumPage,
		},
	})
}

func (handler *Handler) GetDetailNovel(ctx *gin.Context) {
	page := ctx.Query("page")
	novelId := ctx.Param("novel_id")

	novel, numPage, err := handler.Service.GetDetailNovel(novelId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	novel.Id = novelId
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novel":   novel,
			"numPage": numPage,
		},
	})
}

func (handler *Handler) GetNovelByAuthor(ctx *gin.Context) {
	page := ctx.Query("page")
	authorId := ctx.Param("author_id")

	getNovelsResponse, err := handler.Service.GetNovelsByAuthor(authorId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":  getNovelsResponse.Novels,
			"numPage": getNovelsResponse.NumPage,
		},
	})
}

func (handler *Handler) GetNovelsByKeyword(ctx *gin.Context) {
	page := ctx.Query("page")
	keyword := ctx.Query("search")

	getNovelsResponse, err := handler.Service.GetNovelsByKeyword(keyword, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":  getNovelsResponse.Novels,
			"numPage": getNovelsResponse.NumPage,
		},
	})
}

func (handler *Handler) GetDetailChapter(ctx *gin.Context) {
	novelId := ctx.Param("novel_id")
	chapterId := ctx.Param("chapter_id")

	detailChapterResponse, err := handler.Service.GetDetailChapter(novelId, chapterId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":           detailChapterResponse.Novel,
			"current_chapter":  detailChapterResponse.CurrentChapter,
			"previous_chapter": detailChapterResponse.PreviousChapter,
			"next_chapter":     detailChapterResponse.NextChapter,
		},
	})
}
