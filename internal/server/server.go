package server

import (
	"billdb/internal/repository"

	"github.com/labstack/echo/v4"
)

type Config struct {
  Db struct {
    Path string
  }
  TemplatesPath string
  StaticPath string
}

type Server struct {
  Config Config
  Echo     *echo.Echo
	BillRepo repository.BillRepository
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
