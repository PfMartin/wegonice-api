package images

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
