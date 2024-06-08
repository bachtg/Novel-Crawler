package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"novel_crawler/internal/business"
	"novel_crawler/internal/repository"

	"novel_crawler/config"
)

func Start() {
	router := gin.Default()

	truyenFullAdapter := repository.NewTruyenFullAdapter()
	novelService := business.NewService(truyenFullAdapter)
	novelHandler := business.NewHandler(novelService)

	router.GET("/genres", novelHandler.GetAllGenres)
	router.GET("/genres/:genre_id", novelHandler.GetNovelsByGenre)
	router.GET("/novels/:novel_id", novelHandler.GetDetailNovel)
	router.GET("/novels/:novel_id/:chapter_id", novelHandler.GetDetailChapter)
	router.GET("/novels", novelHandler.GetNovelsByKeyword)
	router.GET("/authors/:author_id", novelHandler.GetNovelByAuthor)
	router.GET("/categories/:category_id", novelHandler.GetNovelByCategory)

	config.Cfg.Logger.Info("Server's running on", zap.String("address", config.Cfg.Address))
	_ = router.Run(config.Cfg.Address)
}
