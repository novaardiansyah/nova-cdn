package utils

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

type ImageVersion struct {
	Prefix  string
	Width   uint
	Quality int
}

type ProcessedImage struct {
	FileName string
	FilePath string
	FileSize uint32
}

var ImageVersions = []ImageVersion{
	{Prefix: "small", Width: 300, Quality: 35},
	{Prefix: "medium", Width: 900, Quality: 55},
	{Prefix: "large", Width: 1600, Quality: 65},
}

func ProcessImage(inputPath, outputDir, baseName string) ([]ProcessedImage, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	var results []ProcessedImage

	for _, version := range ImageVersions {
		resized := resize.Resize(version.Width, 0, img, resize.Lanczos3)

		versionFileName := fmt.Sprintf("%s-%s.jpg", version.Prefix, strings.TrimSuffix(baseName, filepath.Ext(baseName)))
		versionFilePath := filepath.Join(outputDir, versionFileName)

		outFile, err := os.Create(versionFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to create output file: %w", err)
		}

		err = jpeg.Encode(outFile, resized, &jpeg.Options{Quality: version.Quality})
		outFile.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to encode jpeg: %w", err)
		}

		info, err := os.Stat(versionFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to stat output file: %w", err)
		}

		relPath := "images/gallery/" + versionFileName

		results = append(results, ProcessedImage{
			FileName: versionFileName,
			FilePath: relPath,
			FileSize: uint32(info.Size()),
		})
	}

	return results, nil
}
