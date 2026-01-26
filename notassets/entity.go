package notassets

import (
	"NotaborEngine/notacollision"
	"NotaborEngine/notagl"
)

type Entity struct {
	Sprite   *Sprite
	Polygon  notagl.Polygon
	Collider notacollision.Collider
}
