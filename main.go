package main

import (
	"NotaborEngine/notacolor"
	"NotaborEngine/notacore"
	"NotaborEngine/notaentity"
	"NotaborEngine/notamath"
	"NotaborEngine/notasdl"
	"NotaborEngine/notatask"
	"NotaborEngine/notatomic"
	"NotaborEngine/notaui"
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

func main() {
	// Engine setup
	settings := &notacore.Settings{
		Vsync:      true,
		SoundLevel: 1,
		Muted:      false,
	}

	drawingLoop := notatask.CreateLoop(60)

	engine, err := notacore.CreateEngine(settings)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	cfg := &notasdl.WindowConfig{
		X:         50,
		Y:         50,
		W:         800,
		H:         600,
		Title:     "Entity Test",
		Type:      notasdl.Windowed,
		Resizable: true,
		TargetFPS: 60,
		Loops:     []*notatask.Loop{drawingLoop},
	}

	win, err := engine.CreateWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	em := engine.EntityManager
	circleRadius := float32(0.25)

	ballVisual, err := win.LoadVisual("quadSprite", "resources/images/hahaha.jpg", notasdl.VisualOptions{
		Width:        circleRadius * 2,
		Height:       circleRadius * 2,
		Mask:         notasdl.MaskCircle,
		CircleRadius: 0.5,
		CircleEdge:   0.01,
	})

	if err != nil {
		log.Fatal(err)
	}

	entity := em.CreateEntity("quad").
		WithVisual(ballVisual).
		WithCollision(notaentity.CircleCollision(circleRadius)).
		WithColor(notacolor.White)

	moveStep := float32(0.05)
	var moveSpeed notatomic.Float32
	moveSpeed.Set(moveStep)

	ui, err := notaui.New(engine, win)
	if err != nil {
		log.Fatal(err)
	}

	var clicks atomic.Int32
	var colorChoice atomic.Int32
	playerName := "nota"
	hudEnabled := true

	ui.Panel("hud-panel").Rect(14, 14, 314, 220)
	ui.Text("hud-title", "NOTA UI").At(24, 24).Scale(2).Color(notacolor.Cyan)
	ui.TextFunc("hud-info", func() string {
		state := "OFF"
		if hudEnabled {
			state = "ON"
		}
		return fmt.Sprintf("PLAYER %s  CLICKS %d  HUD %s", playerName, clicks.Load(), state)
	}).At(24, 48)
	ui.Button("hud-click", "CLICK").Rect(24, 72, 92, 28).OnClick(func() {
		clicks.Add(1)
	})
	ui.Input("hud-name", &playerName).Rect(124, 72, 184, 28).Placeholder("name")
	ui.Slider("hud-speed", &moveStep, 0.01, 0.12).Rect(24, 114, 284, 34).Label("SPEED").OnChange(func(v float32) {
		moveSpeed.Set(v)
	})
	ui.Checkbox("hud-toggle", "HUD ENABLED", &hudEnabled).Rect(24, 154, 160, 24)

	grid := ui.Grid("hud-color-grid", notaui.R(24, 188, 284, 34), 3, 1).Gap(8)
	grid.Button("hud-white", "WHITE", 0, 0).OnClick(func() {
		colorChoice.Store(0)
	})
	grid.Button("hud-red", "RED", 1, 0).OnClick(func() {
		colorChoice.Store(1)
	})
	grid.Button("hud-cyan", "CYAN", 2, 0).OnClick(func() {
		colorChoice.Store(2)
	})

	inputCtx := engine.Input.GetContext()

	moveLeft := notacore.Input("moveLeft", notacore.KeyA, inputCtx)
	moveRight := notacore.Input("moveRight", notacore.KeyD, inputCtx)
	moveUp := notacore.Input("moveUp", notacore.KeyW, inputCtx)
	moveDown := notacore.Input("moveDown", notacore.KeyS, inputCtx)
	winLeft := notacore.Input("winLeft", notacore.KeyQ, inputCtx)
	winRight := notacore.Input("winRight", notacore.KeyE, inputCtx)
	combo := notacore.InputCombo("combo", inputCtx, notacore.KeyE, notacore.KeyQ)

	leftClickSignal := notacore.Input("leftClick", notacore.MouseLeft, inputCtx)

	engine.Input.Start(2400)

	drawingLoop.Do(func() {
		var moveX, moveY float32

		if moveLeft.Held() {
			moveX -= 1
		}
		if moveRight.Held() {
			moveX += 1
		}
		if moveUp.Held() {
			moveY += 1
		}
		if moveDown.Held() {
			moveY -= 1
		}

		if ui.HasKeyboardFocus() {
			moveX = 0
			moveY = 0
		}

		if moveX != 0 || moveY != 0 {
			movement := notamath.Vec2{X: moveX, Y: moveY}.Mul(moveSpeed.Get())
			entity.Move(movement)
		}

		switch colorChoice.Load() {
		case 1:
			entity.WithColor(notacolor.Red)
		case 2:
			entity.WithColor(notacolor.Cyan)
		default:
			entity.WithColor(notacolor.White)
		}

		if winLeft.Held() {
			win.Move(-8, 0)
		}
		if winRight.Held() {
			win.Move(8, 0)
		}

		// Debug input testing removed - use proper logging if needed
		_ = leftClickSignal.Pressed()
		_ = combo.Pressed()

		em.Flush()
		alpha := drawingLoop.Alpha(time.Now())
		err := win.Draw(alpha, nil, entity)
		if err != nil {
			log.Printf("Draw error: %v", err)
			// Skip UI drawing if render fails
			return
		}

		// UI DRAWING NOW IN SAME FRAME
		ui.Draw()
	})

	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}
}
