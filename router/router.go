package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"novel_crawler/config"
	"novel_crawler/internal/business"
	"novel_crawler/internal/repository"
)

func Start() {
	router := gin.Default()

	sourceAdapterManager := repository.SourceAdapterManager{}

	//truyenFullAdapter := repository.NewTruyenFullAdapter()
	tangThuVienAdapter := repository.NewTangThuVienAdapter()

	err := sourceAdapterManager.AddNewSource(&tangThuVienAdapter)
	if err != nil {
		config.Cfg.Logger.Error(err.Error())
		panic(err)
	}

	novelService := business.NewService(&sourceAdapterManager)
	novelHandler := business.NewHandler(novelService)

	router.GET("/genres", novelHandler.GetAllGenres)
	router.GET("/novels/:novel_id", novelHandler.GetDetailNovel)
	router.GET("/novels/:novel_id/:chapter_id", novelHandler.GetDetailChapter)
	router.GET("/novels", novelHandler.GetNovels)
	router.GET("/sources", novelHandler.GetAllSources)
	router.POST("/sources/:domain", novelHandler.RegisterSourceAdapter)
	router.PATCH("/sources", novelHandler.UpdateSourcePriority)

	//router.POST("/downloads", novelHandler.Download)

	config.Cfg.Logger.Info("Server's running on", zap.String("address", config.Cfg.Address))
	_ = router.Run(config.Cfg.Address)
}
