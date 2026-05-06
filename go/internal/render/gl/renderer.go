package gl

import "fmt"

type Renderer struct{}

func (Renderer) Init() error {
	return fmt.Errorf("go-gl renderer not yet wired")
}
