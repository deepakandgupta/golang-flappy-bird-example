package pipe

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const pipeUpImgPath = "./img/pipe.png"
const pipeDownImgPath = "./img/pipe-down.png"

var pipeUpImg *ebiten.Image
var pipeDownImg *ebiten.Image

func Init() int{
	var pipeErr error
	pipeUpImg, _, pipeErr = ebitenutil.NewImageFromFile(pipeUpImgPath)
	if pipeErr != nil {
		log.Fatal(pipeErr)
	}
	pipeWidth := pipeUpImg.Bounds().Max.X

	var pipeDownErr error
	pipeDownImg, _, pipeDownErr = ebitenutil.NewImageFromFile(pipeDownImgPath)
	if pipeDownErr != nil {
		log.Fatal(pipeDownErr)
	}

	return pipeWidth;
}


func Draw(x float32, y float32, pipeGap int,  screen *ebiten.Image){
	pipeUpOp := &ebiten.DrawImageOptions{}
	pipeUpOp.GeoM.Translate(float64(x), float64(y + float32(pipeGap)))
	screen.DrawImage(pipeUpImg, pipeUpOp)	
	
	
	pipeDownHeight := pipeDownImg.Bounds().Max.Y
	pipeDownOp := &ebiten.DrawImageOptions{}
	pipeDownOp.GeoM.Translate(float64(x), float64(y - float32(pipeDownHeight)))
	screen.DrawImage(pipeDownImg, pipeDownOp)	
	
}