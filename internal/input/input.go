package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// IsDown returns true if any of the given keys are currently held down.
func IsDown(keys ...ebiten.Key) bool {
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

// JustPressed returns true if the key was pressed this frame.
func JustPressed(key ebiten.Key) bool {
	return inpututil.IsKeyJustPressed(key)
}
