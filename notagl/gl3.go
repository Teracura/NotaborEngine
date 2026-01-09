package notagl

import (
	"NotaborEngine/notamath"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type DrawOrder3D struct {
	Vertices []notamath.Po3
}

type Renderer3D struct {
	Orders []DrawOrder3D
}

func (r *Renderer3D) Begin() {
	r.Orders = r.Orders[:0]
}

type Shape3D interface {
	AddToOrders(orders *[]DrawOrder3D)
}

func (r *Renderer3D) Submit(s Shape3D) {
	s.AddToOrders(&r.Orders)
}

type vertexFormat3D struct {
	dimension int32 // should be 3
	stride    int32
}

type GLBackend3D struct {
	vao    uint32
	vbo    uint32
	format vertexFormat3D
}

func (b *GLBackend3D) Init() {
	b.format = vertexFormat3D{dimension: 3, stride: int32(unsafe.Sizeof(notamath.Po3{}))}
	gl.CreateVertexArrays(1, &b.vao)
	gl.CreateBuffers(1, &b.vbo)
	gl.VertexArrayVertexBuffer(b.vao, 0, b.vbo, 0, b.format.stride)
	gl.VertexArrayAttribFormat(b.vao, 0, b.format.dimension, gl.FLOAT, false, 0)
	gl.VertexArrayAttribBinding(b.vao, 0, 0)
	gl.EnableVertexArrayAttrib(b.vao, 0)
}

func (b *GLBackend3D) BindVao() {
	gl.BindVertexArray(b.vao)
}

func (b *GLBackend3D) UploadData(vertices interface{}) {
	verts := vertices.([]notamath.Po3)
	gl.NamedBufferData(b.vbo, len(verts)*int(b.format.stride), gl.Ptr(verts), gl.DYNAMIC_DRAW)
}

func (r *Renderer3D) Flush(backend *GLBackend3D) {
	var flat []notamath.Po3
	for _, order := range r.Orders {
		flat = append(flat, order.Vertices...)
	}
	if len(flat) == 0 {
		return
	}
	backend.UploadData(flat)
	backend.BindVao()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(flat)))
}
