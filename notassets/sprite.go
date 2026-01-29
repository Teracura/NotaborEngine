package notassets

import (
	"NotaborEngine/notagl"
)

type Sprite struct {
	Texture             *notagl.Texture
	Name                string
	X, Y                int32
	srcWidth, srcHeight int32
}
