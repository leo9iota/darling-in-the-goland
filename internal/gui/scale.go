package gui

// Scale constants for 640×360 internal resolution.
//
// All GUI elements derive sizing from these constants.
// Pixel fonts render crisply only at their native size (8px)
// or exact integer multiples (16, 24, 32...).
const (
	// FontSmall is 8px — native pixel font size, used for debug overlay.
	FontSmall = 8.0

	// FontMedium is 16px — 2× native, used for HUD coin count and menu buttons.
	FontMedium = 16.0

	// FontLarge is 24px — 3× native, used for menu title (if needed).
	FontLarge = 24.0

	// ShadowOffset is the drop shadow offset in pixels.
	ShadowOffset = 1.0
)
