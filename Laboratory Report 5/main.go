package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
)

// Загрузить изображение из файла
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	return img, err
}

// Сохранить изображение в файл
func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch path[len(path)-3:] {
	case "jpg", "peg":
		return jpeg.Encode(file, img, nil)
	case "png":
		return png.Encode(file, img)
	default:
		panic("Unsupported format")
	}
}

// Преобразование в полутоновое изображение
func grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	gray := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			gray.SetGray(x, y, c)
		}
	}
	return gray
}

// Гистограмма яркости
func computeHistogram(img *image.Gray) [256]int {
	var hist [256]int
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			hist[img.GrayAt(x, y).Y]++
		}
	}
	return hist
}

// Регулировка яркости
func adjustBrightness(img *image.Gray, b int) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := int(img.GrayAt(x, y).Y) + b
			if val < 0 {
				val = 0
			} else if val > 255 {
				val = 255
			}
			result.SetGray(x, y, color.Gray{uint8(val)})
		}
	}
	return result
}

// Негативное преобразование
func negative(img *image.Gray) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := 255 - img.GrayAt(x, y).Y
			result.SetGray(x, y, color.Gray{val})
		}
	}
	return result
}

// Бинаризация
func binarize(img *image.Gray, threshold uint8) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := img.GrayAt(x, y).Y
			if val < threshold {
				result.SetGray(x, y, color.Gray{0})
			} else {
				result.SetGray(x, y, color.Gray{255})
			}
		}
	}
	return result
}

// Линейное растяжение гистограммы
func contrastStretching(img *image.Gray) *image.Gray {
	hist := computeHistogram(img)
	q1, q2 := findQ1Q2(hist)
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	scale := 255.0 / float64(q2-q1)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := float64(img.GrayAt(x, y).Y)
			newVal := (val - float64(q1)) * scale
			if newVal < 0 {
				newVal = 0
			} else if newVal > 255 {
				newVal = 255
			}
			result.SetGray(x, y, color.Gray{uint8(newVal)})
		}
	}
	return result
}

// Нахождение Q1 и Q2 (первые не нулевые значения в гистограмме)
func findQ1Q2(hist [256]int) (int, int) {
	q1, q2 := 0, 255
	for i, v := range hist {
		if v > 0 {
			q1 = i
			break
		}
	}
	for i := 255; i >= 0; i-- {
		if hist[i] > 0 {
			q2 = i
			break
		}
	}
	return q1, q2
}

// Гамма-коррекция
func gammaCorrection(img *image.Gray, gamma float64) *image.Gray {
	bounds := img.Bounds()
	result := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			val := float64(img.GrayAt(x, y).Y) / 255.0
			corrected := math.Pow(val, gamma) * 255.0
			result.SetGray(x, y, color.Gray{uint8(corrected)})
		}
	}
	return result
}

func main() {
	// Пример использования
	img, err := loadImage("input.jpg")
	if err != nil {
		panic(err)
	}

	grayImg := grayscale(img)
	saveImage("gray.jpg", grayImg)

	brightImg := adjustBrightness(grayImg, 50)
	saveImage("bright.jpg", brightImg)

	negImg := negative(grayImg)
	saveImage("negative.jpg", negImg)

	binImg := binarize(grayImg, 128)
	saveImage("binary.jpg", binImg)

	contrastImg := contrastStretching(grayImg)
	saveImage("contrast.jpg", contrastImg)

	gammaImg := gammaCorrection(grayImg, 0.5)
	saveImage("gamma.jpg", gammaImg)
}
