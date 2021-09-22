package basebackground

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const imgPath = "./img/back-base.png"
var img *ebiten.Image

func Init() int{
	var backBaseErr error
	img, _, backBaseErr = ebitenutil.NewImageFromFile(imgPath)
	if backBaseErr != nil {
		log.Fatal(backBaseErr)
	}

	return img.Bounds().Max.Y
}

func Draw(y float64,screen *ebiten.Image){
	backBaseOp := &ebiten.DrawImageOptions{}
	backBaseOp.GeoM.Translate(0, y)
	screen.DrawImage(img, backBaseOp)	
}