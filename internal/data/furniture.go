package data

type Shape int

const (
	Rectangle Shape = iota
	Circle
)

type Furniture struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	X      int64  `json:"x"` // X coordinate relative to Room's left position.
	Y      int64  `json:"y"` // Y coordinate relative to Room's top position.
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Image  string `json:"image,omitempty"` // Path to the image
	Shape  Shape  `json:"shape"`           // To improve collision detection
}
