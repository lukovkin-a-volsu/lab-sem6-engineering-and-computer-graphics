package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type BitmapFileHeader struct {
	BfType      [2]byte
	BfSize      uint32
	BfReserved1 uint16
	BfReserved2 uint16
	BfOffBits   uint32
}

type BitmapInfoHeader struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

func readHeaders(file *os.File) (*BitmapFileHeader, *BitmapInfoHeader, error) {
	var fileHeader BitmapFileHeader
	if err := binary.Read(file, binary.LittleEndian, &fileHeader); err != nil {
		return nil, nil, err
	}

	if string(fileHeader.BfType[:]) != "BM" {
		return nil, nil, errors.New("Файл соответствует формату BMP")
	}

	var infoHeader BitmapInfoHeader
	if err := binary.Read(file, binary.LittleEndian, &infoHeader); err != nil {
		return nil, nil, err
	}

	return &fileHeader, &infoHeader, nil
}

func loadBMP(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileHeader, infoHeader, err := readHeaders(file)
	if err != nil {
		return nil, err
	}

	if infoHeader.BiCompression != 0 {
		return nil, errors.New("BMP изображения с компрессией не поддерживаются")
	}

	var palette []color.RGBA
	if infoHeader.BiBitCount <= 8 {
		paletteSize := int(infoHeader.BiClrUsed)
		if paletteSize == 0 {
			paletteSize = 1 << infoHeader.BiBitCount
		}

		palette = make([]color.RGBA, paletteSize)
		log.Printf("Индексированное изображение: %d бит, палитра из %d цветов", infoHeader.BiBitCount, len(palette))

		for i := 0; i < paletteSize; i++ {
			var entry struct {
				B, G, R, A uint8
			}
			if err := binary.Read(file, binary.LittleEndian, &entry); err != nil {
				return nil, err
			}
			palette[i] = color.RGBA{R: entry.R, G: entry.G, B: entry.B, A: 255}
		}
	} else {
		log.Printf("Полноцветное изображение: %d бит (RGB)", infoHeader.BiBitCount)
	}

	if _, err = file.Seek(int64(fileHeader.BfOffBits), io.SeekStart); err != nil {
		return nil, err
	}

	width := int(infoHeader.BiWidth)
	absHeight := int(int32(math.Abs(float64(infoHeader.BiHeight))))
	img := image.NewRGBA(image.Rect(0, 0, width, absHeight))

	bitsPerPixel := int(infoHeader.BiBitCount)
	bytesPerRow := ((width*bitsPerPixel + 31) / 32) * 4
	buffer := make([]byte, bytesPerRow)

	for y := 0; y < absHeight; y++ {
		if _, err := io.ReadFull(file, buffer); err != nil {
			return nil, err
		}

		destY := y
		if infoHeader.BiHeight > 0 {
			destY = absHeight - 1 - y
		}

		switch bitsPerPixel {
		case 1:
			for x := 0; x < width; x++ {
				byteIdx := x / 8
				bitIdx := 7 - (x % 8)
				bit := (buffer[byteIdx] >> bitIdx) & 1
				img.Set(x, destY, palette[bit])
			}

		case 4:
			for x := 0; x < width; x++ {
				byteIdx := x / 2
				nibble := buffer[byteIdx]
				if x%2 == 0 {
					nibble >>= 4
				} else {
					nibble &= 0x0F
				}
				img.Set(x, destY, palette[nibble])
			}

		case 8:
			for x := 0; x < width; x++ {
				idx := buffer[x]
				img.Set(x, destY, palette[idx])
			}

		case 24:
			for x := 0; x < width; x++ {
				offset := x * 3
				img.Set(x, destY, color.RGBA{
					R: buffer[offset+2],
					G: buffer[offset+1],
					B: buffer[offset],
					A: 255,
				})
			}

		case 32:
			for x := 0; x < width; x++ {
				offset := x * 4
				img.Set(x, destY, color.RGBA{
					R: buffer[offset+2],
					G: buffer[offset+1],
					B: buffer[offset],
					A: buffer[offset+3],
				})
			}

		default:
			return nil, fmt.Errorf("неподдерживаемая глубина пикселей: %d", bitsPerPixel)
		}
	}

	return img, nil
}

func run() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: bmpviewer <file.bmp>")
		return
	}

	img, err := loadBMP(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	cfg := pixelgl.WindowConfig{
		Title:  "BMP Viewer - " + os.Args[1],
		Bounds: pixel.R(0, 0, float64(img.Bounds().Dx()), float64(img.Bounds().Dy())),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	for !win.Closed() {
		win.Clear(color.White)
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
