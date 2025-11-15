package download

import (
	"ChainServer/internal/common/apperror"
	"fmt"
	"os"
)

type DownloadService struct {
}

func NewDownloadService() *DownloadService {
	return &DownloadService{}
}

func (s *DownloadService) DowloadNovachain(filename string) (string, *apperror.AppError) {
	filepath := fmt.Sprintf("/downloads/%s", filename)
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", apperror.NotFound(fmt.Sprintf("%s file not found", filename), nil)
	}

	return filepath, nil
}
