package lib

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func TempFilename(prefix string, extension string) string {
	randBytes := make([]byte, 16)
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+extension)
}

func GenerateQR(uuid, input string) (string, error) {
	//log.Println("Original data:", base64)
	code, err := qr.Encode(input, qr.L, qr.Unicode)
	if err != nil {
		return "", err
	}
	//log.Println("Encoded data: ", code.Content())

	if input != code.Content() {
		return "", errors.New("data differs")
	}

	code, err = barcode.Scale(code, 60, 60)
	if err != nil {
		return "", err
	}
	printPng(os.Stdout, code)

	code, err = barcode.Scale(code, 600, 600)
	if err != nil {
		return "", err
	}
	output := TempFilename(uuid, ".png")
	writePng(output, code)

	//log.Println(`Now open test.png and scan QR code, it will be: "IAV19ysYSl0HUuG5QiCDvdHkowqdGXb0HbqUAWUzHw==" instead of "IAV19ysYSl0HUuG5QiCDvdHkowqdGXb0HbaUAWUzHw==" ('a' is changed to 'q' in 'aUAWUzHw==' part)`)
	return output, nil
}

const BLACK = "\033[40m  \033[0m"
const WHITE = "\033[47m  \033[0m"

func printPng(writer io.Writer, img image.Image) error {

	// Create a new grayscale image
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	for y := 11; y < h-14; y++ {
		for x := 12; x < w-13; x++ {
			if img.At(x, y) == color.White && img.At(x, y+1) == color.White {
				writer.Write([]byte(WHITE))
			} else if img.At(x, y) == color.White && img.At(x, y+1) == color.Black {
				writer.Write([]byte(BLACK))
			} else if img.At(x, y) == color.Black && img.At(x, y+1) == color.White {
				writer.Write([]byte(WHITE))
			} else {
				writer.Write([]byte(BLACK))
			}
		}
		writer.Write([]byte("\n"))
	}

	return nil
}

func writePng(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
