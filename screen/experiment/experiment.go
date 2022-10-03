package experiment

import (
	"ebiten2arkham/config"
	"ebiten2arkham/input"
	"ebiten2arkham/renderer"
	"ebiten2arkham/screen"
	"fmt"
	"github.com/HaBaLeS/arkham-go/runtime"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"math"
	"os"
)

var useSingleImage = true

type ExperimentalScreen struct {
	//01552
	//01005 a,b
	//01111 a,b

	//data/textures/bg168.png
	//data/textures/bg100.png
	bgs           []*ebiten.Image
	tiles         []*BgTile
	card          *renderer.CardSprite
	singleBGImage *ebiten.Image

	camera *renderer.Camera
	filter ebiten.Filter

	worldImage *ebiten.Image
	gamepads   []ebiten.GamepadID
	stroke     *input.Stroke
}

type BgTile struct {
	posx    float64
	posy    float64
	texture *ebiten.Image
}

func NewScreen() screen.Screen {

	rows := 20
	cols := 20

	c1 := renderer.NewCardSprite(runtime.CardDBG().GetCard("01552"))
	c1.Enable()
	c1.X = 150
	c1.Y = 150

	ex := &ExperimentalScreen{
		bgs:           make([]*ebiten.Image, 2),
		worldImage:    ebiten.NewImage(cols*200, rows*200),
		card:          c1,
		singleBGImage: renderer.LoadImage("bg/pexels-pixabay-235985_ar_fixed.jpg"),
		filter:        ebiten.FilterLinear,
		stroke:        &input.Stroke{},
	}

	if useSingleImage {
		x, y := ex.singleBGImage.Size()
		ex.worldImage = ebiten.NewImage(x*2, y*2)
	} else {
		ex.worldImage = ebiten.NewImage(cols*200, rows*200)
	}

	ex.camera = renderer.NewCamera(ex.worldImage)

	ex.bgs[0] = renderer.LoadImage("textures/bg28.png")
	ex.bgs[1] = renderer.LoadImage("textures/bg28.png")

	ex.tiles = make([]*BgTile, rows*cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			pos := j + rows*i
			t := &BgTile{
				posx: float64(200 * j),
				posy: float64(200 * i),
			}
			if pos%2 == 0 && i%2 == 0 {
				t.texture = ex.bgs[0]
			} else if i%2 != 0 && pos%2 != 0 {
				t.texture = ex.bgs[0]
			} else {
				t.texture = ex.bgs[1]
			}
			ex.tiles[pos] = t

			//fmt.Printf("%f:%f\n", t.posx, t.posy)
		}
	}

	return ex
}

func (ex *ExperimentalScreen) Draw(screen *ebiten.Image) {
	ex.worldImage.Clear()

	if useSingleImage {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(2, 2)
		op.Filter = ex.filter
		ex.worldImage.DrawImage(ex.singleBGImage, op)
	} else {
		for _, t := range ex.tiles {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(t.posx, t.posy)
			op.Filter = ex.filter
			ex.worldImage.DrawImage(t.texture, op)
		}
	}

	ex.card.Draw(ex.worldImage)

	ex.camera.Draw(screen)

	worldX, worldY := ex.camera.ScreenToWorld(ebiten.CursorPosition())
	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf("TPS: %0.2f\nMove (WASD/Arrows)\nZoom (QE)\nRotate (R)\nReset (Space)", ebiten.ActualTPS()),
	)

	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("\nCursor World Pos: %.2f,%.2f", worldX, worldY), 0, int(config.Cfg.CameraHeight-32.0),
	)
	ebitenutil.DebugPrintAt(
		screen,
		fmt.Sprintf("\nCamera Pos: %.2f,%.2f - ZF %.2f", ex.camera.Xoff, ex.camera.Yoff, ex.camera.ZoomFactor), 0, int(config.Cfg.CameraHeight-48),
	)
}

func (ex *ExperimentalScreen) Update() error {
	x, y := ex.camera.ScreenToWorld(ebiten.CursorPosition())
	if ex.card.Contains(float64(x), float64(y)) {
		ex.card.Greyout = true
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			ex.stroke.Start()
		}
		if ex.stroke.Active {
			//IDEA ... use a stroke with world coordinates not with screen/camera
			xx, yy := ex.stroke.Delta()
			zf := math.Pow(1.01, float64(ex.camera.ZoomFactor))
			ex.card.X -= float64(xx) * math.Abs(zf)
			ex.card.Y -= float64(yy) * math.Abs(zf)
		}
	} else {
		ex.card.Greyout = false
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		ex.stroke.Stop()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		ex.camera.ZoomFactor -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		ex.camera.ZoomFactor += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyF) {
		ex.filter = ebiten.FilterNearest
	}
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		ex.filter = ebiten.FilterLinear
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		ex.camera.Yoff -= 1 * config.Cfg.CameraMoveFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		ex.camera.Yoff += 1 * config.Cfg.CameraMoveFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		ex.camera.Xoff -= 1 * config.Cfg.CameraMoveFactor
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		ex.camera.Xoff += 1 * config.Cfg.CameraMoveFactor
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		ex.camera.Rotation += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyT) {
		ex.camera.Rotation -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		ex.camera.Reset()
	}

	_, dwx := ebiten.Wheel()
	ex.camera.ZoomFactor += dwx * 5

	ex.handleGamePad()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonMiddle) {
		ex.stroke.Start()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle) {
		if ex.stroke.Active {
			x, y := ex.stroke.Delta()
			ex.camera.Xoff += float64(x)
			ex.camera.Yoff += float64(y)
		}
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonMiddle) {
		ex.stroke.Stop()
	}

	driftx, drifty := ex.camera.ScreenToWorld(0, 0)
	if driftx < 0 {
		ex.camera.Xoff -= driftx
	}

	if drifty < 0 {
		ex.camera.Yoff -= drifty
	}

	return nil
}

func (ex *ExperimentalScreen) Resume() {

}

func (ex *ExperimentalScreen) Pause() {

}

func (ex *ExperimentalScreen) handleGamePad() {
	ex.gamepads = make([]ebiten.GamepadID, 0)
	ex.gamepads = ebiten.AppendGamepadIDs(ex.gamepads)
	if len(ex.gamepads) > 0 {
		ex.camera.ZoomFactor -= ebiten.StandardGamepadButtonValue(ex.gamepads[0], ebiten.StandardGamepadButtonFrontBottomLeft)
		ex.camera.ZoomFactor += ebiten.StandardGamepadButtonValue(ex.gamepads[0], ebiten.StandardGamepadButtonFrontBottomRight)

		//fixme add a threshold to avoid drifrt
		ex.camera.Xoff += ebiten.StandardGamepadAxisValue(ex.gamepads[0], ebiten.StandardGamepadAxisLeftStickHorizontal) * 5
		ex.camera.Yoff += ebiten.StandardGamepadAxisValue(ex.gamepads[0], ebiten.StandardGamepadAxisLeftStickVertical) * 5

		if ebiten.IsStandardGamepadButtonPressed(ex.gamepads[0], ebiten.StandardGamepadButtonRightRight) {
			os.Exit(0)
		}
	}
}
