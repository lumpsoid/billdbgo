package main

import (
	repository "billdb/internal/repository/bill"
	"billdb/internal/server"
	"billdb/internal/server/api"
	"billdb/internal/server/web"
	"context"
	"database/sql"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	cfg, err := server.LoadConfig()
	if err != nil {
		logger.Fatal(err.Error())
		return
	}

	db, err := sql.Open("sqlite3", cfg.DbPath)
	if err != nil {
		logger.Fatal("Error on sqlite3 db open")
		return
	}
	defer db.Close()

	billRepo := repository.NewSqliteBillRepository(db)
	s := server.Server{
		BillRepo: billRepo,
		Config:   cfg,
	}

	pattern := filepath.Join(cfg.TemplatesPath, "*.html")
	t := &server.Template{
		Templates: template.Must(template.ParseGlob(pattern)),
	}

	e := echo.New()

	e.Renderer = t

	// call to /index-style.css will redirect to cfg.StaticPath/index-style.css
	e.Static("/static", cfg.StaticPath)
	e.Static("/uploaded", cfg.QrPath)

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
	webGroup := e.Group("")
	webHandlers := web.NewWebHandlers(cfg, e, billRepo)
	webHandlers.RegisterRoutes(webGroup)
	api.ApiRoutes(&s)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	// Start server
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %v", err)
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
