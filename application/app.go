package api

import (
	"fmt"
	PostHandler "forum/application/post/delivery/http"
	PostRepository "forum/application/post/repository"
	PostUseCase "forum/application/post/usecase"
	ThreadHandler "forum/application/thread/delivery/http"
	ThreadRepository "forum/application/thread/repository"
	ThreadUseCase "forum/application/thread/usecase"
	UserHandler "forum/application/user/delivery/http"
	UserRepository "forum/application/user/repository"
	UserUseCase "forum/application/user/usecase"
	"github.com/apsdehal/go-logger"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"

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

	params := make(map[string]interface{})
	params["search_path"] = "main"
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Db.Host, config.Db.Port),
		User:     config.Db.User,
		Password: config.Db.Password,
		Database: config.Db.Name,
		OnConnect: func(conn *pg.Conn) error {
			_, err := conn.Exec("set search_path=?", "main")

			//_, err = conn.Exec("LOAD 'auto_explain'")
			//_, err = conn.Exec("SET auto_explain.log_analyze TO on;")
			//_, err = conn.Exec("SET auto_explain.log_min_duration TO 300;")
			return err
		},
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

	//router.GET("/debug/*path", pprofhandler.PprofHandler)
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
