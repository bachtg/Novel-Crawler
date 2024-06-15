package business

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"novel_crawler/constant"
	"novel_crawler/internal/model"
)

type Handler struct {
	*Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (handler *Handler) GetAllSources(ctx *gin.Context) {
	sources, err := handler.Service.GetAllSources()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"sources": sources,
		},
	})
}

func (handler *Handler) UpdateSourcePriority(ctx *gin.Context) {
	body := struct {
		Sources []string `json:"sources"`
	}{}
	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    constant.InvalidRequest,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	err = handler.Service.UpdateSourcePriority(body.Sources)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"body": body,
		},
	})
}

func (handler *Handler) GetAllGenres(ctx *gin.Context) {
	genres, err := handler.Service.GetAllGenres()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
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
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	getDetailNovelResponse.Novel.Id = novelId
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novel":   getDetailNovelResponse.Novel,
			"numPage": getDetailNovelResponse.NumPage,
			"perPage": getDetailNovelResponse.PerPage,
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
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
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

	body := struct {
		Domain string `json:"domain"`
	}{}
	_ = ctx.ShouldBindJSON(&body)

	request.SourceDomain = body.Domain

	getDetailChapterResponse, err := handler.Service.GetDetailChapter(request)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"novels":           getDetailChapterResponse.Novel,
			"current_chapter":  getDetailChapterResponse.CurrentChapter,
			"previous_chapter": getDetailChapterResponse.PreviousChapter,
			"next_chapter":     getDetailChapterResponse.NextChapter,
			"sources": getDetailChapterResponse.Sources,
		},
	})
}

func (handler *Handler) RegisterSourceAdapter(ctx *gin.Context) {
	domain := ctx.Param("domain")
	err := handler.Service.RegisterSourceAdapter(domain)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
	})
}

func (handler *Handler) Download(ctx *gin.Context) {
	downloadChapterRequest := &model.DownloadChapterRequest{}
	if err := ctx.ShouldBind(&downloadChapterRequest); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    constant.InvalidRequest,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	downloadChapterResponse, err := handler.Service.Download(downloadChapterRequest)

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    err.(*model.Err).Code,
			"message": err.Error(),
		})
		ctx.Abort()
		return
	}

	filename := fmt.Sprintf("%s.%s", downloadChapterResponse.Filename, downloadChapterRequest.Type)
	ctx.Header("Content-Disposition", "attachment; filename="+filename)
	ctx.Data(http.StatusOK, "application/"+downloadChapterRequest.Type, downloadChapterResponse.BytesData)
}

func (handler *Handler) GetTypes(ctx *gin.Context) {
	types := handler.Service.GetAllTypes();
	ctx.JSON(http.StatusOK, gin.H{
		"code": constant.Success,
		"data": gin.H{
			"types": types,
		},
	})
}
