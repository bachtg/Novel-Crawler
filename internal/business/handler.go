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

	request := &model.GetDetailNovelRequest{
		NovelId: novelId,
		Page:    page,
	}

	getDetailNovelResponse, err := handler.Service.GetDetailNovel(request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	getDetailNovelResponse.Novel.Id = novelId
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novel":   getDetailNovelResponse.Novel,
			"numPage": getDetailNovelResponse.NumPage,
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
	request := &model.GetDetailChapterRequest{
		NovelId:   ctx.Param("novel_id"),
		ChapterId: ctx.Param("chapter_id"),
	}

	getDetailChapterResponse, err := handler.Service.GetDetailChapter(request)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":           getDetailChapterResponse.Novel,
			"current_chapter":  getDetailChapterResponse.CurrentChapter,
			"previous_chapter": getDetailChapterResponse.PreviousChapter,
			"next_chapter":     getDetailChapterResponse.NextChapter,
		},
	})
}
