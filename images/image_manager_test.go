package images

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func CreateImage(t *testing.T, imagePath string) {
	t.Helper()

	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	cyan := color.RGBA{100, 200, 200, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2:
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2:
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	f, err := os.Create(imagePath)
	require.NoError(t, err)
	err = png.Encode(f, img)
	require.NoError(t, err)
}

func TestUnitNewImageManager(t *testing.T) {
	testPath := "test_depot"

	testCases := []struct {
		name            string
		imagesDepotPath string
	}{
		{
			name:            "Success with path not ending on '/'",
			imagesDepotPath: testPath,
		},
		{
			name:            "Success with path ending on '/'",
			imagesDepotPath: testPath + "/",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageManager := NewImageManager(tc.imagesDepotPath)

			require.Equal(t, "test_depot", imageManager.imagesDepotPath)
		})
	}
}

func TestUnitGetImagePath(t *testing.T) {
	imageManager := ImageManager{
		imagesDepotPath: "test_path",
	}

	fileName := "test_file.png"

	require.Equal(t, "test_path/"+fileName, imageManager.GetImagePath(fileName))
}

func TestRemoveImage(t *testing.T) {
	testImageName := "unit_test_image.png"

	testCases := []struct {
		name       string
		fileExists bool
	}{
		{
			name:       "Success with existing file",
			fileExists: true,
		},
		{
			name:       "Fail with non-existing file",
			fileExists: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			imageManager := NewImageManager("./test_images_depot")

			filePath := fmt.Sprintf("%s/%s", imageManager.imagesDepotPath, testImageName)
			if tc.fileExists {
				CreateImage(t, filePath)
				time.Sleep(1 * time.Second)
				require.NoError(t, imageManager.RemoveImage(testImageName))
			} else {
				require.Error(t, imageManager.RemoveImage(testImageName))
			}
		})
	}
}
