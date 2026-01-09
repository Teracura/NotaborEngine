package notagl

import (
	"NotaborEngine/notamath"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type DrawOrder2D struct {
	Vertices []notamath.Po2
}

type Renderer2D struct {
	Orders []DrawOrder2D
}

func (r *Renderer2D) Begin() {
	r.Orders = r.Orders[:0]
}

type Shape2D interface {
	AddToOrders(orders *[]DrawOrder2D, alpha float32)
}

func (r *Renderer2D) Submit(s Shape2D, alpha float32) {
	s.AddToOrders(&r.Orders, alpha)
}

type vertexFormat2D struct {
	dimension int32 // should be 2
	stride    int32
}

type GLBackend2D struct {
	vao    uint32
	vbo    uint32
	format vertexFormat2D
}

func (b *GLBackend2D) Init() {
	b.format = vertexFormat2D{
		dimension: 2,
		stride:    int32(unsafe.Sizeof(notamath.Po2{})),
	}

	gl.CreateVertexArrays(1, &b.vao)
	gl.CreateBuffers(1, &b.vbo)

	gl.VertexArrayVertexBuffer(
		b.vao,
		0,
		b.vbo,
		0,
		b.format.stride,
	)

	gl.VertexArrayAttribFormat(
		b.vao,
		0,
		b.format.dimension,
		gl.FLOAT,
		false,
		0,
	)

	gl.VertexArrayAttribBinding(b.vao, 0, 0)
	gl.EnableVertexArrayAttrib(b.vao, 0)
}

func (r *Renderer2D) Flush(backend *GLBackend2D) {
	var flat []notamath.Po2
	for _, order := range r.Orders {
		flat = append(flat, order.Vertices...)
	}

	if len(flat) == 0 {
		return
	}

	gl.NamedBufferData(
		backend.vbo,
		len(flat)*int(backend.format.stride),
		gl.Ptr(flat),
		gl.DYNAMIC_DRAW,
	)

	gl.BindVertexArray(backend.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(flat)))
}
