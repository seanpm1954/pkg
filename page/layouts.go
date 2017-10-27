package page

import "fmt"

var (
	// layout that will be used in the Render function - defaults to Standard
	Layout = Standard
)

// Layout represents
type Key string

// Suffix adds a suffix to the layout
func (l Key) Suffix(key string) string {
	return fmt.Sprintf("%s.%s", l, key)
}

// Available layouts
const (
	Standard Key = "standard"
)

// SetLayout sets which layout the render method will use
func SetLayout(layout Key) {
	switch layout {
	case Standard:
		layout = Standard
	default:
		panic(fmt.Sprintf("invalid layout: %v", layout))
	}
}
