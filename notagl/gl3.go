package notagl

import (
	"NotaborEngine/notamath"
	"NotaborEngine/notashader"
	"unsafe"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Vertex3D struct {
	Pos   notamath.Po3
	Color notashader.Color
}
type DrawOrder3D struct {
	Vertices []notamath.Po3
}

type Renderer3D struct {
	Orders []DrawOrder3D
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
	provideGlSettings(&b.vao, &b.vbo, b.format.stride, b.format.dimension)
}

func provideGlSettings(vao *uint32, vbo *uint32, stride int32, dimension int32) {
	gl.CreateVertexArrays(1, vao)
	gl.CreateBuffers(1, vbo)
	gl.VertexArrayVertexBuffer(*vao, 0, *vbo, 0, stride)
	gl.VertexArrayAttribFormat(*vao, 0, dimension, gl.FLOAT, false, 0)
	gl.VertexArrayAttribBinding(*vao, 0, 0)
	gl.EnableVertexArrayAttrib(*vao, 0)
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
