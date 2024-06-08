package business

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"novel_crawler/constant"
	"novel_crawler/internal/model"
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

func (handler *Handler) GetNovels(ctx *gin.Context) {
	getNovelsRequest := &model.GetNovelsRequest{
		Page:       ctx.Query("page"),
		Keyword:    ctx.Query("search"),
		AuthorId:   ctx.Query("author"),
		CategoryId: ctx.Query("category"),
		GenreId:    ctx.Query("genre"),
	}

	getNovelsResponse, err := handler.Service.GetNovels(getNovelsRequest)

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
