package api

import (
	"fmt"
	UserHandler "forum/application/user/delivery/http"
	UserRepository "forum/application/user/repository"
	UserUseCase "forum/application/user/usecase"
	"github.com/apsdehal/go-logger"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

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
	"os"
)

type dbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Listen string    `yaml:"listen"`
	Db     *dbConfig `yaml:"db"`
}

type Logger struct {
	InfoLogger  *logger.Logger
	ErrorLogger *logger.Logger
}

type App struct {
	config   Config
	log      *Logger
	doneChan chan bool
	route    *fasthttprouter.Router
	db       *pg.DB
}

func NewApp(config Config) *App {

	infoLogger, _ := logger.New("forum_subd", 1, os.Stdout)
	errorLogger, _ := logger.New("forum_subd", 2, os.Stderr)

	log := &Logger{
		InfoLogger:  infoLogger,
		ErrorLogger: errorLogger,
	}

	infoLogger.SetLogLevel(logger.ErrorLevel)

	router := fasthttprouter.New()


	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Db.Host, config.Db.Port),
		User:     config.Db.User,
		Password: config.Db.Password,
		Database: config.Db.Name,
	})


	UserRep := UserRepository.NewPgRepository(db)
	userCase := UserUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, UserRep)
	UserHandler.NewRest(router, userCase)

	ThreadRep := ThreadRepository.NewPgRepository(db)
	ThreadCase := ThreadUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, ThreadRep)
	ThreadHandler.NewRest(router, ThreadCase)

	ServiceRep := ServiceRepository.NewPgRepository(db)
	ServiceCase := ServiceUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, ServiceRep)
	ServiceHandler.NewRest(router, ServiceCase)

	ForumRep := ForumRepository.NewPgRepository(db)
	ForumCase := ForumUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, ForumRep, UserRep)
	ForumHandler.NewRest(router, ForumCase)

	PostRep := PostRepository.NewPgRepository(db)
	PostCase := PostUseCase.NewUseCase(log.InfoLogger, log.ErrorLogger, PostRep, UserRep, ForumRep, ThreadRep)
	PostHandler.NewRest(router, PostCase)

	//router.GET("/debug/pprof/", func(c *fasthttp.RequestCtx) { pprof.Index(c.Writer, c.Request) })
	//router.GET("/debug/pprof/cmdline", func(c *fasthttp.RequestCtx) { pprof.Cmdline(c.Writer, c.Request) })
	//router.GET("/debug/pprof/profile", func(c *fasthttp.RequestCtx) { pprof.Profile(c.Writer, c.Request) })
	//router.GET("/debug/pprof/symbol", func(c *fasthttp.RequestCtx) { pprof.Symbol(c.Writer, c.Request) })
	//router.GET("/debug/pprof/trace", func(c *fasthttp.RequestCtx) { pprof.Trace(c.Writer, c.Request) })
	//
	//router.POST("/debug/pprof/", func(c *fasthttp.RequestCtx) { pprof.Index(c.Writer, c.Request) })
	//router.POST("/debug/pprof/cmdline", func(c *fasthttp.RequestCtx) { pprof.Cmdline(c.Writer, c.Request) })
	//router.POST("/debug/pprof/profile", func(c *fasthttp.RequestCtx) { pprof.Profile(c.Writer, c.Request) })
	//router.POST("/debug/pprof/symbol", func(c *fasthttp.RequestCtx) { pprof.Symbol(c.Writer, c.Request) })
	//router.POST("/debug/pprof/trace", func(c *fasthttp.RequestCtx) { pprof.Trace(c.Writer, c.Request) })

	app := App{
		config:   config,
		log:      log,
		route:    router,
		doneChan: make(chan bool, 1),
		db:       db,
	}

	return &app
}

func (a *App) Run() {
	_ = fasthttp.ListenAndServe(":5000", a.route.Handler)
}

func (a *App) Close() {
	_ = a.db.Close()
}
