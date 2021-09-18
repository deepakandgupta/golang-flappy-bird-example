package main

import (
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const screenWidth = 570
const screenHeight = 512

// Variable references for all images
var backgroundImg *ebiten.Image
var birdImg *ebiten.Image
var pipeImg *ebiten.Image
var pipeDownImg *ebiten.Image
var backBaseImg *ebiten.Image

const xBird = 100
const birdLowY = screenHeight;
var birdWidth int
var birdHeight int
// Forces on the Bird
var gravity float32 = 2.2
var upVelocity float32 = 60

const pipeNum = 6;
// Pipe Variables
const pipeGap int = 150
var pipeSpeed = 1;

var pipeWidth int
var pipeHeight int

// Using make had some strange effect check why
var pipeCoords [][2]int
const distBwPipe = 150

var backBaseHeight int

func init() {
	initialiseImages()

	birdWidth = birdImg.Bounds().Max.X
	birdHeight = birdImg.Bounds().Max.Y
	
	pipeWidth = pipeImg.Bounds().Max.X
	pipeHeight = pipeImg.Bounds().Max.Y

	backBaseHeight = backBaseImg.Bounds().Max.Y
	// gerating random coordinates for pipes
	for i := 0; i < pipeNum; i++{
		yCoord := randRange(100, screenHeight - pipeGap)
		xCoord := distBwPipe*(i+1)
		pc := [2]int{xCoord, yCoord}
		pipeCoords = append(pipeCoords, pc)
	}
}

type Game struct{
	spacePressed bool
	mousePressed bool

	yBird float32
}

func (g *Game) Update() error {
	g.spacePressed =  inpututil.IsKeyJustPressed(ebiten.KeySpace)
	g.mousePressed =  inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	for i := 0; i < len(pipeCoords); i++{
		pipeCoords[i][0] -= pipeSpeed
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen)
	for _, coord  := range pipeCoords {
		drawPipe(coord[0], coord[1], screen)
	}

	drawBackBase(screen)

	g.applyGravityAndVelOnBird(screen)
	removeAndSpawnPipe();
	// ebitenutil.DebugPrintAt(screen, "number", pipeCoords[1][0], pipeCoords[1][1])
	g.checkBirdCollision(screen);

}

func printmsg(screen *ebiten.Image, msg string){
	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight + backBaseHeight)
	ebiten.SetWindowTitle("Skallapy bird!")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) checkBirdCollision(screen *ebiten.Image){
	// pipeLeft :=  pipeCoords[1][0]
	// pipeTop := pipeCoords[1][1]

	// pipeBottom := pipeTop + pipeGap
	// ebitenutil.DebugPrintAt(screen, "this", pipeLeft, pipeBottom)	
	if(checkXCoordPipeBird() && !checkYCoordPipeBird(g.yBird)){
		printmsg(screen, "Not Safe")
		// upVelocity = 0
		gravity = 0
		pipeSpeed = 0	
	}
}

func checkXCoordPipeBird() bool{
	birdLeft :=  float32(xBird)
	birdRight := birdLeft + float32(birdWidth)
	
	pipeLeft :=  float32(pipeCoords[1][0])
	pipeRight := pipeLeft + float32(pipeWidth)

	return ((pipeLeft <= birdLeft || pipeLeft <= birdRight) &&
		(pipeRight >= birdLeft || pipeRight >= birdRight) )
}


func checkYCoordPipeBird(yBird float32) bool{
	
	birdTop := yBird	
	birdBottom := birdTop + float32(birdHeight)

	pipeTop := float32(pipeCoords[1][1])
	pipeBottom := pipeTop + float32(pipeGap)

	return birdTop >= pipeTop && birdBottom <= pipeBottom
}

func removeAndSpawnPipe(){
	if(pipeCoords[0][0] < -distBwPipe){
		// Adding a new Pipe
		yCoord := randRange(100, screenHeight - pipeGap)
		xCoord := distBwPipe + pipeCoords[(len(pipeCoords) -1)][0]
		anotherPipe := [2]int{xCoord, yCoord}
		pipeCoords =  append(pipeCoords, anotherPipe)
		// Removing the last pipe
		pipeCoords = pipeCoords[1:]
	}
}

func (g *Game) applyGravityAndVelOnBird(screen *ebiten.Image){
	g.drawBird(screen)
	g.yBird += gravity;
	if g.mousePressed || g.spacePressed {
		g.yBird -= upVelocity
	}
}

func drawBackground(screen *ebiten.Image){
	// backOp.GeoM.Scale(1, 1)

	backOp := &ebiten.DrawImageOptions{}
	backOp.GeoM.Translate(0, 0)
	screen.DrawImage(backgroundImg, backOp)

	var backWidth = backgroundImg.Bounds().Max.X
	
	backOp2 := &ebiten.DrawImageOptions{}
	backOp2.GeoM.Translate(float64(backWidth), 0)
	screen.DrawImage(backgroundImg, backOp2)
}

func drawBackBase(screen *ebiten.Image){
	backBaseOp := &ebiten.DrawImageOptions{}
	backBaseOp.GeoM.Translate(0, float64(screenHeight))
	screen.DrawImage(backBaseImg, backBaseOp)
}



func (g *Game) drawBird(screen *ebiten.Image){
	// Bird options
	birdOp := &ebiten.DrawImageOptions{}
	// Bouding birds on screen up and down
	if(g.yBird >= float32(birdLowY - birdHeight)){
		g.yBird = float32(birdLowY - birdHeight)
	} else if(g.yBird <=0){
		g.yBird = 0;
	}
	birdOp.GeoM.Translate(xBird, float64(g.yBird))
	screen.DrawImage(birdImg, birdOp)	
}

func drawPipe(xCoord int, yCoord int, screen *ebiten.Image){
	// fmt.Println(xCoord)

	pipeOp := &ebiten.DrawImageOptions{}
	pipeOp.GeoM.Translate(float64(xCoord), float64(yCoord - pipeHeight))
	screen.DrawImage(pipeDownImg, pipeOp)	
	
	pipe2Op := &ebiten.DrawImageOptions{}
	pipe2Op.GeoM.Translate(float64(xCoord), float64(yCoord + pipeGap))
	screen.DrawImage(pipeImg, pipe2Op)	
}

func randRange(min, max int) int{
	rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min + 1) + min
}

func initialiseImages(){
	var backErr error
	backgroundImg, _, backErr = ebitenutil.NewImageFromFile("./img/background.png")
	if backErr != nil {
		log.Fatal(backErr)
	}

	var birdErr error
	birdImg, _, birdErr = ebitenutil.NewImageFromFile("./img/bird.png")
	if birdErr != nil {
		log.Fatal(birdErr)
	}

	var pipeErr error
	pipeImg, _, pipeErr = ebitenutil.NewImageFromFile("./img/pipe.png")
	if pipeErr != nil {
		log.Fatal(pipeErr)
	}

	var pipeDownErr error
	pipeDownImg, _, pipeDownErr = ebitenutil.NewImageFromFile("./img/pipe-down.png")
	if pipeDownErr != nil {
		log.Fatal(pipeDownErr)
	}

	var backBaseErr error
	backBaseImg, _, backBaseErr = ebitenutil.NewImageFromFile("./img/back-base.png")
	if backBaseErr != nil {
		log.Fatal(backBaseErr)
	}
}