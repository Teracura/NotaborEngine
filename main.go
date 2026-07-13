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

	// ─── Grid layout ──────────────────────────────────────────────
	// Floating HUD overlay — 30% wide, 50% tall, anchored top-right.
	// Resizing the window repositions and re-sizes it automatically.
	hud := ui.Grid("hud").Columns(3).Gap(6).Padding(8)
	hud.SetBounds(AnchorTopRight, 0.3, 0.5)

	hud.Text("hud-title", "NOTA UI").ColSpan(3).Scale(2).Color(wrapColor(Purple))
	hud.TextFunc("hud-info", func() string {
		state := "OFF"
		if hudEnabled {
			state = "ON"
		}
		return fmt.Sprintf("PLAYER %s  CLICKS %d  HUD %s", playerName, clicks.Get(), state)
	}).ColSpan(3)
	hud.Button("hud-click", "EXIT").OnClick(func() {
		win.Close()
	})
	hud.Input("hud-name", &playerName).ColSpan(2).Placeholder("name")
	hud.Slider("hud-speed", &moveStep, 0.01, 0.12).ColSpan(3).Label("SPEED").OnChange(func(v float32) {
		moveSpeed.Set(v)
	})
	hud.Checkbox("hud-toggle", "HUD ENABLED", &hudEnabled).ColSpan(3)
	hud.Button("hud-white", "WHITE").OnClick(func() {
		colorChoice.Set(0)
	})
	hud.Button("hud-red", "RED").OnClick(func() {
		colorChoice.Set(1)
	})
	hud.Button("hud-cyan", "CYAN").OnClick(func() {
		colorChoice.Set(2)
	})

	// ─── Input ────────────────────────────────────────────────────

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

		_ = leftClickSignal.Pressed()
		_ = combo.Pressed()

		em.Flush()
		alpha := drawingLoop.Alpha()
		err := win.Draw(alpha, (*Camera2D)(nil), entity)
		if err != nil {
			log.Printf("Draw error: %v", err)
			return
		}

		ui.Draw()
	})

	if err := engine.Run(); err != nil {
		log.Fatal(err)
	}
}
