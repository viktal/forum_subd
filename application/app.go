package api

import (
	"context"
	"fmt"
	"forum/application/common"
	"github.com/apsdehal/go-logger"
	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	UserHandler "forum/application/user/delivery/http"
	UserRepository "forum/application/user/repository"
	UserUseCase "forum/application/user/usecase"

	ThreadHandler "forum/application/thread/delivery/http"
	ThreadRepository "forum/application/thread/repository"
	ThreadUseCase "forum/application/thread/usecase"

	PostHandler "forum/application/post/delivery/http"
	PostRepository "forum/application/post/repository"
	PostUseCase "forum/application/post/usecase"

	ForumHandler "forum/application/forum/delivery/http"
	ForumRepository "forum/application/forum/repository"
	ForumUseCase "forum/application/forum/usecase"

	ServiceHandler "forum/application/service/delivery/http"
	ServiceRepository "forum/application/service/repository"
	ServiceUseCase "forum/application/service/usecase"

	"github.com/go-pg/pg/v9"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Listen  string    `yaml:"listen"`
	Db      *dbConfig `yaml:"db"`
	DocPath string    `yaml:"docPath"`
	Redis   string    `yaml:"redis_address"`
}

type Logger struct {
	InfoLogger  *logger.Logger
	ErrorLogger *logger.Logger
}

type App struct {
	config   Config
	log      *Logger
	doneChan chan bool
	route    *gin.Engine
	db       *pg.DB
}

func NewApp(config Config) *App {

	infoLogger, _ := logger.New("forum_subd", 1, os.Stdout)
	errorLogger, _ := logger.New("forum_subd", 2, os.Stderr)

	log := &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}

	infoLogger.SetLogLevel(logger.DebugLevel)

	r := gin.New()
	r.Use(common.RequestLogger(log.InfoLogger))
	r.Use(common.ErrorLogger(log.ErrorLogger))
	r.Use(common.ErrorMiddleware())
	r.Use(common.Recovery(log.ErrorLogger))

	// Only for requests WITHOUT credentials, the literal value "*" can be specified
	corsMiddleware := cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://127.0.0.1")
		},
		MaxAge: time.Hour,
	})

	r.Use(corsMiddleware)

	r.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	if config.DocPath != "" {
		r.Static("/doc/api", config.DocPath)
	} else {
		log.ErrorLogger.Warning("Document path is undefined")
	}

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Db.Host, config.Db.Port),
		User:     config.Db.User,
		Password: config.Db.Password,
		Database: config.Db.Name,
	})


	gin.Default()
	govalidator.SetFieldsRequiredByDefault(false)
	api := r.Group("/api")

	UserRep := UserRepository.NewPgRepository(db)
	userCase := UserUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, UserRep)
	UserHandler.NewRest(api.Group("/user"), userCase)

	ThreadRep := ThreadRepository.NewPgRepository(db)
	ThreadCase := ThreadUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, ThreadRep)
	ThreadHandler.NewRest(api.Group("/thread"), ThreadCase)

	ServiceRep := ServiceRepository.NewPgRepository(db)
	ServiceCase := ServiceUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, ServiceRep)
	ServiceHandler.NewRest(api.Group("/service"), ServiceCase)

	ForumRep := ForumRepository.NewPgRepository(db)
	ForumCase := ForumUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, ForumRep, UserRep)
	ForumHandler.NewRest(api.Group("/forum"), ForumCase)

	PostRep := PostRepository.NewPgRepository(db)
	PostCase := PostUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, PostRep, UserRep, ForumRep, ThreadRep)
	PostHandler.NewRest(api.Group("/post"), PostCase)


	app := App{
		config:   config,
		log:      log,
		route:    r,
		doneChan: make(chan bool, 1),
		db:       db,
	}

	return &app
}

func (a *App) Run() {

	srv := &http.Server{
		Addr:    a.config.Listen,
		Handler: a.route,
	}

	go func() {
		a.log.InfoLogger.Infof("Start listening on %s", a.config.Listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.log.ErrorLogger.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
	case <-a.doneChan:
	}
	a.log.InfoLogger.Info("Shutdown Server (timeout of 1 seconds) ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		mes := fmt.Sprint("Server Shutdown:", err)
		a.log.ErrorLogger.Fatal(mes)
	}

	<-ctx.Done()
	a.log.InfoLogger.Info("Server exiting")
}

func (a *App) Close() {
	a.db.Close()
	a.doneChan <- true
}
