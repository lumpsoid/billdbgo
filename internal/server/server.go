package server

import (
	repository "billdb/internal/repository/bill"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/segmentio/ksuid"
)

type Server struct {
	Config   *Config
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

func Put(path string, handler func(s *Server) echo.HandlerFunc) func(s *Server) *echo.Route {
	return func(s *Server) *echo.Route {
		return s.Echo.PUT(path, handler(s))
	}
}

func IndexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1 // Return -1 if the element is not found
}

func CheckFormFile(file *multipart.FileHeader) error {
	if file.Size > 1048576 {
		return fmt.Errorf("file size is more then 1Mb: %d", file.Size)
	}
	fileType := file.Header.Get("Content-Type")
	typeArray := []string{"image/jpeg", "image/jpg", "image/png", "image/avif"}
	typeSupported := IndexOf(typeArray, fileType)
	if typeSupported == -1 {
		return errors.New("file type is not supported")
	}
	return nil
}

func UploadFileToServer(folderPath string, file *multipart.FileHeader) (string, error) {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer src.Close()

	// Define the destination path
	tmpName := ksuid.New().String()
	dstPath := filepath.Join(folderPath, tmpName)

	// Create the destination file
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer dst.Close()

	// Copy the file content to the destination file
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", fmt.Errorf("error copying file: %w", err)
	}

	return dstPath, nil
}
