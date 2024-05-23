package images

import (
	"fmt"
	"os"
)

type ImageManager struct {
	imagesDepotPath string
}

func NewImageManager(imagesDepotPath string) *ImageManager {
	return &ImageManager{
		imagesDepotPath: imagesDepotPath,
	}
}

func (imageManager *ImageManager) GetImagePath(imageName string) string {
	return fmt.Sprintf("%s/%s", imageManager.imagesDepotPath, imageName)
}

func (imageManager *ImageManager) RemoveImage(imageName string) error {
	imagePath := imageManager.GetImagePath(imageName)

	return os.Remove(imagePath)
}
