package app

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/controller"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/middleware"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/service"
	"github.com/tuantranpham204/CyberDiner.git/src/server/internal/store"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/logger"
	"github.com/tuantranpham204/CyberDiner.git/src/server/pkg/util"
	customvalidator "github.com/tuantranpham204/CyberDiner.git/src/server/pkg/validator"
	"gorm.io/gorm"
)

type App struct {
	Config         *Config
	DB             *gorm.DB
	Router         *gin.Engine
	JWT            *util.JWTManager
	Denylist       store.TokenDenylist
	AuthController    *controller.AuthController
	ProfileController *controller.ProfileController
	DocsController    *controller.DocsController
}

func New(cfg *Config) (*App, error) {
	db, err := store.NewDB(cfg.Database.DSN())
	if err != nil {
		return nil, err
	}
	if err := store.AutoMigrate(db); err != nil {
		return nil, err
	}

	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := customvalidator.Register(v); err != nil {
			return nil, err
		}
	}

	userStore := store.NewUserStore(db)

	var denylist store.TokenDenylist
	if cfg.Redis.URL != "" {
		client, rerr := store.NewRedisClient(cfg.Redis.URL)
		if rerr != nil {
			logger.L().Warnw("redis_unreachable_falling_back_to_memory",
				"url", cfg.Redis.URL, "error", rerr)
			denylist = store.NewMemoryDenylist()
		} else {
			logger.L().Infow("redis_connected", "url", cfg.Redis.URL)
			denylist = store.NewRedisDenylist(client, "jwt:denylist:")
		}
	} else {
		denylist = store.NewMemoryDenylist()
	}

	jwtMgr := util.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours, cfg.JWT.Issuer)
	authSvc := service.NewAuthService(userStore, denylist, jwtMgr)
	profileSvc := service.NewProfileService(userStore)
	authCtl := controller.NewAuthController(authSvc)
	profileCtl := controller.NewProfileController(profileSvc)
	docsCtl := controller.NewDocsController(cfg.Docs.SpecPath)

	r := gin.New()
	r.Use(middleware.Recovery(), middleware.RequestLogger(), middleware.CORS(cfg.CORS.AllowedOrigins))

	return &App{
		Config:         cfg,
		DB:             db,
		Router:         r,
		JWT:            jwtMgr,
		Denylist:       denylist,
		AuthController:    authCtl,
		ProfileController: profileCtl,
		DocsController:    docsCtl,
	}, nil
}
