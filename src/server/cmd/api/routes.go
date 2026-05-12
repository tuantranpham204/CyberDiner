package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/app"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	return app.NewRouter(db)
}
