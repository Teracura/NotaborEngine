package notagl

type Renderer[T any] interface {
	Begin()
	Flush(backend *T)
}

type GlBackend interface {
	Init()
	BindVao()
	UploadData(vertices interface{})
}
