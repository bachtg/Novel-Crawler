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

	novels, numPage, err := handler.Service.GetNovelsByGenre(genreId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":  novels,
			"numPage": numPage,
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

	novels, numPage, err := handler.Service.GetNovelsByCategory(categoryId, page)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": err.Error(),
		})
		ctx.Abort()
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":  novels,
			"numPage": numPage,
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
