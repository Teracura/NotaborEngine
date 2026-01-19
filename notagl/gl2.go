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

func (r *Renderer2D) Submit(p Polygon, alpha float32) {
	var temp []DrawOrder2D

	p.AddToOrders(&temp, alpha)

	for _, order := range temp {
		tris := Triangulate2D(order.Vertices)
		if len(tris) == 0 {
			continue
		}
		r.Orders = append(r.Orders, DrawOrder2D{
			Vertices: tris,
		})
	}
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
	b.format = vertexFormat2D{dimension: 2, stride: int32(unsafe.Sizeof(notamath.Po2{}))}
	gl.CreateVertexArrays(1, &b.vao)
	gl.CreateBuffers(1, &b.vbo)
	gl.VertexArrayVertexBuffer(b.vao, 0, b.vbo, 0, b.format.stride)
	gl.VertexArrayAttribFormat(b.vao, 0, b.format.dimension, gl.FLOAT, false, 0)
	gl.VertexArrayAttribBinding(b.vao, 0, 0)
	gl.EnableVertexArrayAttrib(b.vao, 0)
}

func (b *GLBackend2D) BindVao() {
	gl.BindVertexArray(b.vao)
}

func (b *GLBackend2D) UploadData(vertices interface{}) {
	verts := vertices.([]notamath.Po2)
	gl.NamedBufferData(b.vbo, len(verts)*int(b.format.stride), gl.Ptr(verts), gl.DYNAMIC_DRAW)
}

func (r *Renderer2D) Flush(backend *GLBackend2D) {
	var flat []notamath.Po2
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

func Triangulate2D(polygon []notamath.Po2) []notamath.Po2 {
	n := len(polygon)
	if n < 3 {
		return nil
	}

	verts := append([]notamath.Po2{}, polygon...)

	// Enforce CCW winding
	if !IsCCW(verts) {
		for i, j := 0, len(verts)-1; i < j; i, j = i+1, j-1 {
			verts[i], verts[j] = verts[j], verts[i]
		}
	}

	var result []notamath.Po2

	for len(verts) > 3 {
		earFound := false

		for i := 0; i < len(verts); i++ {
			prev := verts[(i-1+len(verts))%len(verts)]
			curr := verts[i]
			next := verts[(i+1)%len(verts)]

			if IsEar(prev, curr, next, verts) {
				result = append(result, prev, curr, next)

				verts = append(verts[:i], verts[i+1:]...)
				earFound = true
				break
			}
		}

		if !earFound {
			return nil
		}
	}

	result = append(result, verts[0], verts[1], verts[2])
	return result
}
