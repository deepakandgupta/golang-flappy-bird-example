package main

import (
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var backgroundImg *ebiten.Image

var birdImg *ebiten.Image

func init() {
	var berr error
	backgroundImg, _, berr = ebitenutil.NewImageFromFile("./img/background.png")
	if berr != nil {
		log.Fatal(berr)
	}

	var birdErr error
	birdImg, _, birdErr = ebitenutil.NewImageFromFile("./img/bird.png")
	if birdErr != nil {
		log.Fatal(birdErr)
	}

}

type Game struct{
	keys []ebiten.Key
	mousePressed bool

	yBird float64
}

func (g *Game) Update() error {
	g.keys = inpututil.PressedKeys();
	// inpututil.IsKeyJustPressed(ebiten.KeySpace)
	g.mousePressed =  inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	g.yBird += 2;
	if g.mousePressed {
		g.yBird -= 60
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Backgorund Image
	backOp := &ebiten.DrawImageOptions{}
	backOp.GeoM.Translate(300, 50)
	backOp.GeoM.Scale(1, 1)
	screen.DrawImage(backgroundImg, backOp)
	
	// Bird options
	birdOp := &ebiten.DrawImageOptions{}
	birdOp.GeoM.Translate(300, g.yBird)
	// birdOp.GeoM.Rotate(1);
	screen.DrawImage(birdImg, birdOp)

	keyStrs := []string{}
	for _, p := range g.keys {
		keyStrs = append(keyStrs, p.String())
		// if(p.String() == "Space"){
		// 	ebitenutil.DebugPrint(screen, "Jump Up")
		// }
	}
	
	ebitenutil.DebugPrint(screen, strings.Join(keyStrs, ", "))
	
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1080, 720
}

func main() {
	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Skallapy bird!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}