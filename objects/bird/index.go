package bird

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const imgPath = "./img/bird.png"
var img *ebiten.Image

func Init() (int, int){
	var birdErr error
	img, _, birdErr = ebitenutil.NewImageFromFile(imgPath)
	if birdErr != nil {
		log.Fatal(birdErr)
	}

	birdWidth := img.Bounds().Max.X
	birdHeight := img.Bounds().Max.Y
	return birdWidth, birdHeight
}

func Draw(x float64, y float32,screen *ebiten.Image){
	birdOp := &ebiten.DrawImageOptions{}
	birdOp.GeoM.Translate(x, float64(y))
	screen.DrawImage(img, birdOp)	
}