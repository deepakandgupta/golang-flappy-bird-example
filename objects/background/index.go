package background

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const imgPath = "./img/background.png"
var img *ebiten.Image

func Init(){
	var backErr error
	img, _, backErr = ebitenutil.NewImageFromFile(imgPath)
	if backErr != nil {
		log.Fatal(backErr)
	}
}

func Draw(screen *ebiten.Image){
	// backOp.GeoM.Scale(1, 1)

	backOp := &ebiten.DrawImageOptions{}
	backOp.GeoM.Translate(0, 0)
	screen.DrawImage(img, backOp)

	var backWidth = img.Bounds().Max.X
	
	backOp2 := &ebiten.DrawImageOptions{}
	backOp2.GeoM.Translate(float64(backWidth), 0)
	screen.DrawImage(img, backOp2)
}