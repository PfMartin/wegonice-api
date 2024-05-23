package images

import (
	"fmt"
	"os"
	"strings"
)

type ImageManager struct {
	imagesDepotPath string
}

func NewImageManager(imagesDepotPath string) *ImageManager {
	fixedPath := imagesDepotPath

	if strings.HasSuffix(imagesDepotPath, "/") {
		fixedPath = fixedPath[:len(imagesDepotPath)-1]
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
