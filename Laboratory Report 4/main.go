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
	Signature  [2]byte
	FileSize   uint32
	Reserved   uint32
	DataOffset uint32
}

type BitmapInfoHeader struct {
	Size            uint32
	Width           int32
	Height          int32
	Planes          uint16
	BitCount        uint16
	Compression     uint32
	ImageSize       uint32
	XpixelsPerM     int32
	YpixelsPerM     int32
	ColorsUsed      uint32
	ColorsImportant uint32
}

func readHeaders(file *os.File) (*BitmapFileHeader, *BitmapInfoHeader, error) {
	var fileHeader BitmapFileHeader
	if err := binary.Read(file, binary.LittleEndian, &fileHeader); err != nil {
		return nil, nil, err
	}

	if string(fileHeader.Signature[:]) != "BM" {
		return nil, nil, errors.New("not a BMP file")
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

	if infoHeader.Compression != 0 {
		return nil, errors.New("compressed BMP files are not supported")
	}

	var palette []color.RGBA
	if infoHeader.BitCount <= 8 {
		paletteSize := int(infoHeader.ColorsUsed)
		if paletteSize == 0 {
			paletteSize = 1 << infoHeader.BitCount
		}

		palette = make([]color.RGBA, paletteSize)
		for i := 0; i < paletteSize; i++ {
			var entry struct {
				B, G, R, A uint8
			}
			if err := binary.Read(file, binary.LittleEndian, &entry); err != nil {
				return nil, err
			}
			palette[i] = color.RGBA{entry.R, entry.G, entry.B, 255}
		}
	}

	if _, err = file.Seek(int64(fileHeader.DataOffset), io.SeekStart); err != nil {
		return nil, err
	}

	width := int(infoHeader.Width)
	absHeight := int(int32(math.Abs(float64(infoHeader.Height))))
	img := image.NewRGBA(image.Rect(0, 0, width, absHeight))

	bitsPerPixel := int(infoHeader.BitCount)
	bytesPerRow := ((width*bitsPerPixel + 31) / 32) * 4
	buffer := make([]byte, bytesPerRow)

	for y := 0; y < absHeight; y++ {
		if _, err := io.ReadFull(file, buffer); err != nil {
			return nil, err
		}

		destY := y
		if infoHeader.Height > 0 {
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
			return nil, fmt.Errorf("unsupported bit depth: %d", bitsPerPixel)
		}
	}

	return img, nil
}

func run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: bmpviewer <file.bmp>")
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
