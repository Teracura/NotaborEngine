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
	Vertices []Vertex3D
}

type Renderer3D struct {
	Orders []DrawOrder3D
}

func (r *Renderer3D) Submit(m *Mesh, alpha float32) {
	m.AddToOrders(&r.Orders, alpha)
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
	b.format = vertexFormat3D{
		dimension: 3,
		stride:    int32(unsafe.Sizeof(Vertex3D{})),
	}

	gl.CreateVertexArrays(1, &b.vao)
	gl.CreateBuffers(1, &b.vbo)

	gl.VertexArrayVertexBuffer(b.vao, 0, b.vbo, 0, b.format.stride)

	gl.VertexArrayAttribFormat(
		b.vao,
		0,
		3,
		gl.FLOAT,
		false,
		0,
	)
	gl.VertexArrayAttribBinding(b.vao, 0, 0)
	gl.EnableVertexArrayAttrib(b.vao, 0)

	colorOffset := uint32(unsafe.Sizeof(notamath.Po3{}))

	gl.VertexArrayAttribFormat(
		b.vao,
		1,
		4,
		gl.FLOAT,
		false,
		colorOffset,
	)
	gl.VertexArrayAttribBinding(b.vao, 1, 0)
	gl.EnableVertexArrayAttrib(b.vao, 1)
}

func (b *GLBackend3D) UploadData(vertices interface{}) {
	verts := vertices.([]Vertex3D)
	gl.NamedBufferData(
		b.vbo,
		len(verts)*int(b.format.stride),
		gl.Ptr(verts),
		gl.DYNAMIC_DRAW,
	)
}

func (r *Renderer3D) Flush(backend *GLBackend3D) {
	var flat []Vertex3D
	for _, order := range r.Orders {
		tris := order.Vertices
		flat = append(flat, tris...)
	}

	if len(flat) == 0 {
		return
	}
	backend.UploadData(flat)
	gl.BindVertexArray(backend.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(flat)))
}
