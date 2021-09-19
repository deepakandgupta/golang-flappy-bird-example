package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Forces on the Bird
var gravity float32 = 0
var upVelocity float32 = 0

var isBirdAlive = false
var score = 0

const screenWidth = 570
const screenHeight = 512

const xBird = 100
const birdLowY = screenHeight;
var birdWidth int
var birdHeight int

const birdExtraBuffer = 3


const pipeNum = 6;
// Pipe Variables
const pipeGap int = 100
var pipeSpeed float32 = 0;

var pipeWidth int
// var pipeHeight int

var pipeDownHeight int

// Using make had some strange effect check why
var pipeCoords [][2]float32
const distBwPipe = 250

var backBaseHeight int

const pipeIndexForCollision = 0;


// Variable references for all images
var backgroundImg *ebiten.Image
var birdImg *ebiten.Image
var pipeImg *ebiten.Image
var pipeDownImg *ebiten.Image
var backBaseImg *ebiten.Image


var gameFont font.Face

const startGameText = "Press 'Space' to start\n 'Mouse Left' to play"
const gameOverText = "             Game Over\n Press 'Space' to restart"

var isGameStart bool

func init() {
	isGameStart = true

	initialiseImages()
	initialiseFont()

	birdWidth = birdImg.Bounds().Max.X
	birdHeight = birdImg.Bounds().Max.Y
	
	pipeWidth = pipeImg.Bounds().Max.X
	// pipeHeight = pipeImg.Bounds().Max.Y

	pipeDownHeight = pipeDownImg.Bounds().Max.Y

	backBaseHeight = backBaseImg.Bounds().Max.Y

	pipeCoords = pipeCoords[:0]
	score = 0
	// gerating random coordinates for pipes
	for i := 0; i < pipeNum; i++{
		yCoord := randRange(100, screenHeight - pipeGap)
		var xCoord float32
		if(i==0){
			xCoord = 500
		} else{
			xCoord = pipeCoords[i-1][0] + distBwPipe
		}
		pc := [2]float32{xCoord, float32(yCoord)}
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

func restartGame(){
	isGameStart = false

	pipeCoords = pipeCoords[:0]
	score = 0
	// gerating random coordinates for pipes
	for i := 0; i < pipeNum; i++{
		yCoord := randRange(100, screenHeight - pipeGap)
		var xCoord float32
		if(i==0){
			xCoord = 500
		} else{
			xCoord = pipeCoords[i-1][0] + distBwPipe
		}
		pc := [2]float32{xCoord, float32(yCoord)}
		pipeCoords = append(pipeCoords, pc)
	}

	gravity = 2.3
	upVelocity  = 50
	isBirdAlive = true
	pipeSpeed = 1.5
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawBackground(screen)
	for _, coord  := range pipeCoords {
		drawPipe(coord[0], coord[1], screen)
	}

	drawBackBase(screen)

	g.applyGravityAndVelOnBird(screen)
	removeAndSpawnPipe();
	g.checkBirdCollision(screen);
	incremenetScore()
	drawScore(screen);

	if(!isBirdAlive && g.spacePressed){
		g.yBird = screenWidth/2
		restartGame()
		} else if(!isBirdAlive && isGameStart) {
			text.Draw(screen, startGameText, gameFont, 100, screenHeight/2, color.White)
			g.yBird = screenWidth/2
	}

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

func incremenetScore(){
	birdLeft :=  float32(xBird)
	
	pipeLeft :=  float32(pipeCoords[pipeIndexForCollision][0])
	pipeRight := pipeLeft + float32(pipeWidth)

	if(isBirdAlive && birdLeft == float32(math.Floor(float64(pipeRight))) ||
	birdLeft == float32(math.Ceil(float64(pipeRight)))){
		score++
	}
}

// We only need to check the first pipe for collision
// The bird does not move left/right
func (g *Game) checkBirdCollision(screen *ebiten.Image){
	if(checkXCoordPipeBird() && !checkYCoordPipeBird(g.yBird)){
		setGameOver(screen)
	}
}

// Check if bird x and pipe's x coordinate are matched
func checkXCoordPipeBird() bool{
	birdLeft :=  float32(xBird) + birdExtraBuffer
	birdRight := birdLeft + float32(birdWidth) - 2*birdExtraBuffer
	
	pipeLeft :=  float32(pipeCoords[pipeIndexForCollision][0])
	pipeRight := pipeLeft + float32(pipeWidth)

	return ((pipeLeft <= birdLeft || pipeLeft <= birdRight) &&
		(pipeRight >= birdLeft || pipeRight >= birdRight) )
}

// Check if bird y and pipe's y coordinate are matched
func checkYCoordPipeBird(yBird float32) bool{
	birdTop := yBird + birdExtraBuffer
	birdBottom := birdTop + float32(birdHeight) - 2*birdExtraBuffer

	pipeTop := float32(pipeCoords[pipeIndexForCollision][1])
	pipeBottom := pipeTop + float32(pipeGap)

	return birdTop >= pipeTop && birdBottom <= pipeBottom
}

// On game over, remove all control but still draw the current state
func setGameOver(screen *ebiten.Image){
	text.Draw(screen, gameOverText, gameFont, 100, screenHeight/2, color.White)
	gravity = 0
	pipeSpeed = 0	
	upVelocity = 0
	isBirdAlive = false
}

func removeAndSpawnPipe(){
	// If the pipe is off screen, remove that pipe
	if(pipeCoords[0][0] < -float32(pipeWidth)){
		
		// Adding a new Pipe
		yCoord := randRange(100, screenHeight - pipeGap)
		xCoord := distBwPipe + pipeCoords[(len(pipeCoords) -1)][0]
		anotherPipe := [2]float32{xCoord, float32(yCoord)}
		pipeCoords =  append(pipeCoords, anotherPipe)
		
		// Removing first pipe
		pipeCoords = pipeCoords[1:]
	}
}

func (g *Game) applyGravityAndVelOnBird(screen *ebiten.Image){
	g.drawBird(screen)
	g.yBird += gravity;
	if g.mousePressed {
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

func drawScore(screen *ebiten.Image){
	text.Draw(screen, "Score: "+fmt.Sprint(score), gameFont, screenWidth -180, 30, color.White)
}

func (g *Game) drawBird(screen *ebiten.Image){
	// Bird options
	birdOp := &ebiten.DrawImageOptions{}
	// Bounding birds on screen up and down
	// If bird touches either upper or lower boundary -  game over
	if(g.yBird >= float32(birdLowY - birdHeight)){
		g.yBird = float32(birdLowY - birdHeight)
		setGameOver(screen)
	} else if(g.yBird <=0){
		g.yBird = 0;
		setGameOver(screen)
	}
	birdOp.GeoM.Translate(xBird, float64(g.yBird))
	screen.DrawImage(birdImg, birdOp)	
}

func drawPipe(xCoord float32, yCoord float32, screen *ebiten.Image){
	// fmt.Println(xCoord)

	pipeOp := &ebiten.DrawImageOptions{}
	pipeOp.GeoM.Translate(float64(xCoord), float64(yCoord - float32(pipeDownHeight)))
	screen.DrawImage(pipeDownImg, pipeOp)	
	
	pipe2Op := &ebiten.DrawImageOptions{}
	pipe2Op.GeoM.Translate(float64(xCoord), float64(yCoord + float32(pipeGap)))
	screen.DrawImage(pipeImg, pipe2Op)	
}

// Simple Random function to give random between range
func randRange(min, max int) int{
	// Without seed, GO return similar random values
	rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min + 1) + min
}


// Images that will be used in game, same images can be used multiple times
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

// Adding font to use in game
func initialiseFont(){
	fontBytes, err := ioutil.ReadFile("./fonts/yoster.ttf")
	if err != nil {
		fmt.Println(err)
		return
	}


	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 60
	gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
	})
	if err != nil {
		log.Fatal(err)
	}
}
