package images

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
)

type ImageManager struct {
	imagesDepotPath string
}

func NewImageManager(imagesDepotPath string) *ImageManager {
	fixedPath := imagesDepotPath

	if strings.HasSuffix(imagesDepotPath, "/") {
		fixedPath = fixedPath[:len(imagesDepotPath)-1]
	}

	if _, err := os.Stat(fixedPath); os.IsNotExist(err) {
		err = os.MkdirAll(fixedPath, os.ModePerm)
		if err != nil {
			log.Fatal().Msgf("Failed to create image depot path: %s", err)
		}
	}

	return &ImageManager{
		imagesDepotPath: fixedPath,
	}
}

func (imageManager *ImageManager) GetImagePath(imageName string) string {
	return fmt.Sprintf("%s/%s", imageManager.imagesDepotPath, imageName)
}

func (imageManager *ImageManager) RemoveImage(imageName string) error {
	imagePath := imageManager.GetImagePath(imageName)

	return os.Remove(imagePath)
}

func (imageManager *ImageManager) CreateUniqueName(imageName string) string {
	uniqueID := uuid.New().String()

	if flag.Lookup("test.v") != nil {
		uniqueID = "unique"
	}
	return fmt.Sprintf("%s-%s", uniqueID, imageName)
}
