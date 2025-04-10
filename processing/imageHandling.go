package processing

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"time"

	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

var face font.Face

func InitFont() error {
	fontBytes, err := os.ReadFile("fonts/firacode.ttf")
	if err != nil {
		return err
	}
	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		return err
	}
	const fontSize = 30
	face, err = opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	return err
}

func downloadImage(url string) (image.Image, error) {
	client := http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received bad status code: %d", resp.StatusCode)
	}

	img, format, err := image.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", url, err)
	}

	if format != "jpeg" && format != "png" {
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}

	return img, nil
}

func cropCenter(img image.Image, size int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	var crop image.Rectangle
	if w > h {
		x0 := (w - h) / 2
		crop = image.Rect(x0, 0, x0+h, h)
	} else {
		y0 := (h - w) / 2
		crop = image.Rect(0, y0, w, y0+w)
	}

	cropped := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(cropped, cropped.Bounds(), img, crop, draw.Over, nil)
	return cropped
}

func blankImage(size int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Gray16{}}, image.Point{}, draw.Src)
	return img
}

func drawLabel(dst draw.Image, text string, x, y, maxWidth int) {
	if len(text) > 24 {
		text = text[:23] + "…"
	}
	textWidth := font.MeasureString(face, text).Ceil()
	startX := x + (maxWidth-textWidth)/2
	startY := y

	paddingX := 6
	paddingY := 4
	boxX0 := startX - paddingX
	boxY0 := startY - face.Metrics().Ascent.Ceil() - paddingY
	boxX1 := startX + textWidth + paddingX
	boxY1 := startY + paddingY

	draw.Draw(dst, image.Rect(boxX0, boxY0, boxX1, boxY1), &image.Uniform{color.Black}, image.Point{}, draw.Over)

	drawer := &font.Drawer{
		Dst:  dst,
		Src:  &image.Uniform{color.White},
		Face: face,
		Dot:  fixed.P(startX, startY),
	}
	drawer.DrawString(text)
}

func GenerateAnimeGridImage(animeMap map[int64]Anime, cellSize int, labelHeight int, outputPath string) *image.RGBA {

	var animeList []Anime
	for _, a := range animeMap {
		animeList = append(animeList, a)
	}
	sort.Slice(animeList, func(i, j int) bool {
		return animeList[i].Count > animeList[j].Count
	})

	total := len(animeList)
	gridSize := int(math.Ceil(math.Sqrt(float64(total))))
	finalImg := image.NewRGBA(image.Rect(0, 0, gridSize*cellSize, gridSize*(cellSize+labelHeight)))

	for i := 0; i < gridSize*gridSize; i++ {
		x := (i % gridSize) * cellSize
		y := (i / gridSize) * (cellSize + labelHeight)

		var img image.Image
		var name string

		if i < total {
			a := animeList[i]
			name = a.Name
			downloaded, err := downloadImage(a.Cover)
			if err != nil {
				log.Printf("Failed to download image: %v", err)
				img = blankImage(cellSize)
			} else {
				img = cropCenter(downloaded, cellSize)
			}
		} else {
			img = blankImage(cellSize)
			name = ""
		}

		draw.Draw(finalImg, image.Rect(x, y, x+cellSize, y+cellSize), img, image.Point{0, 0}, draw.Over)
		drawLabel(finalImg, name, x, y+30, cellSize)

	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil
	}
	defer outFile.Close()

	return finalImg
}
