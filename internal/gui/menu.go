package gui

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Menu constants — scaled for 640×360 internal resolution.
const (
	menuWidth        = 200.0
	menuHeight       = 150.0
	menuButtonWidth  = 150.0
	menuButtonHeight = 28.0
	menuButtonMargin = 10.0
	menuFontSize     = 10.0
)

// Button is a clickable menu button.
type Button struct {
	Text   string
	X, Y   float64
	W, H   float64
	Hover  bool
	Action func() error
}

// Menu is a pause menu overlay with buttons.
type Menu struct {
	Active  bool
	buttons []*Button
	font    *text.GoTextFace
	x, y    float64
}

// NewMenu creates a centered pause menu.
func NewMenu(screenW, screenH float64, toggleFn func()) (*Menu, error) {
	fontSrc, err := LoadFont("assets/fonts/public-pixel-font.ttf")
	if err != nil {
		return nil, err
	}

	face := &text.GoTextFace{
		Source: fontSrc,
		Size:   menuFontSize,
	}

	mx := (screenW - menuWidth) / 2
	my := (screenH - menuHeight) / 2
	btnX := mx + (menuWidth-menuButtonWidth)/2

	m := &Menu{
		font: face,
		x:    mx,
		y:    my,
	}

	m.buttons = []*Button{
		{
			Text:   "Resume",
			X:      btnX,
			Y:      my + 25,
			W:      menuButtonWidth,
			H:      menuButtonHeight,
			Action: func() error { toggleFn(); return nil },
		},
		{
			Text:   "Settings",
			X:      btnX,
			Y:      my + 25 + menuButtonHeight + menuButtonMargin,
			W:      menuButtonWidth,
			H:      menuButtonHeight,
			Action: func() error { log.Println("Settings not implemented"); return nil },
		},
		{
			Text:   "Quit",
			X:      btnX,
			Y:      my + 25 + (menuButtonHeight+menuButtonMargin)*2,
			W:      menuButtonWidth,
			H:      menuButtonHeight,
			Action: func() error { return ebiten.Termination },
		},
	}

	return m, nil
}

// Toggle flips the menu active state.
func (m *Menu) Toggle() {
	m.Active = !m.Active
}

// Update checks mouse hover on buttons.
func (m *Menu) Update() {
	if !m.Active {
		return
	}

	mx, my := ebiten.CursorPosition()
	fx, fy := float64(mx), float64(my)

	for _, b := range m.buttons {
		b.Hover = fx >= b.X && fx <= b.X+b.W && fy >= b.Y && fy <= b.Y+b.H
	}
}

// Draw renders the pause menu overlay.
func (m *Menu) Draw(screen *ebiten.Image) {
	if !m.Active {
		return
	}

	// Semi-transparent background
	vector.DrawFilledRect(screen, float32(m.x), float32(m.y), menuWidth, menuHeight, color.NRGBA{51, 51, 51, 230}, true)

	// Border
	vector.StrokeRect(screen, float32(m.x), float32(m.y), menuWidth, menuHeight, 1, color.White, true)

	// Title
	title := "Menu"
	tw, _ := MeasureText(title, m.font)
	DrawText(screen, title, m.x+(menuWidth-tw)/2, m.y+10, m.font, color.White)

	// Buttons
	for _, b := range m.buttons {
		// Button background
		bgClr := color.NRGBA{77, 77, 77, 255}
		if b.Hover {
			bgClr = color.NRGBA{102, 102, 102, 255}
		}
		vector.DrawFilledRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), bgClr, true)

		// Button border
		vector.StrokeRect(screen, float32(b.X), float32(b.Y), float32(b.W), float32(b.H), 1, color.White, true)

		// Button text (centered)
		tw, th := MeasureText(b.Text, m.font)
		DrawText(screen, b.Text, b.X+(b.W-tw)/2, b.Y+(b.H-th)/2, m.font, color.White)
	}
}

// MousePressed checks if a button was clicked and fires its action.
func (m *Menu) MousePressed(x, y int) error {
	if !m.Active {
		return nil
	}

	fx, fy := float64(x), float64(y)
	for _, b := range m.buttons {
		if fx >= b.X && fx <= b.X+b.W && fy >= b.Y && fy <= b.Y+b.H {
			return b.Action()
		}
	}
	return nil
}
