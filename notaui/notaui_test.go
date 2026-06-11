package notaui

import "testing"

func TestRectContains(t *testing.T) {
	r := R(10, 20, 100, 50)

	if !r.Contains(10, 20) {
		t.Fatal("expected top-left edge to be inside")
	}
	if !r.Contains(110, 70) {
		t.Fatal("expected bottom-right edge to be inside")
	}
	if r.Contains(111, 70) {
		t.Fatal("expected point past right edge to be outside")
	}
	if r.Contains(10, 71) {
		t.Fatal("expected point past bottom edge to be outside")
	}
}

func TestGridCell(t *testing.T) {
	grid := (&Grid{
		bounds:  R(10, 20, 300, 200),
		columns: 3,
		rows:    2,
	}).Gap(10).Padding(5)

	got := grid.Cell(1, 1)
	want := R(115, 125, 90, 90)

	if got != want {
		t.Fatalf("unexpected cell: got %+v want %+v", got, want)
	}
}

func TestGridCellClamps(t *testing.T) {
	grid := &Grid{
		bounds:  R(0, 0, 100, 100),
		columns: 2,
		rows:    2,
	}

	if got, want := grid.Cell(-5, 99), R(0, 50, 50, 50); got != want {
		t.Fatalf("unexpected clamped cell: got %+v want %+v", got, want)
	}
}

func TestSliderSetFromMouseClamps(t *testing.T) {
	value := float32(0)
	slider := &Slider{
		value:  &value,
		min:    10,
		max:    20,
		bounds: R(100, 0, 200, 20),
	}

	slider.setFromMouse(50)
	if value != 10 {
		t.Fatalf("expected min value, got %f", value)
	}

	slider.setFromMouse(400)
	if value != 20 {
		t.Fatalf("expected max value, got %f", value)
	}

	slider.setFromMouse(200)
	if value != 15 {
		t.Fatalf("expected midpoint value, got %f", value)
	}
}

func TestSliderOnChangeCallback(t *testing.T) {
	value := float32(0)
	got := float32(-1)
	slider := (&Slider{
		value:  &value,
		min:    0,
		max:    100,
		bounds: R(0, 0, 100, 20),
	}).OnChange(func(v float32) {
		got = v
	})

	callback := slider.setFromMouse(25)
	if callback == nil {
		t.Fatal("expected callback")
	}
	callback()

	if got != 25 {
		t.Fatalf("expected callback value 25, got %f", got)
	}
}
