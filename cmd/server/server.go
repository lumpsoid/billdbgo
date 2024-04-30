package main

import (
	"billdb/internal/repository"
	"billdb/internal/server"
	flutter "billdb/internal/server/flutterapi"
	web "billdb/internal/server/webapi"
	"context"
	"database/sql"
  _ "modernc.org/sqlite"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func main() {
  var cfg server.Config
  data, err := os.ReadFile("./config.yaml")
  err = yaml.Unmarshal(data, &cfg)
  if err != nil {
    log.Fatal("Error on yaml unmarshal")
    return
  }

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	db, err := sql.Open("sqlite", cfg.Db.Path)
	if err != nil {
		logger.Fatal("Error on sqlite3 db open")
		return
	}
  defer db.Close()

	billRepo := repository.NewSqliteBillRepository(db)
	s := server.Server{
		BillRepo: billRepo,
    Config: cfg,
	}

	t := &server.Template{
		Templates: template.Must(template.ParseGlob(cfg.TemplatesPath)),
	}

	e := echo.New()

	e.Renderer = t

  // call to /index-style.css will redirect to cfg.StaticPath/index-style.css
  e.Static("/", cfg.StaticPath)

	e.Logger.SetLevel(log.INFO)
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogError:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
				zap.Error(v.Error),
			)
			return nil
		},
	}))
	s.Echo = e
	// handlers
	web.RegisterWebRoutes(&s)
  flutter.FlutterApiRoutes(&s)
  

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":1323"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
