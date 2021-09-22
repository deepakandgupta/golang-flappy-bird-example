package gamefont

import (
	"fmt"
	"io/ioutil"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const fontPath = "./fonts/yoster.ttf"

func Init() font.Face{
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}


	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 60
	gameFont, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
	})
	if err != nil {
		log.Fatal(err)
	}

	return gameFont;
}