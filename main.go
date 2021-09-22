package main

import (
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/deepakandgupta/skalappy-bird/objects/background"
	"github.com/deepakandgupta/skalappy-bird/objects/basebackground"
	"github.com/deepakandgupta/skalappy-bird/objects/bird"
	"github.com/deepakandgupta/skalappy-bird/objects/gamefont"
	"github.com/deepakandgupta/skalappy-bird/objects/pipe"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
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

// Using make had some strange effect check why
var pipeCoords [][2]float32
const distBwPipe = 250

var backBaseHeight int

const pipeIndexForCollision = 0;

var myGameFont font.Face

const startGameText = "Press 'Space' to start\n 'Mouse Left' to play"
const gameOverText = "             Game Over\n Press 'Space' to restart"

var isGameStart bool

type Game struct{
	spacePressed bool
	mousePressed bool

	yBird float32
}

func init() {
	isGameStart = true
	
	initialiseObjects()

	initPipes()

}

func (g *Game) Update() error {
	g.spacePressed =  inpututil.IsKeyJustPressed(ebiten.KeySpace)
	g.mousePressed =  inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)

	movePipes()
	removeAndSpawnPipe();
	incremenetScore()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	background.Draw(screen)
	for _, coord  := range pipeCoords {
		pipe.Draw(coord[0], coord[1], pipeGap, screen)
	}

	basebackground.Draw( float64(screenHeight), screen)
	g.applyGravityAndVelOnBird(screen)
	g.checkBirdCollision(screen);
	drawScore(screen);

	if(!isBirdAlive && g.spacePressed){
		g.yBird = screenWidth/2
		restartGame()
		} else if(!isBirdAlive && isGameStart) {
			text.Draw(screen, startGameText, myGameFont, 100, screenHeight/2, color.White)
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

func initPipes(){
	pipeCoords = pipeCoords[:0]
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

func movePipes(){
	for i := 0; i < len(pipeCoords); i++{
		pipeCoords[i][0] -= pipeSpeed
	}
}

func restartGame(){
	isGameStart = false
	score = 0

	initPipes()

	gravity = 2.3
	upVelocity  = 50
	isBirdAlive = true
	pipeSpeed = 1.5
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
	if(g.yBird >= float32(birdLowY - birdHeight)){
		g.yBird = float32(birdLowY - birdHeight)
		setGameOver(screen)
	} else if(g.yBird <=0){
		g.yBird = 0;
		setGameOver(screen)
	}
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
	text.Draw(screen, gameOverText, myGameFont, 100, screenHeight/2, color.White)
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
	bird.Draw(float64(xBird), g.yBird, screen)
	g.yBird += gravity;
	if g.mousePressed {
		g.yBird -= upVelocity
	}
}

func drawScore(screen *ebiten.Image){
	text.Draw(screen, "Score: "+fmt.Sprint(score), myGameFont, screenWidth -180, 30, color.White)
}

// Simple Random function to give random between range
func randRange(min, max int) int{
	// Without seed, GO return similar random values
	rand.Seed(time.Now().UnixNano())
    return rand.Intn(max - min + 1) + min
}


// Images that will be used in game, same images can be used multiple times
func initialiseObjects(){
	background.Init();
	myGameFont = gamefont.Init()
	birdWidth, birdHeight = bird.Init()
	pipeWidth = pipe.Init();
	backBaseHeight = basebackground.Init()
}

