package notagl

import (
	"NotaborEngine/notamath"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type DrawOrder struct {
	Vertices []notamath.Po2
}
type Renderer struct {
	Orders []DrawOrder
}

func (r *Renderer) Begin() {
	r.Orders = r.Orders[:0]
}

type Shape interface {
	AddToOrders(orders *[]DrawOrder)
}

type GLBackend struct {
	vao uint32
	vbo uint32
}

func (b *GLBackend) Init() {
	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.vbo)

	gl.BindVertexArray(b.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(
		0, // location
		2, // vector dimensions
		gl.FLOAT,
		false,
		2*4, // sizeof(Po2) (dimension amount of elements of 4 bytes)
		gl.PtrOffset(0),
	)

	gl.BindVertexArray(0)
}

func (r *Renderer) Submit(s Shape) {
	s.AddToOrders(&r.Orders)
}

func (r *Renderer) Flush(backend *GLBackend) {
	var flat []notamath.Po2

	for _, order := range r.Orders {
		flat = append(flat, order.Vertices...)
	}

	if len(flat) == 0 {
		return
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, backend.vbo)
	gl.BufferData(
		gl.ARRAY_BUFFER,
		len(flat)*int(unsafe.Sizeof(notamath.Po2{})),
		gl.Ptr(flat),
		gl.DYNAMIC_DRAW,
	)

	gl.BindVertexArray(backend.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(flat)))
}
