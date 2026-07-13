package main

import (
	"fmt"
	"log"
)

func main() {
	// Engine setup
	settings := &Settings{
		Vsync:      true,
		SoundLevel: 1,
		Muted:      false,
	}

	drawingLoop := NewLoop(60)

	engine, err := CreateEngine(settings)
	if err != nil {
		log.Fatal(err)
	}
	defer engine.Shutdown()

	cfg := &WindowConfig{
		X:         50,
		Y:         50,
		W:         800,
		H:         600,
		Title:     "Entity Test",
		Type:      WindowWindowed,
		Resizable: true,
		TargetFPS: 60,
	}

	cfg.Loops = append(cfg.Loops, drawingLoop.handle)

	win, err := engine.CreateWindow(cfg)
	if err != nil {
		log.Fatal(err)
	}

	em := engine.EntityManager()
	circleRadius := float32(0.25)

	ballVisual, err := win.LoadVisual("quadSprite", "resources/images/hahaha.jpg", VisualOptions{
		Width:        circleRadius * 2,
		Height:       circleRadius * 2,
		Mask:         MaskCircle,
		CircleRadius: 0.5,
		CircleEdge:   0.01,
	})

	if err != nil {
		log.Fatal(err)
	}

	entity := em.CreateEntity("quad").
		WithVisual(ballVisual).
		WithCollision(CircleCollision(circleRadius)).
		WithColor(White)

	moveStep := float32(0.05)
	var moveSpeed AtomicFloat32
	moveSpeed.Set(moveStep)

	ui, err := NewUI(engine, win)
	if err != nil {
		log.Fatal(err)
	}

	var clicks AtomicInt32
	var colorChoice AtomicInt32
	playerName := ""
	hudEnabled := true

	ui.Panel("hud-panel").Rect(14, 14, 314, 220)
	ui.Text("hud-title", "NOTA UI").At(24, 24).Scale(2).Color(wrapColor(Purple))
	ui.TextFunc("hud-info", func() string {
		state := "OFF"
		if hudEnabled {
			state = "ON"
		}
		return fmt.Sprintf("PLAYER %s  CLICKS %d  HUD %s", playerName, clicks.Get(), state)
	}).At(24, 48)
	ui.Button("hud-click", "EXIT").Rect(24, 72, 92, 28).OnClick(func() {
		win.Close()
	})
	ui.Input("hud-name", &playerName).Rect(124, 72, 184, 28).Placeholder("name")
	ui.Slider("hud-speed", &moveStep, 0.01, 0.12).Rect(24, 114, 284, 34).Label("SPEED").OnChange(func(v float32) {
		moveSpeed.Set(v)
	})
	ui.Checkbox("hud-toggle", "HUD ENABLED", &hudEnabled).Rect(24, 154, 160, 24)

	grid := ui.Grid("hud-color-grid", R(24, 188, 284, 34), 3, 1).Gap(8)
	grid.Button("hud-white", "WHITE", 0, 0).OnClick(func() {
		colorChoice.Set(0)
	})
	grid.Button("hud-red", "RED", 1, 0).OnClick(func() {
		colorChoice.Set(1)
	})
	grid.Button("hud-cyan", "CYAN", 2, 0).OnClick(func() {
		colorChoice.Set(2)
	})

	inputCtx := engine.Input().GetContext()

	moveLeft := Input("moveLeft", KeyA, inputCtx)
	moveRight := Input("moveRight", KeyD, inputCtx)
	moveUp := Input("moveUp", KeyW, inputCtx)
	moveDown := Input("moveDown", KeyS, inputCtx)
	winLeft := Input("winLeft", KeyQ, inputCtx)
	winRight := Input("winRight", KeyE, inputCtx)
	combo := InputCombo("combo", inputCtx, KeyE, KeyQ)

	leftClickSignal := Input("leftClick", MouseLeft, inputCtx)

	engine.Input().Start(2400)

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
			movement := Vec2{X: moveX, Y: moveY}.Mul(moveSpeed.Get())
			entity.Move(movement)
		}

		switch colorChoice.Get() {
		case 1:
			entity.WithColor(Red)
		case 2:
			entity.WithColor(Cyan)
		default:
			entity.WithColor(White)
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
		alpha := drawingLoop.Alpha()
		err := win.Draw(alpha, (*Camera2D)(nil), entity)
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
