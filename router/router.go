package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"novel_crawler/internal/repository"
	"novel_crawler/internal/repository/source_adapter"

	"novel_crawler/config"
	"novel_crawler/internal/business"
	"novel_crawler/middleware"
)

func Start() {
	router := gin.Default()
	router.Use(middleware.CorsMiddleware())

	sourceAdapterManager := source_adapter.SourceAdapterManager{}

	truyenFullAdapter := &source_adapter.TruyenFullAdapter{}
	truyenFullAdapterWrapper := truyenFullAdapter.Connect()

	tangThuVienAdapter := &source_adapter.TangThuVienAdapter{}
	tangThuVienAdapterWrapper := tangThuVienAdapter.Connect()

	err := sourceAdapterManager.AddNewSource(&tangThuVienAdapterWrapper, &truyenFullAdapterWrapper)
	if err != nil {
		config.Cfg.Logger.Error(err.Error())
		panic(err)
	}

	exporterManager := repository.ExporterManager{}
	PDFExporter := repository.NewPDFExporter()
	EpubExporter := repository.NewEpubExporter()
	err = exporterManager.AddNewExporter(&PDFExporter, &EpubExporter)
	if err != nil {
		config.Cfg.Logger.Error(err.Error())
		panic(err)
	}

	novelService := business.NewService(&sourceAdapterManager, &exporterManager)
	novelHandler := business.NewHandler(novelService)

	router.GET("/novels/:novel_id", novelHandler.GetDetailNovel)
	router.GET("/novels/:novel_id/:chapter_id", novelHandler.GetDetailChapter)
	router.GET("/novels", novelHandler.GetNovels)

	router.GET("/genres", novelHandler.GetAllGenres)

	router.GET("/sources", novelHandler.GetAllSources)
	router.POST("/sources/:source_id", novelHandler.RegisterNewSourceAdapter)
	router.PATCH("/sources", novelHandler.UpdateSourcePriority)
	router.DELETE("/sources/:domain", novelHandler.RemoveSourceAdapter)

	router.GET("/types", novelHandler.GetTypes)
	router.DELETE("/types/:type_id", novelHandler.DeleteType)

	router.POST("/downloads", novelHandler.Download)

	config.Cfg.Logger.Info("Server's running on", zap.String("address", config.Cfg.Address))
	_ = router.Run()
}
