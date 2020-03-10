package app

import (
	"datapreparation/pkg/cryptohelper"
	"datapreparation/profiles"
	"datapreparation/profiles/delivery"
	"datapreparation/profiles/usecase"
	"time"

	"github.com/buaazp/fasthttprouter"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type App struct {
	profilesUC profiles.UseCase

	server *fasthttp.Server
	logger *zap.Logger
}

func NewApp(profilesRepo profiles.Repository, dec cryptohelper.Decrypt) *App {
	server := &fasthttp.Server{}
	profilesUC := usecase.Init(profilesRepo, dec)

	return &App{
		server:     server,
		profilesUC: profilesUC,
		logger:     &zap.Logger{},
	}
}

func (a *App) WithLogger(logger *zap.Logger) *App {
	a.logger = logger
	return a
}

func (a *App) Run(port string) {
	a.server.Handler = a.newRouter().Handler

	a.server.ListenAndServe(":" + port)
}

func (a *App) newRouter() *fasthttprouter.Router {
	mw := a.mw
	router := fasthttprouter.New()

	profiles := delivery.NewHandler(a.profilesUC)

	router.POST("/profile", mw(profiles.ImportICCID))
	return router
}

// ------ MIDDLEWARE

func (a *App) mw(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		t := time.Now()
		reqid := uuid.NewV4().String()
		a.logger.Info("newrequest",
			zap.String("transaction-id", reqid),
			zap.String("path", string(ctx.Path())),
			zap.String("method", string(ctx.Method())),
		)
		ctx.SetUserValue("reqid", reqid)

		h(ctx)
		ctx.Response.Header.Set("Content-Type", "application/json")

		a.logger.Info("response",
			zap.String("transaction-id", reqid),
			zap.Int("status", ctx.Response.StatusCode()),
			zap.Duration("spent", time.Since(t)),
		)
	}
}
