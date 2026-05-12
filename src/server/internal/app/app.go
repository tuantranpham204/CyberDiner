package app

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/controller"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/middleware"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/model"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/service"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}
	return db, nil
}

func NewRouter(db *gorm.DB) *gin.Engine {
	denylist := store.NewInMemoryDenylist()
	userStore := store.NewUserStore(db)
	authSvc := service.NewAuthService(userStore, denylist)
	authCtrl := controller.NewAuthController(authSvc)

	r := gin.Default()
	r.Use(corsMiddleware())

	v1 := r.Group("/api/v1")
	auth := v1.Group("/auth")

	auth.POST("/signup", authCtrl.SignUp)
	auth.POST("/signin", authCtrl.SignIn)

	protected := auth.Group("", middleware.JWTAuth(denylist))
	protected.POST("/signout", authCtrl.SignOut)

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept-Language")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
