package driver

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mcsymiv/godriver/config"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.NewSource(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newScreenShotCommand() Command {
	return Command{
		Path:   PathDriverScreenshot,
		Method: http.MethodGet,
	}
}

func screenshot(d Driver) error {

	data := new(struct{ Value string })
	d.execute(defaultStrategy{Command{
		Path:         PathDriverScreenshot,
		Method:       http.MethodGet,
		ResponseData: data,
	}})

	decodedImage, err := base64.StdEncoding.DecodeString(data.Value)
	if err != nil {
		return fmt.Errorf("error on decode base64 string: %v", err)
	}

	// Create an image.Image from decoded bytes
	img, err := png.Decode(strings.NewReader(string(decodedImage)))
	if err != nil {
		return fmt.Errorf("error on decode: %v", err)
	}

	// Create a new file for the output JPEG image
	// TODO: upd randSeq, use meaninful screenshot name
	outputFile, err := os.Create(fmt.Sprintf("%s/%s_%s.jpg", config.TestSetting.ArtifactScreenshotsPath, randSeq(8), time.Now().Format("2006_01_02_15:04:05")))
	if err != nil {
		return fmt.Errorf("error on create file: %v", err)
	}
	defer outputFile.Close()

	// Encode the image as JPEG
	err = jpeg.Encode(outputFile, img, nil)
	if err != nil {
		return fmt.Errorf("error on encode: %v", err)
	}

	return nil
}

func (d Driver) Screenshot() {
	err := screenshot(d)
	if err != nil {
		log.Println("error on screenshot:", err)
	}
}
