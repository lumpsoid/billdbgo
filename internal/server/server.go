package server

import (
	repository "billdb/internal/repository/bill"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

type Config struct {
	DbPath        string
	TemplatesPath string
	StaticPath    string
}

type Server struct {
	Config   *Config
	Echo     *echo.Echo
	BillRepo repository.BillRepository
}

func LoadConfig() (*Config, error) {
	var cfg Config
	var present bool

	cfg.DbPath, present = os.LookupEnv("BILLDB_DB_PATH")
	if !present {
		return nil, fmt.Errorf("BILLDB_DB_PATH not set")
	}
	cfg.TemplatesPath, present = os.LookupEnv("BILLDB_TEMPLATE_PATH")
	if !present {
		return nil, fmt.Errorf("BILLDB_TEMPLATE_PATH not set")
	}
	cfg.StaticPath, present = os.LookupEnv("BILLDB_STATIC_PATH")
	if !present {
		return nil, fmt.Errorf("BILLDB_STATIC_PATH not set")
	}
	return &cfg, nil
}

func Get(path string, handler func(s *Server) echo.HandlerFunc) func(s *Server) *echo.Route {
	return func(s *Server) *echo.Route {
		return s.Echo.GET(path, handler(s))
	}
}

//func Post(path string, handler echo.HandlerFunc) func(s *Server) *echo.Route {
//	return func(s *Server) *echo.Route {
//		return s.Echo.POST(path, handler)
//	}
//}

func Post(path string, handler func(s *Server) echo.HandlerFunc) func(s *Server) *echo.Route {
	return func(s *Server) *echo.Route {
		return s.Echo.POST(path, handler(s))
	}
}

func Put(path string, handler func(s *Server) echo.HandlerFunc) func(s *Server) *echo.Route {
	return func(s *Server) *echo.Route {
		return s.Echo.PUT(path, handler(s))
	}
}
