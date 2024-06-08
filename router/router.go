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

	truyenFullAdapter := repository.NewSourceAdapter("truyenfull.vn")
	novelService := business.NewService(truyenFullAdapter)
	novelHandler := business.NewHandler(novelService)

	router.GET("/genres", novelHandler.GetAllGenres)
	router.GET("/genres/:genre_id", novelHandler.GetNovelsByGenre)
	router.GET("/categories/:category_id", novelHandler.GetNovelByCategory)
	router.GET("/novels/:novel_id", novelHandler.GetDetailNovel)
	router.GET("/novels/:novel_id/:chapter_id", novelHandler.Test)

	config.Cfg.Logger.Info("Server's running on", zap.String("address", config.Cfg.Address))
	_ = router.Run(config.Cfg.Address)
}
